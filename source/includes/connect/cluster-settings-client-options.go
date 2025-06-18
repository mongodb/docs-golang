package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}
	// start-client-options
	// Sets client options with cluster settings
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerSelectionTimeout(10 * time.Second).
		SetLocalThreshold(15 * time.Millisecond)

	// Creates a new client and connects to the server
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// end-client-options

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
}
