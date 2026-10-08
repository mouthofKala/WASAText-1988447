package database

import (
	"database/sql"
	"errors"
	"time"
)

func (db *appdbimpl) AddToGroup(targetuserIDs []string, groupID string, userID string) (Chat, error) {
	// 1. check that userID is in group themselves
	var slimshady bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chat_members
			WHERE chat_id = ? AND user_id = ?)
	`, groupID, userID).Scan(&slimshady)
	if !errors.Is(err, nil) {
		return Chat{}, err
	}
	if !slimshady {
		return Chat{}, ErrUserNotFound // 403
	}

	// 2. checking the chat is actually a group
	var grouporchat string
	var chatname string
	var creationdate time.Time
	var photo string
	err = db.c.QueryRow(`
		SELECT group_or_chat, chat_name, creation_date, photo
		FROM chats
		WHERE chat_id = ?
	`, groupID).Scan(&grouporchat, &chatname, &creationdate, &photo)

	if errors.Is(err, sql.ErrNoRows) {
		return Chat{}, sql.ErrNoRows // bad request
	}
	if !errors.Is(err, nil) {
		return Chat{}, err // internal server error
	}
	if grouporchat != group {
		return Chat{}, ErrNotaGroup // 403
	}

	// ADDENDUM: CHECKING THAT THE AMT SELECTED IS NOT BIGGER THAN 100
	// 1.a fetch all memberships and count them
	// 2.a check that len(targetuserIDs)<100 - amt of membershipt
	membercount := 0
	err = db.c.QueryRow(`
		SELECT COUNT(*)
		FROM chat_members
		WHERE chat_id = ?
		GROUP BY user_id
	`, groupID).Scan(&membercount)
	if err != nil {
		return Chat{}, ErrFetchingMembers // 500
	}
	if len(targetuserIDs)+membercount > 100 {
		return Chat{}, ErrBadReq // 400 cuz adding too much
	}

	// 3. check that targetuserID actually exists
	for _, targetuserID := range targetuserIDs {
		var targetuserIDexists bool
		err = db.c.QueryRow(`
			SELECT EXISTS(
				SELECT 1
				FROM users
				WHERE user_id = ?
			)
		`, targetuserID).Scan(&targetuserIDexists)
		if !errors.Is(err, nil) {
			return Chat{}, err // bad req?
		}
		if !targetuserIDexists {
			return Chat{}, ErrUserNotFound // 400
		}
		// 4. checking if user is ALREADY in the group
		var alrpresent bool
		err = db.c.QueryRow(`
			SELECT EXISTS(
				SELECT 1
				FROM chat_members
				WHERE chat_id = ? AND user_id = ?)
		`, groupID, targetuserID).Scan(&alrpresent)

		if err != nil {
			return Chat{}, err
		}
		if alrpresent {
			return Chat{}, ErrUserAlreadyIn // bad request
		}
	}

	// 5. now we are sure we are adding (a) new, unique member(s) to a unique gc
	// TRANSACTION

	tx, err := db.c.Begin()
	if err != nil {
		return Chat{}, ErrTX
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, targetuserID := range targetuserIDs {
		_, err = tx.Exec(`
			INSERT INTO chat_members (chat_id, user_id)
			VALUES (?,?)`, groupID, targetuserID)
		if !errors.Is(err, nil) {
			return Chat{}, ErrExec // ISE
		}
	}

	if err = tx.Commit(); err != nil {
		return Chat{}, ErrTX
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

	if errors.Is(err, nil) {
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
	if !errors.Is(err, nil) {
		return Chat{}, ErrFetchingMembers // internal server error
	}
	defer rows.Close()

	for rows.Next() {
		var member string

		err = rows.Scan(&member)
		if !errors.Is(err, nil) {
			return Chat{}, ErrFetchingMembers // internal server error
		}
		members = append(members, member)
	}
	if err := rows.Err(); !errors.Is(err, nil) {
		return Chat{}, ErrFetchingMembers
	}
	return Chat{
		ChatID:        groupID,
		GroupOrChat:   true,
		ChatName:      chatname,
		Members:       members,
		CreationDate:  creationdate,
		Photo:         photo,
		MostRecentMsg: recentmsgptr,
	}, nil
}
