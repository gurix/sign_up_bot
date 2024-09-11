package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client         *mongo.Client
	ChatCollection *mongo.Collection
)

// InitMongoDB initializes the MongoDB client and connects to the database
func InitMongoDB() {
	var err error
	Client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	ChatCollection = Client.Database("chatbotdb").Collection("sessions")
}
