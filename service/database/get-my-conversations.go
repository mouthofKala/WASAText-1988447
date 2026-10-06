package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetMyConversations(userID string) ([]Chat, error) {
	var chats []Chat
	rows, err := db.c.Query(`
		SELECT
			chats.chat_name,
			chats.chat_id,
			chats.photo,
			chats.creation_date,
			chats.group_or_chat,
			messages.message_id,
			messages.sender_id,
			messages.content,
			messages.photo,
			messages.timestamp

		FROM chats
		JOIN chat_members ON chats.chat_id=chat_members.chat_id
		LEFT JOIN messages ON messages.message_id = (
			SELECT m.message_id
			FROM messages m
			WHERE m.chat_id = chats.chat_id
			ORDER BY m.timestamp DESC
			LIMIT 1	
		)
		WHERE chat_members.user_id=?
		ORDER BY messages.timestamp DESC
	`, userID)
	// LEFT JOIN to return even if the chat has 0 msg
	if !errors.Is(err, nil) {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat Chat
		var preview MessagePreview
		var messageID sql.NullString
		var senderID sql.NullString
		var content sql.NullString
		var photo sql.NullString
		var timestamp sql.NullTime

		err = rows.Scan(
			&chat.ChatName,
			&chat.ChatID,
			&chat.Photo,
			&chat.CreationDate,
			&chat.GroupOrChat,
			&messageID, //   msg may not exist
			&senderID,
			&content,
			&photo,
			&timestamp,
		)
		if !errors.Is(err, nil) {
			return nil, err
		}

		if messageID.Valid {
			preview.MessageID = messageID.String
			preview.SenderID = senderID.String
			if timestamp.Valid {
				preview.Timestamp = timestamp.Time
			}
			if timestamp.Valid {
				preview.Timestamp = timestamp.Time
			}
			if content.Valid {
				preview.Content = &content.String
			}
			if photo.Valid {
				preview.Photo = &photo.String
			}
			chat.MostRecentMsg = &preview
		}
		chats = append(chats, chat)
	}
	if err := rows.Err(); !errors.Is(err, nil) {
		return nil, err
	}
	return chats, nil
}
