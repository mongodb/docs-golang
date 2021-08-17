package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace the uri string with your MongoDB deployment's connection string.
const uri = "mongodb+srv://<username>:<password>@<cluster-address>/test?w=majority"

func main() {

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

	collCount, collCountErr := coll.EstimatedDocumentCount(context.TODO())
	count, countErr := coll.CountDocuments(context.TODO(), filter)
	// end countDocuments

	if collCountErr != nil {
		panic(collCountErr)
	}
	if countErr != nil {
		panic(countErr)
	}

	// When you run this file, it should print:
	// Estimated number of documents in the movies collection:: 23541
	// Number of movies from China: 303
	
	fmt.Printf("Estimated number of documents in the movies collection:: %d\n", collCount)
	fmt.Printf("Number of movies from China: %d\n", count)
}
