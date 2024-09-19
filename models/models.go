package models

import "time"

// ChatMessage represents a chat message structure.
type ChatMessage struct {
	CreatedAt time.Time `json:"createdAt"`
	Message   string    `json:"message"`
	Result    string    `json:"result"`
}

// ChatDialog represents a session structure to be stored in MongoDB.
type ChatDialog struct {
	SessionID string        `bson:"session_id"`
	Messages  []ChatMessage `bson:"messages"`
}
