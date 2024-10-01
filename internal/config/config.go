package config

import (
	"context"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Printf("Failed to create MongoDB client: %v", err)
		return nil, err
	}

	// Ping the primary
	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("Failed to ping MongoDB Atlas: %v", err)
		return nil, err
	}

	log.Println("Connected to MongoDB Atlas successfully!")

	db := client.Database(dbName)

	return &MongoDB{
		Client:   client,
		Database: db,
	}, nil
}

// Disconnect closes the connection to the database
func (m *MongoDB) Disconnect(ctx context.Context) error {
	if err := m.Client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect from MongoDB Atlas: %v", err)
		return err
	}
	log.Println("Disconnected from MongoDB Atlas successfully")
	return nil
}