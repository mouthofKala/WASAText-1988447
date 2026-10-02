package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetUserProfile(userID string) (User, error) {

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

	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound // 404
	}

	return user, err // 500
}
