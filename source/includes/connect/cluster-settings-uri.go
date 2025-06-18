package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// start-uri-variable
// Connection string with cluster settings options
const (
	uri = "mongodb://localhost:27017/?serverSelectionTimeoutMS=10000&localThresholdMS=15"
)

// end-uri-variable

func main() {
	// start-apply-uri
	// Creates a new client and connects to the server
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	// end-apply-uri

	fmt.Println("Connected to MongoDB with cluster settings options")
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
}
