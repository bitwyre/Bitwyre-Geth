package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func SetupMongo() (*mongo.Client, error) {
	cwd, _ := os.Getwd()
	fmt.Println("Current working directory:", cwd)

	// Get username and password from environment variables
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	// Create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	connectionURI := fmt.Sprintf("mongodb+srv://%s:%s@gethcluster.bqax6oq.mongodb.net/?retryWrites=true&w=majority", username, password)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI).SetServerAPIOptions(serverAPI))
	if err != nil {
		log.Fatal(err)
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected and pinged.")

	return client, nil
}
