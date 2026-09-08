package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) SetMyUsername(targetusername string, userID string) error {
	var usernameExists string
	err := db.c.QueryRow(`
		SELECT username
		FROM users
		WHERE username = ?
	`, targetusername).Scan(&usernameExists)

	if err == nil {
		return ErrUsernameUnavailable
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err1 := db.c.Exec(`
		UPDATE users
		SET username = ?
		WHERE user_id = ?
	`, targetusername, userID)

	return err1
}
