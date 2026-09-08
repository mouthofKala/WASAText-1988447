package database

import (
	"database/sql"
	"errors"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/globaltime"
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) MakeChat(otheruser string, creatorID string) (Chat, error) {

	exists, err := db.userExists(otheruser)
	if err != nil {
		return Chat{}, err //bad request
	}
	if !exists {
		return Chat{}, ErrUserNotFound //bad request
	}

	//2. DOES A PRIVATE CHAT WITH THIS USER ALREADY EXIST?
	var existingChatID string
	err = db.c.QueryRow(`
		SELECT cm1.chat_id
		FROM chat_members cm1
		JOIN chat_members cm2
			ON cm1.chat_id = cm2.chat_id
		WHERE cm1.user_id = ?
			AND cm2.user_id = ?
			AND (
				SELECT COUNT(*)
				FROM chat_members cm3
				WHERE cm3.chat_id = cm1.chat_id	
			) = 2
		LIMIT 1
	`, creatorID, otheruser).Scan(&existingChatID)

	if err == nil {
		return Chat{}, ErrChatAlreadyExists //conflict
	}
	if errors.Is(err, sql.ErrNoRows) {
		return Chat{}, err
	}

	//get other user's username:
	var otherusername string
	err = db.c.QueryRow(`
		SELECT username
		FROM users
		WHRERE user_id = ?	
	`, otheruser).Scan(&otherusername)

	//make chat in database
	chatID, err := uuid.NewV4()
	if err != nil {
		return Chat{}, ErrGeneratingID //internal server error w/ msg
	}
	creationdate := globaltime.Now()

	tx, err := db.c.Begin()
	if err != nil {
		return Chat{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.Exec(`
		INSERT INTO chats(
			chat_id,
			group_or_chat,
			chat_name,
			creation_date,
			photo)
		VALUES (?,?,?,?,?)
	`, chatID.String(), "private", otheruser, creationdate, "")

	if err != nil {
		return Chat{}, err
	}

	_, err = tx.Exec(`
		INSERT INTO chat_members (chat_id, user_id)
		VALUES (?,?),(?,?)
	`, chatID.String(), creatorID, chatID.String(), otheruser)
	if err != nil {
		return Chat{}, err
	}

	if err = tx.Commit(); err != nil {
		return Chat{}, err
	}

	//return chat
	return Chat{
		ChatID:       chatID.String(),
		GroupOrChat:  "private",
		ChatName:     otherusername,
		CreationDate: creationdate,
		Photo:        "",
	}, nil

}
