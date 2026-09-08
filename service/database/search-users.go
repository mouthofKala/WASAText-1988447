package database

func (db *appdbimpl) SearchUsers(searchkey string) ([]User, error) {
	rows, err := db.c.Query(`
		SELECT user_id, username, photo
		FROM users
		WHERE username LIKE ?
	`, "%"+searchkey+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userlist []User
	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.Photo,
		)
		if err != nil {
			return nil, err
		}
		userlist = append(userlist, user)
	}
	return userlist, rows.Err()
}
