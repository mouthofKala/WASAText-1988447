package database

import (
	"database/sql"
	"errors"
	"time"
)

func (db *appdbimpl) AddToGroup(targetuserID string, groupID string, userID string) (Chat, error) {
	//1. check that userID is in group themselves
	var slimshady bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chat_members
			WHERE chat_id = ? AND user_id = ?)
	`, groupID, userID).Scan(&slimshady)
	if err != nil {
		return Chat{}, err
	}
	if !slimshady {
		return Chat{}, ErrUserNotFound //403
	}

	//2. checking the chat is actually a group
	var grouporchat string
	var chatname string
	var creationdate time.Time
	var photo string
	err = db.c.QueryRow(`
		SELECT group_or_chat, chat_name, creation_date, photo
		FROM chats
		WHERE chat_id = ?
	`, groupID).Scan(&grouporchat, &chatname, creationdate, photo)

	if errors.Is(err, sql.ErrNoRows) {
		return Chat{}, sql.ErrNoRows //bad request
	}
	if err != nil {
		return Chat{}, err //internal server error
	}
	if grouporchat != "group" {
		return Chat{}, ErrNotaGroup //403
	}

	//3. check that targetuserID actually exists
	var targetuserIDexists bool
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE user_id = ?
	`, targetuserID).Scan(&targetuserIDexists)

	if err != nil {
		return Chat{}, err //bad req?
	}

	//4. checking if user is ALREADY in the group
	var alrpresent bool
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chat_members
			WHERE chat_id = ? AND user_id = ?)
	`, groupID, targetuserID).Scan(&alrpresent)

	if alrpresent {
		return Chat{}, ErrUserAlreadyIn //bad request
	}

	//5. now we are sure we are adding a new, unique member to a unique gc
	_, err = db.c.Exec(`
		INSERT INTO chat_members (chat_id, user_id)
		VALUES (?,?)`, groupID, targetuserID)
	if err != nil {
		return Chat{}, ErrModifyingChat //ISE
	}

	var recentmsg MessagePreview
	var recentmsgptr *MessagePreview
	err = db.c.QueryRow(`
		SELECT
			message_id,
			sender_id,
			content,
			photo,
			timestamp
		FROM messages
		WHERE chat_id = ?
		ORDER BY timestamp DESC
		LIMIT 1
	`, groupID).Scan(
		&recentmsg.MessageID,
		&recentmsg.SenderID,
		&recentmsg.Content,
		&recentmsg.Photo,
		&recentmsg.Timestamp)

	if err == nil {
		recentmsgptr = &recentmsg
	} else if !errors.Is(err, sql.ErrNoRows) {
		return Chat{}, ErrFetchingChat
	}

	var members []string
	rows, err := db.c.Query(`
		SELECT user_id
		FROM chat_members
		WHERE chat_id = ?	
	`, groupID)
	if err != nil {
		return Chat{}, ErrFetchingMembers //internal server error
	}
	defer rows.Close()

	for rows.Next() {
		var member string

		err = rows.Scan(&member)
		if err != nil {
			return Chat{}, ErrFetchingMembers //internal server error
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return Chat{}, ErrFetchingMembers
	}
	return Chat{
		ChatID:        groupID,
		GroupOrChat:   "group",
		ChatName:      chatname,
		Members:       members,
		CreationDate:  creationdate,
		Photo:         photo,
		MostRecentMsg: recentmsgptr,
	}, nil
}
