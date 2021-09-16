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

	// begin updatemany
	myCollection := client.Database("sample_airbnb").Collection("listingsAndReviews")
	filter := bson.D{{"address.market", "Sydney"}}
	update := bson.D{{"$mul", bson.D{{"price", 1.15}}}}

	result, err := myCollection.UpdateMany(context.TODO(), filter, update)
	// end updatemany

	if err != nil {
		panic(err)
	}

	// When you run this file for the first time, it should print:
	// Number of documents replaced: 609
	fmt.Printf("Documents updated: %v\n", result.ModifiedCount)
}
