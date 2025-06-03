package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	//start-connection-pool-uri
	// Connection string with connection pool options
	uri := "mongodb://localhost:27017/?maxPoolSize=50&minPoolSize=10&maxIdleTimeMS=30000"

	// Creates a new client and connect to the server
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	//end-connection-pool-uri

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
}
