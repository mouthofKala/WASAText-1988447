package database

func (db *appdbimpl) SetGroupName(newname string, chatID string, userID string) error {
	var privacy string
	err := db.c.QueryRow(`
		SELECT group_or_chat
		FROM chats
		WHERE chat_id = ?
	`, chatID).Scan(&privacy)
	if err != nil {
		return err //500
	}
	if privacy != "group" {
		return ErrNotaGroup //403
	}

	//if requestor not a member 403
	var membership bool
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM chat_members
			WHERE chat_id = ?
			AND user_id = ?)`, chatID, userID).Scan(&membership)
	if err != nil {
		return err //500
	}
	if !membership {
		return ErrNotaGroupMember //403
	}

	_, err = db.c.Exec(`
		UPDATE chats
		SET chat_name = ?
		WHERE chat_id = ?
	`, newname, chatID)
	return err //500
}
