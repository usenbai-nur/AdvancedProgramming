package infrastructure

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func InitDatabase() error {
	_ = godotenv.Load()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
		log.Printf("MONGODB_URI is not set, trying default %s", uri)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri).SetServerSelectionTimeout(2 * time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Mongo connect failed: %v. Continuing with in-memory order storage.", err)
		return nil
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("Mongo ping failed: %v. Continuing with in-memory order storage.", err)
		return nil
	}

	Client = client
	Database = client.Database("carstore")
	log.Println("Connected to MongoDB")
	return nil
}

func CloseDatabase() {
	if Client != nil {
		_ = Client.Disconnect(context.TODO())
		log.Println("🔌 MongoDB connection closed")
	}
}
