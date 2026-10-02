package database

import (
	"errors"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/globaltime"
)

func (db *appdbimpl) MakeChatGroup(chatID string, chatname string, photouri string, members []string, isgroup bool) (Chat, error) {

	// make the group or chat
	creationdate := globaltime.Now()
	grouporchat := privatechat
	// true->it's a group
	if isgroup {
		grouporchat = groupchat
	}

	tx, err := db.c.Begin()
	if !(errors.Is(err, nil)) {
		return Chat{}, ErrTX
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
	`, chatID, grouporchat, chatname, creationdate, "storage/data"+photouri)

	if !(errors.Is(err, nil)) {
		return Chat{}, ErrExec
	}
	// private chats have a black pic as scaffold, but frontend shows user the other
	// user's pic; likewise, chatname is a scaffold for priv chats, which otherwise
	// show the other user's username visually

	for _, memberID := range members {
		_, err = tx.Exec(`
			INSERT INTO chat_members (chat_id, user_id)
			VALUES (?,?)
		`, chatID, memberID)

		if !(errors.Is(err, nil)) {
			return Chat{}, ErrExec
		}
	}

	if err = tx.Commit(); !(errors.Is(err, nil)) {
		return Chat{}, ErrTX // INTERNAL SERVER ERR
	}

	return Chat{
		ChatID:       chatID,
		GroupOrChat:  grouporchat,
		ChatName:     chatname,
		CreationDate: creationdate,
		Photo:        photouri,
	}, nil

}
