package store

import (
	"context"
	"time"

	"github.com/gurix/sign_up_bot/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client         *mongo.Client
	ChatCollection *mongo.Collection
)

const timeout = 10 * time.Second

// InitMongoDB initializes the MongoDB client and connects to the database.
func InitMongoDB() {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	Client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	ChatCollection = Client.Database("chatbotdb").Collection("sessions")
}

// UpdateChatCollection updates or inserts chat messages in the MongoDB.
func UpdateChatCollection(sessionID string, chatMessage models.ChatMessage) error {
	_, err := ChatCollection.UpdateOne(
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

func GetDialog(sessionID string) (models.ChatDialog, error) {
	var result models.ChatDialog
	err := ChatCollection.FindOne(
		context.TODO(),
		bson.M{"session_id": sessionID},
	).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.ChatDialog{}, nil
		}

		return models.ChatDialog{}, err // Return any other errors
	}

	return result, nil
}
