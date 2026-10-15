package database

import "errors"

func (db *appdbimpl) SearchUsers(searchkey string, userID string) ([]User, error) {
	usern := ""
	err := db.c.QueryRow(`
		SELECT username
		FROM users
		WHERE user_id = ?
	`).Scan(&usern)
	if err != nil {
		return nil, err
	}
	rows, err := db.c.Query(`
		SELECT user_id, username, photo
		FROM users
		WHERE username LIKE ?
		AND username !=  ?
	`, "%"+searchkey+"%", usern)
	if !errors.Is(err, nil) {
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
		if !errors.Is(err, nil) {
			return nil, err
		}
		userlist = append(userlist, user)
	}
	return userlist, rows.Err()
}
