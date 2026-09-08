package database

import "database/sql"

func (db *appdbimpl) LeaveGroup(userID string, chatID string) (bool, error) {
	var chatIDexist bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chats
			WHERE chat_id = ?)`,
		chatID).Scan(&chatIDexist)
	if err != nil {
		return false, err //500
	}
	if !chatIDexist {
		return false, ErrBadReq //400
	}

	var membership bool
	err = db.c.QueryRow(`
	SELECT EXISTS(
		SELECT 1
		FROM chat_members
		WHERE chat_id = ? AND user_id = ?)`,
		chatID, userID).Scan(&membership)
	if err != nil {
		return false, err //500
	}
	if !membership {
		return false, ErrNotaGroupMember //403
	}

	//is userID the last member?
	rows, err := db.c.Query(`
		SELECT user_id
		FROM chat_members
		WHERE chat_id = ?
	`, chatID)

	if err == sql.ErrNoRows {
		return false, ErrNoMembersSelected //500
	}
	if err != nil {
		return false, err //500
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var member string

		err = rows.Scan(&member)
		if err != nil {
			return false, ErrFetchingMembers //internal server error
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return false, ErrFetchingMembers //500
	}

	if len(members) > 1 {
		_, err = db.c.Exec(`
			DELETE FROM chat_members
			WHERE user_id = ? and chat_id = ?
		`, userID, chatID)
		if err != nil {
			return false, err //500
		}

		return false, nil
	}

	if len(members) == 1 {
		//begin transaction to remove member and group
		tx, err := db.c.Begin()
		if err != nil {
			return true, ErrTX
		}
		defer tx.Rollback()

		_, err = tx.Exec(`
			DELETE FROM chat_members
			WHERE user_id = ? and chat_id = ?		
		`, userID, chatID)

		if err != nil {
			return true, ErrTXcm
		}

		_, err = tx.Exec(`
			DELETE FROM chats
			WHERE chat_id = ?		
		`, chatID)
		if err != nil {
			return true, ErrTXc
		}

		_, err = tx.Exec(`
			DELETE FROM messages
			WHERE chat_id = ?
		`, chatID)
		if err != nil {
			return true, ErrTXm
		}

		_, err = tx.Exec(`
			DELETE FROM msg_notreads
			WHERE chat_id = ?
		`, chatID)
		if err != nil {
			return true, ErrTXnr
		}

		_, err = tx.Exec(`
			DELETE FROM reactions
			WHERE chat_id = ?
		`, chatID)
		if err != nil {
			return true, ErrTXr
		}

		if err := tx.Commit(); err != nil {
			return true, ErrTX
		}

		return true, nil
	}

	//in whatsapp,it is possible to keep the group you have left
	//but that's likely because server has deleted chat whereas
	//you keep a copy of it locally
	return false, nil
}
