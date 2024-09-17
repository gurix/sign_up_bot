package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gurix/sign_up_bot/models"
	"github.com/gurix/sign_up_bot/sessions"
	"github.com/gurix/sign_up_bot/store"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetChats(w http.ResponseWriter, r *http.Request) {

	// Initialize session
	session := sessions.GetSession(w, r)

	messages := sessions.GetMessagesFromSession(session)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// HandleChat handles incoming chat messages and stores them in the session and MongoDB
func HandleChat(w http.ResponseWriter, r *http.Request) {
	// Initialize session
	session := sessions.GetSession(w, r)

	// Retrieve the session ID
	sessionID := session.ID

	// Parse incoming message
	var chatMessage models.ChatMessage
	if err := json.NewDecoder(r.Body).Decode(&chatMessage); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Simulate bot response
	chatMessage.Result = "You said: \"" + chatMessage.Message + "\""

	// Append messages to session
	sessions.AppendMessageToSession(session, chatMessage)

	// Save session
	if err := sessions.SaveSession(r, w, session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store or update conversation in MongoDB
	_, err := store.ChatCollection.UpdateOne(
		context.TODO(),
		bson.M{"session_id": sessionID},
		bson.M{
			"$set":  bson.M{"session_id": sessionID},
			"$push": bson.M{"messages": chatMessage}, // Use $push to append the new message
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatMessage)
}
