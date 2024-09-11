package sessions

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/gurix/sign_up_bot/models"
)

var store = sessions.NewCookieStore([]byte(SecretKey()))

func SecretKey() string {
	// Fetch the value of the environment variable "GEHEIMNIS"
	secret := os.Getenv("GEHEIMNIS")

	// Check if the environment variable is set
	if secret == "" {
		log.Fatalf("The environment variable 'GEHEIMNIS' is not set.")
	}

	return secret
}

func init() {
	gob.Register([]models.ChatMessage{}) // Register the slice type for session storage
}

// GetSession retrieves the session from the request
func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "chat-session")
}

// SaveSession saves the session data
func SaveSession(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return sessions.Save(r, w)
}

// AppendMessageToSession appends a message to the session
func AppendMessageToSession(session *sessions.Session, message models.ChatMessage) []models.ChatMessage {
	messages, ok := session.Values["messages"].([]models.ChatMessage)
	if !ok {
		messages = []models.ChatMessage{}
	}
	messages = append(messages, message)
	session.Values["messages"] = messages
	return messages
}

