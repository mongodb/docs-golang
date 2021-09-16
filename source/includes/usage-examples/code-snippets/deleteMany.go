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
		log.Fatal("You must set your `MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/")
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

	// begin deleteMany
	coll := client.Database("sample_mflix").Collection("movies")
	filter := bson.D{{"runtime", bson.D{{"$gt", 800}}}}

	results, err := coll.DeleteMany(context.TODO(), filter)
	// end deleteMany

	if err != nil {
		panic(err)
	}

	// When you run this file for the first time, it should print:
	// Documents deleted: 4
	fmt.Printf("Documents deleted: %d\n", results.DeletedCount)
}
