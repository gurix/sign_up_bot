package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

// Export the store if you need it in main.go.
var Store = sessions.NewCookieStore([]byte(SecretKey()))

func SecretKey() string {
	// Fetch the value of the environment variable "GEHEIMNIS"
	secret := os.Getenv("GEHEIMNIS")

	// Check if the environment variable is set
	if secret == "" {
		log.Fatalf("The environment variable 'GEHEIMNIS' is not set.")
	}

	return secret
}

// GenerateRandomID creates a random unique identifier.
func GenerateRandomID() (string, error) {
	bytes := make([]byte, 16) // 16 bytes = 128 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// DialogIDMiddleware is a chi middleware that retrieves or creates a dialog ID.
func DialogIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the session, or create a new one
		session, err := Store.Get(r, "dialog-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the session has the dialog_id
		dialogID, ok := session.Values["dialog_id"]
		if !ok {
			// If not, generate a new dialog ID
			newDialogID, err := GenerateRandomID()
			if err != nil {
				http.Error(w, "Failed to generate dialog ID", http.StatusInternalServerError)
				return
			}

			// Store the dialog ID in the session
			session.Values["dialog_id"] = newDialogID
			dialogID = newDialogID

			// Save the session to persist the dialog ID
			if err := session.Save(r, w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Attach the dialog ID to the request context
		ctx := context.WithValue(r.Context(), "dialog_id", dialogID)
		next.ServeHTTP(w, r.WithContext(ctx)) // Pass the request to the next handler
	})
}
