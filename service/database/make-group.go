package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/globaltime"
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) MakeGroup(chatname string, photouri string, members []string) (Chat, error) {
	if len(members) == 0 {
		return Chat{}, ErrNoMembersSelected
	}

	//CHECK FOR DUPLICATE USERS?
	for _, memberID := range members {
		exists, err := db.userExists(memberID)
		if err != nil {
			return Chat{}, err
		}
		if !exists {
			return Chat{}, ErrUserNotFound
		}
	}

	if photouri == "" {
		photouri = defaultGroupPhoto
	}
	//make the group
	chatID, err := uuid.NewV4()
	if err != nil {
		return Chat{}, err
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
		INSERT INTO chats (
			chat_id,
			group_or_chat,
			chat_name,
			creation_date,
			photo)
		VALUES (?,?,?,?,?)
	`, chatID.String(), "group", chatname, creationdate, photouri)

	if err != nil {
		return Chat{}, err
	}
	//PHOTOURI IS OPTIONAL, IS IT FRONT END RESPONSIBILITY TO INTERPRET ""
	//AS "no photo"?

	for _, memberID := range members {
		_, err = tx.Exec(`
			INSERT INTO chat_members (chat_id, user_id)
			VALUES (?,?)
		`, chatID.String(), memberID)

		if err != nil {
			return Chat{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return Chat{}, err //INTERNAL SERVER ERR
	}

	return Chat{
		ChatID:       chatID.String(),
		GroupOrChat:  "group",
		ChatName:     chatname,
		CreationDate: creationdate,
		Photo:        photouri,
	}, nil
}
