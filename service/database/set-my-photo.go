package database

func (db *appdbimpl) SetMyPhoto(photouri string, userID string) (string, error) {

	//get old pfp
	var oldPhoto string
	err := db.c.QueryRow(`
		SLEECT photo
		FROM users
		WHERE user_id = ?	
	`, userID).Scan(&oldPhoto)
	if err != nil {
		return "", err //500
	}

	_, err = db.c.Exec(`
		UPDATE users
		SET photo = ?
		WHERE user_id = ?
	`, photouri, userID)

	if err != nil {
		return "", err //500
	}
	return oldPhoto, nil

}
