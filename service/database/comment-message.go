package database

import (
	"database/sql"
)

func (db *appdbimpl) CommentMessage(msgID string, chatID string, userID string, reactionID string, emoji rune) (Reaction, error) {
	//does chatID exist? 400
	var exist bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chats
			WHERE chat_id = ?)`, chatID).Scan(&exist)
	if !exist {
		return Reaction{}, ErrBadReq
	}
	if err != nil {
		return Reaction{}, err
	}

	//is userID in chat ID? 403
	var membership bool
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chat_members
			WHERE chat_id = ? AND user_id = ?)`,
		chatID, userID).Scan(&membership)
	if !membership {
		return Reaction{}, ErrForbidden
	}
	if err != nil {
		return Reaction{}, err
	}

	//does msgID exist? 400
	var msgexists bool
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM messages
			WHERE chat_id = ? AND message_id = ?)`,
		chatID, msgID).Scan(&msgexists)
	if !msgexists {
		return Reaction{}, ErrBadReq
	}
	if err != nil {
		return Reaction{}, err
	}

	//is there a reaction already?
	var reac Reaction
	err = db.c.QueryRow(`
		SELECT reaction_id
		FROM reactions
		WHERE message_id = ? AND user_id = ?`,
		msgID, userID).Scan(&reac.ReactionID)
	if err == sql.ErrNoRows {
		//there is NOT a reaction already
		_, err = db.c.Exec(`
			INSERT INTO reactions
			(reaction_id, message_id, chat_id, user_id, emoji)
			VALUES (?,?,?,?,?)`,
			reactionID, msgID, chatID, userID, string(emoji))
		if err != nil {
			return Reaction{}, ErrExec
		}
		return Reaction{
			ReactionID: reactionID,
			MessageID:  msgID,
			ChatID:     chatID,
			UserID:     userID,
			Emoji:      string(emoji),
		}, nil
	}
	if err != nil {
		return Reaction{}, err //500
	}

	_, err = db.c.Exec(`
		UPDATE reactions
		SET emoji = ?
		WHERE reaction_id = ?
			AND message_id = ?
			AND user_id = ?
	`, string(emoji), reac.ReactionID, msgID, userID)
	if err != nil {
		return Reaction{}, ErrExec
	}

	//make a reaction ID in handler

	return Reaction{
		ReactionID: reac.ReactionID,
		MessageID:  msgID,
		ChatID:     chatID,
		UserID:     userID,
		Emoji:      string(emoji),
	}, nil
}
