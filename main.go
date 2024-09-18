package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gurix/sign_up_bot/handlers"
	"github.com/gurix/sign_up_bot/store"
)

func main() {
	// Initialize MongoDB
	store.InitMongoDB()

	// Setup Chi router
	request := chi.NewRouter()
	request.Use(middleware.Logger)
	request.Use(middleware.Recoverer)

	// Serve static files (index.html)
	request.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})

	// API endpoint to handle chat requests
	request.Post("/v1/chat", handlers.ChatInput)

	request.Get("/v1/chats", handlers.GetChats)

	// Fetch the port from the environment variable or use 8000 as a default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	const readWriteTimeout = 10 * time.Second

	const idleTimeout = 6 * readWriteTimeout

	const oneMB = 1 << 20

	// Start the server
	log.Printf("Server is running on http://localhost:%s\n", port)
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           request,
		ReadTimeout:       readWriteTimeout,
		WriteTimeout:      readWriteTimeout,
		IdleTimeout:       idleTimeout,
		ReadHeaderTimeout: readWriteTimeout,
		MaxHeaderBytes:    oneMB,
		TLSConfig:         nil,
		TLSNextProto:      nil,
		ErrorLog:          log.New(os.Stderr, "SERVER: ", log.LstdFlags),
		ConnState: func(_ net.Conn, state http.ConnState) {
			log.Printf("Connection state changed: %v", state) // Logging connection states
		},
		DisableGeneralOptionsHandler: true,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}
	log.Fatal(server.ListenAndServe())
}
