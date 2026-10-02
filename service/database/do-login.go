package database

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) DoLogin(username string) (string, error) {
	var userID string

	err := db.c.QueryRow(`
		SELECT user_id
		FROM users
		WHERE username = ?
	`, username).Scan(&userID)

	// user exists->authenticated
	if errors.Is(err, nil) {
		return userID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	// create a new user
	newID, err := uuid.NewV4()
	if !errors.Is(err, nil) {
		return "", err
	}
	userID = newID.String()

	_, err = db.c.Exec(`
		INSERT INTO users (user_id, username, photo)
		VALUES (?,?,?)`,
		userID, username, "photouriofblackpic.jpg")

	if !errors.Is(err, nil) {
		return "", err
	}

	return userID, nil

}
