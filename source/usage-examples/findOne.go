// This program finds a document titled 'The Room' in your MongoDB deployment.
package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {

	// Replace the uri string with your MongoDB deployment's connection string.
	uri := "mongodb+srv://<user>:<password>@<cluster-url>?retryWrites=true&w=majority"

	// Create an empty context
	ctx := context.TODO()

	// Connect to your MongoDB deployment
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// Disconnect Client
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// Define Collection
	coll := client.Database("sample_mflix").Collection("movies")
	var result bson.M
	// Initialize an empty *FindOneoptions value
	opts := options.FindOne()
	// Project only the `title` and `imdb` fields in the returned document
	opts = opts.SetProjection(bson.M{"_id": 0, "title": 1, "imdb": 1})
	// Sort matched documents in descending order by rating
	opts = opts.SetSort(bson.D{{"rating", -1}})
	// Retrieve your document using a query and your configuration options
	err = coll.FindOne(ctx, bson.M{"title": "The Room"}, opts).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	// Print your matched document
	fmt.Printf("%v\n", result)
}
