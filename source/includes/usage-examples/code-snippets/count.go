package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// begin countDocuments
	coll := client.Database("sample_mflix").Collection("movies")
	filter := bson.D{{"countries", "China"}}

	estCount, estCountErr := coll.EstimatedDocumentCount(context.TODO())
	count, countErr := coll.CountDocuments(context.TODO(), filter)
	// end countDocuments

	if estCountErr != nil {
		panic(estCountErr)
	}
	if countErr != nil {
		panic(countErr)
	}

	// When you run this file, it should print:
	// Estimated number of documents in the movies collection: 23541
	// Number of movies from China: 303
	fmt.Printf("Estimated number of documents in the movies collection: %d\n", estCount)
	fmt.Printf("Number of movies from China: %d\n", count)
}
