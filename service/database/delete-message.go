package database

import "errors"

func (db *appdbimpl) DeleteMessage(msgID string, userID string, chatID string) error {
	// does msg exist?
	var exist bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM messages
			WHERE message_id = ? AND chat_id = ?)	
	`, msgID, chatID).Scan(&exist)
	if !errors.Is(err, nil) {
		return err // 500
	}
	if !exist {
		return ErrBadReq // 400
	}
	// 400

	// select msg by msgid and userid
	var auth bool
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM messages
			WHERE message_id = ? AND sender_id = ? AND chat_id = ?)	
	`, msgID, userID, chatID).Scan(&auth)
	if !errors.Is(err, nil) {
		return err // 500
	}
	if !auth {
		return ErrForbidden // 403
	}

	// begin transaction
	tx, err := db.c.Begin()
	if !errors.Is(err, nil) {
		return ErrTX
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.Exec(`
		DELETE FROM messages
		WHERE message_id = ? AND chat_id = ?
	`, msgID, chatID)
	if !errors.Is(err, nil) {
		return ErrExec
	}
	_, err = tx.Exec(`
		DELETE FROM msg_notreads
		WHERE message_id = ?
	`, msgID)
	if !errors.Is(err, nil) {
		return ErrExec
	}
	if err = tx.Commit(); !errors.Is(err, nil) {
		return ErrTX
	}

	return nil

}
