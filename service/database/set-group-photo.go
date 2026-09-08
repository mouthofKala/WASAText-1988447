package database

func (db *appdbimpl) SetGroupPhoto(photouri string, chatID string, userID string) (string, error) {
	var privacy string
	var oldPhoto string
	err := db.c.QueryRow(`
		SELECT group_or_chat,photo
		FROM chats
		WHERE chat_id = ?
	`, chatID).Scan(&privacy, &oldPhoto)
	if err != nil {
		return "", err //500
	}
	if privacy != "group" {
		return "", ErrNotaGroup //403
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
		return "", err //500
	}
	if !membership {
		return "", ErrNotaGroupMember //403
	}

	_, err = db.c.Exec(`
		UPDATE chats
		SET photo = ?
		WHERE chat_id = ?`, photouri, chatID)

	if err != nil {
		return "", err //500
	}
	return oldPhoto, nil
}
