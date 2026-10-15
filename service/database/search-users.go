package database

import "errors"

func (db *appdbimpl) SearchUsers(searchkey string, userID string) ([]User, error) {

	rows, err := db.c.Query(`
		SELECT user_id, username, photo
		FROM users
		WHERE username LIKE ?
		AND user_id !=  ?
	`, "%"+searchkey+"%", userID)
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
