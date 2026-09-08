package database

func (db *appdbimpl) UncommentMessage(msgID string, chatID string, userID string, reactionID string) error {

	var chatexist bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chats
			WHERE chat_id = ?)
	`, chatID).Scan(&chatexist)

	if err != nil {
		return err //500
	}

	var msgexist bool
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM messages
			WHERE message_id = ?)
	`, msgID).Scan(&msgexist)

	if err != nil {
		return err //500
	}

	if !chatexist || !msgexist {
		return ErrBadReq
	}

	//is userid owner of the reaction? 400
	var isowner bool
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM reactions
			WHERE reaction_id = ? AND user_id = ?)`,
		reactionID, userID).Scan(&isowner)
	if err != nil {
		return err
	}
	if !isowner {
		return ErrForbidden
	}

	//remove from db
	_, err = db.c.Exec(`
		DELETE FROM reactions
		WHERE reaction_id = ?`, reactionID)

	if err != nil {
		return ErrExec
	}
	return nil
}
