package database

import "time"

type User struct {
	UserID   string
	Username string
	Photo    string
}

type Chat struct {
	ChatID        string
	GroupOrChat   bool
	ChatName      string
	Members       []string
	CreationDate  time.Time
	Photo         string
	MostRecentMsg *MessagePreview
}

type Message struct {
	MessageID string
	ChatID    string
	SenderID  string
	Content   *string // pointer to text, can distinguish NULL and ""
	Photo     *string // pointer to img
	Status    string
	Timestamp time.Time
	ReplyTo   *string // ptr to msgID
	FwdFrom   *string // ptr2 chatID
} // basically use ptrs when info is optional

type Conversation struct {
	Chat     Chat
	Messages []Message
}

type MessagePreview struct {
	MessageID string
	SenderID  string
	Timestamp time.Time
	Content   *string
	Photo     *string
}

type Reaction struct {
	ReactionID string
	MessageID  string
	ChatID     string
	UserID     string
	Emoji      string
}
