package database

func (db *appdbimpl) userExists(userID string) (bool, error) {
	var exists bool
	err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM users
			where user_id=?)
	`, userID).Scan(&exists)
	return exists, err
}
