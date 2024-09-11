package main

import (
  "context"
  "encoding/json" 
  "encoding/gob"
  "fmt"
  "net/http"
  "time"
  "os"

  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
  "github.com/gorilla/sessions"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/bson"
)

// ChatMessage represents a chat message structure
type ChatMessage struct {
  CreatedAt time.Time `json:"created_at"`
  Message string `json:"message"`
  Result string `json:"result"`
}

// ChatDialog represents a session structure to be stored in MongoDB
type ChatDialog struct {
  SessionID string   `bson:"session_id"`
  Messages  []ChatMessage `bson:"messages"`
}

var store = sessions.NewCookieStore([]byte("your-secret-key"))
var client *mongo.Client
var chatCollection *mongo.Collection

func init() {

  gob.Register(&[]ChatMessage{})
}

func main() {
  // MongoDB setup
  var err error
  client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
  if err != nil {
    panic(err)
  }
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  err = client.Connect(ctx)
  if err != nil {
    panic(err)
  }
  chatCollection = client.Database("chatbotdb").Collection("sessions")

  // Setup Chi router
  r := chi.NewRouter()
  r.Use(middleware.Logger)
  r.Use(middleware.Recoverer)

  // Serve static files (index.html)
  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "frontend/index.html")
  })

  // API endpoint to handle chat requests
  r.Post("/v1/chat", handleChat)

  // Fetch the port from the environment variable or use 8000 as a default
  port := os.Getenv("PORT")
  if port == "" {
    port = "8000"
  }

  // Start the server
  fmt.Printf("Server is running on http://localhost:%s\n", port)
  http.ListenAndServe(":"+port, r)
}

// handleChat handles incoming chat messages and stores them in the session and MongoDB
func handleChat(w http.ResponseWriter, r *http.Request) {
  // Initialize session
  session, err := store.Get(r, "chat-session")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  
  // Check if session is new and persist it if necessary
  if session.IsNew {
    session.Values["ID"] = fmt.Sprintf("%d", time.Now().UnixNano())
    session.Save(r, w)
  }
  
  // Parse incoming message
  var chatMessage ChatMessage
  if err := json.NewDecoder(r.Body).Decode(&chatMessage); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  
  // Simulate bot response
  chatMessage.Result = "You said: \"" + chatMessage.Message + "\""
  chatMessage.CreatedAt = time.Now()

  // Append messages to session
  messages, ok := session.Values["messages"].([]ChatMessage)
  if !ok {
    messages = []ChatMessage{}
  }
  messages = append(messages, chatMessage)
  session.Values["messages"] = messages

  // Save session
  if err := sessions.Save(r, w); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // Store conversation in MongoDB
  _, err = chatCollection.UpdateOne(
    context.TODO(),
    bson.M{"session_id": session.Values["ID"]},
    bson.M{
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
