package infrastructure

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func InitDatabase() error {
	uri := "mongodb+srv://norvi:nurik777@cluster0.y5nglmd.mongodb.net/carstore?retryWrites=true&w=majority"

	clientOptions := options.Client().ApplyURI(uri)
	var err error
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return fmt.Errorf("Atlas failed: %w", err)
	}

	if err = Client.Ping(context.TODO(), nil); err != nil {
		return fmt.Errorf("Atlas ping failed: %w", err)
	}

	Database = Client.Database("carstore")
	log.Println("MongoDB Atlas connected! norvi@cluster0")
	return nil
}

func CloseDatabase() {
	if Client != nil {
		Client.Disconnect(context.TODO())
	}
}
