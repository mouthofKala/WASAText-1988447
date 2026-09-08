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

	//user exists->authenticated
	if err == nil {
		return userID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	//create a new user
	newID, err2 := uuid.NewV4()
	if err2 != nil {
		return "", err2
	}
	userID = newID.String()

	_, err1 := db.c.Exec(`
		INSERT INTO users (user_id, username, photo)
		VALUES (?,?,?)`,
		userID, username, "photouriofblackpic.jpg")

	if err1 != nil {
		return "", err1
	}

	return userID, nil

}
