package database

import "time"

type User struct {
	UserID   string
	Username string
	Photo    string
}

type Chat struct {
	ChatID        string          `json:"chatid"`
	GroupOrChat   bool            `json:"grouporchat"`
	ChatName      string          `json:"chatname"`
	Members       []string        `json:"members"`
	CreationDate  time.Time       `json:"creationdate"`
	Photo         string          `json:"photo"`
	MostRecentMsg *MessagePreview `json:"mostrecentmsg"`
}

type Message struct {
	MessageID string    `json:"messageid"`
	ChatID    string    `json:"chatid"`
	SenderID  string    `json:"userid"`
	SenderN   string    `json:"username"`
	Content   *string   `json:"content"`  // pointer to text, can distinguish NULL and ""
	Photo     *string   `json:"photouri"` // link in memory to img
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	ReplyTo   *string   `json:"replyto"` // ptr to msgID
	FwdFrom   *string   `json:"fwdfrom"` // ptr2 chatID
} // basically use ptrs when info is optional

type Conversation struct {
	Chat     Chat      `json:"chat"`
	Messages []Message `json:"messages"`
}

type MessagePreview struct {
	MessageID string    `json:"messageid"`
	SenderID  string    `json:"userid"`
	SenderN   string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
	Content   *string   `json:"content"`
	Photo     *string   `json:"photo"`
}

type Reaction struct {
	ReactionID string `json:"reactionid"`
	MessageID  string `json:"messageid"`
	ChatID     string `json:"chatid"`
	UserID     string `json:"userid"`
	Emoji      string `json:"emoji"`
}
