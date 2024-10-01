package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gurix/sign_up_bot/llm"
	"github.com/gurix/sign_up_bot/models"
	"github.com/gurix/sign_up_bot/store"
)

// sendJSONResponse is a helper function to encode data to JSON and send it as a response.
func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getDialogID(request *http.Request) string {
	return request.Context().Value("dialog_id").(string)
}

func GetChats(writer http.ResponseWriter, request *http.Request) {
	// Fetch all messages of a session
	dialog, err := store.GetDialog(getDialogID(request))
	if err != nil {
		if err.Error() == "no dialog found" {
			http.NotFound(writer, request)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Send response back to the client
		sendJSONResponse(writer, dialog.Messages)
	}
}

func ChatInput(writer http.ResponseWriter, request *http.Request) {
	var ai llm.Ai

	llmClient := llm.OllamaConnect()
	ai.Client = llmClient

	// Parse incoming message
	var chatMessage models.ChatMessage
	if err := json.NewDecoder(request.Body).Decode(&chatMessage); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	dialog, err := store.GetDialog(getDialogID(request))
	if err != nil && err.Error() != "no dialog found" {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}
	// Retrieve the result from the model
	resp, err := ai.GenerateResponse(context.Background(), chatMessage.Message, dialog)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	// Just use the first choice for the moment
	chatMessage.Result = llm.GetFirstContent(resp)

	// Update chat collection
	if err := store.UpdateChatCollection(getDialogID(request), chatMessage); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}

	// Send response back to the client
	sendJSONResponse(writer, chatMessage)
}
