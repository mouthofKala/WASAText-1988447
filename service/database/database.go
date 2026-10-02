/*
To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.
For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

const DefaultGroupPhoto = "storage/data/pfp/black.jpg"

var ErrUserNotFound = errors.New("user not found")
var ErrUsernameUnavailable = errors.New("username already in use")
var ErrChatAlreadyExists = errors.New("chat already exists")
var ErrNoMembersSelected = errors.New("no members selected for chat or group")
var ErrGeneratingID = errors.New("error generating chat ID")
var ErrUserAlreadyIn = errors.New("error is already in the group")
var ErrNotaGroup = errors.New("not a group")
var ErrFetchingMembers = errors.New("error fetching members of group")
var ErrFetchingChat = errors.New("error fetching chat")
var ErrModifyingChat = errors.New("error modifying chat")
var ErrInvalidImage = errors.New("image not valid or unavailable")
var ErrInvalidMsg = errors.New("message empty or invalid")
var ErrNotaGroupMember = errors.New("user is not a member of the group")
var ErrFetchingMessages = errors.New("error searching or scanning messages from db")
var ErrDecodingReq = errors.New("error decoding request")
var ErrExec = errors.New("error doing database Exec")
var ErrForbidden = errors.New("user not authorised to do this action")
var ErrBadReq = errors.New("resource does not exist")
var ErrStorage = errors.New("error with deletion or addition to storage")
var ISE = "500 Internal Server Error"
var BR = "400 Bad Request"
var Un = "401 Unauthorized"
var NF = "404 Not Found"
var C = "409 Conflict"
var F = "403 Forbidden"
var ErrTX = errors.New("error performing a transaction")
var ErrTXcm = errors.New("error deleting chat members during transaction")
var ErrTXc = errors.New("error deleting chat during transaction")
var ErrTXm = errors.New("error deleting messages during transaction")
var ErrTXnr = errors.New("error deleting not read msgs during transaction")
var ErrTXr = errors.New("error deleting reactions during transaction")

// DATAVASE METHODS
type AppDatabase interface {
	GetUserProfile(userID string) (User, error)
	SearchUsers(searchkey string) ([]User, error)
	SetMyUsername(targetusername string, userID string) error
	SetMyPhoto(photouri string, userID string) (string, error)
	DoLogin(username string) (string, error)
	GetMyConversations(userID string) ([]Chat, error)
	MakeChatGroup(chatID string, chatname string, photouri string, members []string, isgroup bool) (Chat, error)
	// MakeChat(otheruser string, creatorID string) (Chat, error)
	AddToGroup(targetuserIDs []string, groupID string, userID string) (Chat, error)
	SetGroupPhoto(photouri string, chatID string, userID string) (string, error)
	SetGroupName(newname string, chatID string, userID string) error
	GetConversation(userID string, chatID string) (Conversation, error)
	SendMessage(userID, chat_id string, content *string, photo *string, replyto *string, fwdfrom *string, messageID string) (Message, error)
	DeleteMessage(msgID string, userID string, chatID string) error
	ForwardMessage(userID, targetchat_id string, fwdmsgID string, chatID string, newmsgID string) (Message, error)
	CommentMessage(msgID string, chatID string, userID string, reactionID string, emoji rune) (Reaction, error)
	UncommentMessage(msgID string, chatID string, userID string, reactionID string) error
	LeaveGroup(userID string, chatID string) (bool, error)
	UserExists(userID string) (bool, error)
	ValidateChatGroup(members []string, isgroup bool) error
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required")
	}

	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            user_id TEXT PRIMARY KEY,
            username TEXT NOT NULL UNIQUE,
            photo TEXT
        );

        CREATE TABLE IF NOT EXISTS chats (
            chat_id TEXT PRIMARY KEY,
            group_or_chat BIT NOT NULL,
            chat_name TEXT NOT NULL,
            creation_date DATETIME NOT NULL,
            photo TEXT
        );

        CREATE TABLE IF NOT EXISTS chat_members (
            chat_id TEXT NOT NULL,
            user_id TEXT NOT NULL,
            PRIMARY KEY (chat_id, user_id)
        );

        CREATE TABLE IF NOT EXISTS messages (
            message_id TEXT PRIMARY KEY,
            chat_id TEXT NOT NULL,
            sender_id TEXT NOT NULL,
            content TEXT,
            photo TEXT,
            status TEXT NOT NULL,
            timestamp DATETIME NOT NULL,
            reply_to TEXT,
            fwd_from TEXT
			CHECK (content IS NOT NULL OR photo IS NOT NULL)
        );

		CREATE TABLE IF NOT EXISTS msg_notreads(
			message_id TEXT NOT NULL,
			chat_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			PRIMARY KEY (message_id, user_id)
		);

        CREATE TABLE IF NOT EXISTS reactions (
            reaction_id TEXT PRIMARY KEY,
            message_id TEXT NOT NULL,
			chat_id TEXT NOT NULL,
            user_id TEXT NOT NULL,
            emoji TEXT NOT NULL
        );
    `)

	if err != nil {
		return nil, fmt.Errorf("creating database tables: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
