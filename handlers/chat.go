package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/gurix/sign_up_bot/models"
	s "github.com/gurix/sign_up_bot/sessions"
	"github.com/gurix/sign_up_bot/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// sendJSONResponse is a helper function to encode data to JSON and send it as a response.
func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetChats(writer http.ResponseWriter, request *http.Request) {
	// Initialize session
	session := initializeSession(writer, request)

	// Fetch all messages of a session
	messages := s.GetMessagesFromSession(session)

	// Send response back to the client
	sendJSONResponse(writer, messages)
}

// initializeSession initializes the session and returns it.
func initializeSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session := s.GetSession(w, r)
	w.Header().Set("Content-Type", "application/json")

	return session
}

// updateChatCollection updates or inserts chat messages in the MongoDB.
func updateChatCollection(sessionID string, chatMessage models.ChatMessage) error {
	_, err := store.ChatCollection.UpdateOne(
		context.TODO(),
		bson.M{"session_id": sessionID},
		bson.M{
			"$set":  bson.M{"session_id": sessionID},
			"$push": bson.M{"messages": chatMessage},
		},
		options.Update().SetUpsert(true),
	)

	return err
}

func ChatInput(writer http.ResponseWriter, request *http.Request) {
	// Initialize session
	session := initializeSession(writer, request)

	// Retrieve the session ID
	sessionID := session.ID

	// Parse incoming message
	var chatMessage models.ChatMessage
	if err := json.NewDecoder(request.Body).Decode(&chatMessage); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	// Simulate bot response
	chatMessage.Result = "You said: \"" + chatMessage.Message + "\""

	// Append messages to session
	s.AppendMessageToSession(session, chatMessage)

	// Save session
	if err := s.SaveSession(request, writer); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}

	// Update chat collection in MongoDB
	if err := updateChatCollection(sessionID, chatMessage); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}

	// Send response back to the client
	sendJSONResponse(writer, chatMessage)
}
