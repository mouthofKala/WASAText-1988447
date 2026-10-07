package database

import (
	"errors"
)

func (db *appdbimpl) GetConversation(userID string, chatID string, before string) (Conversation, error) {
	var convo Conversation

	err := db.c.QueryRow(`
		SELECT
			c.chat_id,
			c.group_or_chat,
			c.chat_name,
			c.creation_date,
			c.photo
		FROM chats c
		JOIN chat_members cm
			ON c.chat_id = cm.chat_id
		WHERE c.chat_id = ?
		AND cm.user_id = ?`, chatID, userID).Scan(
		&convo.Chat.ChatID,
		&convo.Chat.GroupOrChat,
		&convo.Chat.ChatName,
		&convo.Chat.CreationDate,
		&convo.Chat.Photo,
	)

	if convo.Chat.ChatID == "" {
		return Conversation{}, ErrFetchingChat
	}

	if !errors.Is(err, nil) {
		return Conversation{}, err // generic 500
	}

	// search members and most recent msg
	rows, err := db.c.Query(`
			SELECT user_id
			FROM chat_members
			WHERE chat_id = ?
		`, chatID)
	if !errors.Is(err, nil) {
		return Conversation{}, ErrFetchingMembers // 500 this includes norows
	}
	defer rows.Close()

	for rows.Next() {
		var memberID string
		if err := rows.Scan(&memberID); !errors.Is(err, nil) {
			return Conversation{}, ErrFetchingMembers // 500
		}
		convo.Chat.Members = append(convo.Chat.Members, memberID)
	}
	if err := rows.Err(); !errors.Is(err, nil) {
		return Conversation{}, ErrFetchingMembers // 500
	}

	// search message history
	if before == "" {
		rows, err = db.c.Query(`
			SELECT
				message_id,
				chat_id,
				sender_id,
				content,
				photo,
				status,
				timestamp,
				reply_to,
				fwd_from
			FROM (
				SELECT *
				FROM messages
				WHERE chat_id = ?
				ORDER BY timestamp DESC
				LIMIT 50
			)
			ORDER BY timestamp ASC;
			`, chatID)
	} else {
		rows, err = db.c.Query(`
			SELECT
				message_id,
				chat_id,
				sender_id,
				content,
				photo,
				status,
				timestamp,
				reply_to,
				fwd_from
			FROM (
				SELECT *
				FROM messages
				WHERE chat_id = ?
					AND timestamp < ?
				ORDER BY timestamp DESC
				LIMIT 50
			)
			ORDER BY timestamp ASC;
			`, chatID, before)
	}

	if !errors.Is(err, nil) {
		return Conversation{}, ErrFetchingMessages
	}
	defer rows.Close()
	for rows.Next() {
		var message Message

		if err := rows.Scan(
			&message.MessageID,
			&message.ChatID,
			&message.SenderID,
			&message.Content,
			&message.Photo,
			&message.Status,
			&message.Timestamp,
			&message.ReplyTo,
			&message.FwdFrom,
		); !errors.Is(err, nil) {
			return Conversation{}, ErrFetchingMessages // 500?
		}

		convo.Messages = append(convo.Messages, message)
	}
	if err := rows.Err(); !errors.Is(err, nil) {
		return Conversation{}, ErrFetchingMessages // 500
	}

	// depopulate reads in this chat for this user
	_, err = db.c.Exec(`
		DELETE FROM msg_notreads
		WHERE chat_id = ? AND user_id = ?`, chatID, userID)
	if !errors.Is(err, nil) {
		return Conversation{}, err // 500
	}

	// now with the fetched messages, for each message we need to fetch rows in unreads
	for i := range convo.Messages {
		var gotrows bool
		err = db.c.QueryRow(`
			SELECT EXISTS(
				SELECT 1
				FROM msg_notreads
				WHERE message_id = ?)`, convo.Messages[i].MessageID).Scan(&gotrows)
		if !errors.Is(err, nil) {
			return Conversation{}, err
		}
		if !gotrows {
			_, err = db.c.Exec(`
				UPDATE messages
				SET status = ?
				WHERE message_id = ?`, "read", convo.Messages[i].MessageID)
			if !errors.Is(err, nil) {
				return Conversation{}, err
			}
			convo.Messages[i].Status = "read"
		}
	}

	return convo, nil
}
