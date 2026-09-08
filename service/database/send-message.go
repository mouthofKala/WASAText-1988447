package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/globaltime"
)

func (db *appdbimpl) SendMessage(
	userID string,
	chat_id string,
	content *string,
	photo *string,
	replyto *string,
	fwdfrom *string,
	messageID string) (Message, error) {
	//db check: does target chat exist?
	var chatExists bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chats
			WHERE chat_id = ?)`, chat_id).Scan(&chatExists)
	if err != nil {
		return Message{}, err //500
	}
	if !chatExists {
		return Message{}, ErrBadReq //MODIFY IN API HANDLER
	}

	//db check: is userID a member?
	var membership bool
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chat_members
			WHERE chat_id = ? AND user_id = ?)`,
		chat_id, userID).Scan(&membership)
	if err != nil {
		return Message{}, err //500
	}
	if !membership {
		return Message{}, ErrNotaGroupMember //403
	}

	//target chat existing is implied if you assume there are no
	//"illegal" actions on the db which leave with a deleted chat
	//but no deleted membership, so it's best to check both dbs

	//is message VALID (must have at least text or img)
	if (content == nil || *content == "") && photo == nil {
		return Message{}, ErrInvalidMsg //400
	}

	timestamp := globaltime.Now()

	tx, err := db.c.Begin()
	if err != nil {
		return Message{}, ErrExec
	}

	//put msg in db ...
	_, err = tx.Exec(`
		INSERT INTO messages (
			message_id,
			chat_id,
			sender_id,
			content,
			photo,
			status,
			timestamp,
			reply_to,
			fwd_from)
		VALUES (?,?,?,?,?,?,?,?,?)`,
		messageID, chat_id, userID, content, photo, "sent", timestamp, replyto, fwdfrom)
	if err != nil {
		return Message{}, ErrExec //gneeric 500 bc uniquely identifiable
	}

	//CREATE NOTREADS
	//fetch members of the chat except sender:
	rows, err := tx.Query(`
		SELECT user_id
		FROM chat_members
		WHERE chat_id = ? AND user_id != ?`, chat_id, userID)
	if err != nil {
		return Message{}, err //500
	}
	defer rows.Close()

	for rows.Next() {
		var memberID string
		if err = rows.Scan(&memberID); err != nil {
			return Message{}, ErrExec
		}
		_, err = tx.Exec(`
			INSERT INTO msg_notreads(message_id, chat_id, user_id)
			VALUES (?,?,?)`, messageID, chat_id, memberID)

		if err != nil {
			return Message{}, ErrExec //500
		}
	}
	if err := rows.Err(); err != nil {
		return Message{}, ErrExec //500
	}
	if err := tx.Commit(); err != nil {
		return Message{}, ErrExec
	}

	return Message{
		MessageID: messageID,
		ChatID:    chat_id,
		SenderID:  userID,
		Content:   content,
		Photo:     photo,
		Status:    "sent",
		Timestamp: timestamp,
		ReplyTo:   replyto,
		FwdFrom:   fwdfrom,
	}, nil
}
