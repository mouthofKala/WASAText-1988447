package database

import (
	"errors"
	"log"

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

	// db check: does target chat exist?
	var chatExists bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chats
			WHERE chat_id = ?)`, chat_id).Scan(&chatExists)
	if !errors.Is(err, nil) {
		log.Printf("error checking chat existence: %v", err)
		return Message{}, err // 500
	}
	if !chatExists {
		log.Printf("chat may not exist in DB: %v", err)
		return Message{}, ErrBadReq
	}

	// db check: is userID a member?
	var membership bool
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chat_members
			WHERE chat_id = ? AND user_id = ?)`,
		chat_id, userID).Scan(&membership)
	if !errors.Is(err, nil) {
		log.Printf("error fetching user membership in db: %v", err)
		return Message{}, err // 500
	}
	if !membership {
		log.Printf("user may not be a member in DB: %v", err)
		return Message{}, ErrNotaGroupMember // 403
	}

	if replyto != nil {
		var replytoexist bool
		err := db.c.QueryRow(`
			SELECT EXISTS(
				SELECT 1
				FROM messages
				WHERE message_id = ? AND chat_id = ?)`,
			replyto, chat_id).Scan(&replytoexist)
		if err != nil {
			log.Printf("error fetching msg user is replying to: %v", err)
			return Message{}, err
		}
		if !replytoexist {
			log.Printf("msg reply to refers to does not exist: %v", err)
			return Message{}, ErrInvalidMsg
		}
	}

	// is message VALID (must have at least text or img)
	if (content == nil || *content == "") && photo == nil {
		log.Printf("msg is empty!!!: %v", err)
		return Message{}, ErrInvalidMsg // 400
	}

	// fetch username
	username := ""
	err = db.c.QueryRow(`
		SELECT username
		FROM users
		WHERE user_id = ?
	`, userID).Scan(&username)
	if err != nil {
		log.Printf("this user has no registered username??: %v", err)
		return Message{}, err
	}

	timestamp := globaltime.Now()

	tx, err := db.c.Begin()
	if !errors.Is(err, nil) {
		return Message{}, ErrTX
	}
	defer func() {
		_ = tx.Rollback()
	}()

	//  put msg in db ...
	_, err = tx.Exec(`
		INSERT INTO messages (
			message_id,
			chat_id,
			sender_id,
			sender_n,
			content,
			photo,
			status,
			timestamp,
			reply_to,
			fwd_from
			)
		VALUES (?,?,?,?,?,?,?,?,?,?)`,
		messageID, chat_id, userID, username, content, photo, "sent", timestamp, replyto, fwdfrom)
	if !errors.Is(err, nil) {
		log.Printf("error inserting msg in db: %v", err)
		return Message{}, ErrExec //  gneeric 500 bc uniquely identifiable
	}

	//  CREATE NOTREADS
	//  fetch members of the chat except sender:
	rows, err := tx.Query(`
		SELECT user_id
		FROM chat_members
		WHERE chat_id = ? AND user_id != ?`, chat_id, userID)
	if !errors.Is(err, nil) {
		log.Printf("error fetching users for unreads: %v", err)
		return Message{}, err //  500
	}
	defer func() {
		if errclose := rows.Close(); errclose != nil {
			log.Printf("error closing rows: %v", errclose)
		}
	}()

	for rows.Next() {
		var memberID string
		if err = rows.Scan(&memberID); !errors.Is(err, nil) {
			return Message{}, ErrTX
		}
		_, err = tx.Exec(`
			INSERT INTO msg_notreads(message_id, chat_id, user_id)
			VALUES (?,?,?)`, messageID, chat_id, memberID)

		if !errors.Is(err, nil) {
			log.Printf("error adding msg unreads to db: %v", err)
			return Message{}, ErrExec // 500
		}
	}
	if err := rows.Err(); !errors.Is(err, nil) {
		log.Printf("error during transaction: %v", err)
		return Message{}, ErrTX // 500
	}
	if err := tx.Commit(); !errors.Is(err, nil) {
		return Message{}, ErrTX
	}

	return Message{
		MessageID: messageID,
		ChatID:    chat_id,
		SenderID:  userID,
		SenderN:   username,
		Content:   content,
		Photo:     photo,
		Status:    "sent",
		Timestamp: timestamp,
		ReplyTo:   replyto,
		FwdFrom:   fwdfrom,
		Reacts:    []Reaction{},
	}, nil
}
