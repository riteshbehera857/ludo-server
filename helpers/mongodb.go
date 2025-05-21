package helpers

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
	clientError    error
)

// InitializeNewMongoClient initializes the MongoDB client singleton
func InitializeNewMongoClient(uri string) (*mongo.Client, error) {
	clientOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientOptions := options.Client().ApplyURI(uri)
		clientInstance, clientError = mongo.Connect(ctx, clientOptions)
		if clientError != nil {
			log.Printf("Failed to connect to MongoDB: %v", clientError)
			return
		}

		// Ping the database
		if err := clientInstance.Ping(ctx, nil); err != nil {
			clientError = err
			log.Printf("Failed to ping MongoDB: %v", err)
			return
		}

		log.Println("Successfully connected to MongoDB")
	})

	return clientInstance, clientError
}

// GetMongoClient returns the MongoDB client singleton instance
func GetMongoClient() *mongo.Client {
	if clientInstance == nil {
		log.Fatal("MongoDB client not initialized. Call InitializeNewMongoClient first")
	}
	return clientInstance
}
