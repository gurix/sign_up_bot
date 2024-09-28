package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gorilla/sessions"
)

// GenerateRandomID creates a random unique identifier.
func GenerateRandomID() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 128 bits

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// DialogIDMiddleware takes the sessionHandling store as a parameter.
func DialogIDMiddleware(sessionHandling *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			// Retrieve the session, or create a new one
			session, err := sessionHandling.Get(request, "dialog-session")
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			// Check if the session has the dialog_id
			dialogID, ok := session.Values["dialog_id"]
			if !ok {
				// If not, generate a new dialog ID
				newDialogID, err := GenerateRandomID()
				if err != nil {
					http.Error(writer, "Failed to generate dialog ID", http.StatusInternalServerError)
					return
				}

				// Store the dialog ID in the session
				session.Values["dialog_id"] = newDialogID
				dialogID = newDialogID

				// Save the session to persist the dialog ID
				if err := session.Save(request, writer); err != nil {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			// Attach the dialog ID to the request context
			ctx := context.WithValue(request.Context(), "dialog_id", dialogID)
			next.ServeHTTP(writer, request.WithContext(ctx)) // Pass the request to the next handler
		})
	}
}
