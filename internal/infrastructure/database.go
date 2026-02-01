package infrastructure

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func InitDatabase() error {
	_ = godotenv.Load() // Пытаемся загрузить .env

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return fmt.Errorf("MONGODB_URI not found in environment or .env file")
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return err
	}

	Client = client
	Database = client.Database("carstore")
	log.Println("✅ Connected to MongoDB Atlas")
	return nil
}

func CloseDatabase() {
	if Client != nil {
		Client.Disconnect(context.TODO())
		log.Println("🔌 MongoDB connection closed")
	}
}
