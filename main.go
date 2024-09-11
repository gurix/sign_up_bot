package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gurix/sign_up_bot/handlers"
	"github.com/gurix/sign_up_bot/store"
)

func main() {
	// Initialize MongoDB
	store.InitMongoDB()

	// Setup Chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static files (index.html)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})

	// API endpoint to handle chat requests
	r.Post("/v1/chat", handlers.HandleChat)

	// Fetch the port from the environment variable or use 8000 as a default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Start the server
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, r)
}
