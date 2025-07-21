// Updates documents that match a query filter by using the Go driver
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Defines a Restaurant struct as a model for documents in the "restaurants" collection
type Restaurant struct {
	ID            bson.ObjectID `bson:"_id"`
	Name          string        `bson:"name"`
	Cuisine       string        `bson:"cuisine"`
	AverageRating float64       `bson:"avg_rating,omitempty"`
}

// Create a filter struct to specify the documents to update
type UpdateManyRestaurantFilter struct {
	Cuisine string `bson:"cuisine"`
	Borough string `bson:"borough"`
}

// Defines a RestaurantUpdate struct to specify the fields to update
type RestaurantUpdateMany struct {
	AverageRating float64 `bson:"avg_rating"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/connect/mongoclient/#environment-variable")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("sample_restaurants").Collection("restaurants")
	filter := UpdateManyRestaurantFilter{
		Cuisine: "Pizza",
		Borough: "Brooklyn",
	}

	// Creates instructions to update the values of the "AverageRating" field
	update := bson.D{{"$set", RestaurantUpdateMany{AverageRating: 4.5}}}

	// Updates documents in which the value of the "Cuisine"
	// field is "Pizza"
	result, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}

	// Prints the number of updated documents
	fmt.Printf("Documents updated: %v\n", result.ModifiedCount)

	// When you run this file for the first time, it should print output similar to the following:
	// Documents updated: 296
}
