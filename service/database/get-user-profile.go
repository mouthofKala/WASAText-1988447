package database

func (db *appdbimpl) GetUserProfile(userID string) (User, error) {

	//CHECK IF USERID IS IN LIST OF USERS, GIVE BAD REQUEST ERR

	var user User
	err := db.c.QueryRow(`
		SELECT user_id, username, photo
		FROM users
		WHERE user_id=?
	`, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.Photo,
	)

	//IF NO ROWS, GIVE OUT ERROR NO ROWS FOR 404
	return user, err
}
