package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) ForwardMessage(userID, targetchat_id string, fwdmsgID string, chatID string, newmsgID string) (Message, error) {
	//check that both chats exist

	//check that the user is a member of chat_id and targetchat_id, otherwise 403
	var membership1 bool
	var membership2 bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chat_members
			WHERE user_id = ? AND chat_id = ?)`, userID, chatID).Scan(&membership1)
	if err != nil {
		return Message{}, err //500
	}
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chat_members
			WHERE user_id = ? AND chat_id = ?)`, userID, targetchat_id).Scan(&membership2)
	if err != nil {
		return Message{}, err
	}
	if !(membership1 && membership2) {
		return Message{}, ErrForbidden
	}

	//check that the msg exists, otherwise 400
	var fwdmsg Message
	err = db.c.QueryRow(`
		SELECT user_id, content, photo, timestamp
		FROM messages
		WHERE chat_id = ? AND message_id = ?	
	`, chatID, fwdmsgID).Scan(
		&fwdmsg.SenderID,
		&fwdmsg.Content,
		&fwdmsg.Photo,
		&fwdmsg.Timestamp)
	if errors.Is(err, sql.ErrNoRows) {
		return Message{}, ErrBadReq
	}
	if err != nil {
		return Message{}, err //500
	}

	msg, err := db.SendMessage(userID, targetchat_id, fwdmsg.Content, fwdmsg.Photo, nil, &chatID, newmsgID)
	if err != nil {
		return Message{}, err
	}
	//use snedmsg as a helper
	return msg, nil
}
