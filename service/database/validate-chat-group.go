package database

import (
	"database/sql"
	"errors"
)

const (
	privatechat = false
	groupchat   = true
)

func (db *appdbimpl) ValidateChatGroup(members []string, isgroup bool) error {
	if len(members) == 0 {
		return ErrNoMembersSelected
	}
	if len(members) != 2 && !isgroup {
		return ErrBadReq
	}

	if len(members)+1 > 100 {
		return ErrTooMany
	}
	for _, memberID := range members {
		exists, err := db.UserExists(memberID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrUserNotFound
		}

	}
	// if private, check if a private chat with those members exists
	if !isgroup {
		var existingChatID string
		err := db.c.QueryRow(`
			SELECT cm1.chat_id
			FROM chat_members cm1
			JOIN chat_members cm2
				ON cm1.chat_id = cm2.chat_id
			WHERE cm1.user_id = ?
				AND cm2.user_id = ?
				AND (
					SELECT COUNT(*)
					FROM chat_members cm3
					WHERE cm3.chat_id = cm1.chat_id	
				) = 2
			LIMIT 1
		`, members[0], members[1]).Scan(&existingChatID)

		if err == nil {
			return ErrChatAlreadyExists // 409
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if !(errors.Is(err, nil)) {
			return err // 500
		}

	}
	return nil

}
