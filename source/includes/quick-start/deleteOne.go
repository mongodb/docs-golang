package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace the uri string with your MongoDB deployment's connection string.
const uri = "mongodb+srv://<username>:<password>@<cluster-address>/test?w=majority"
	
func main() {

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// begin deleteOne
	coll := client.Database("sample_mflix").Collection("movies")
	result, err := coll.DeleteOne(ctx, bson.D{{"title", "Twilight"}})
	// end deleteOne

	if err != nil {
		log.Panic(err)
	}

	// When you run this file for the first time, it should print "Number of documents deleted: 1"
	fmt.Printf("Number of documents deleted: %d\n", result)
}
