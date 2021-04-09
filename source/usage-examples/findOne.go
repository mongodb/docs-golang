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
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// Disconnect the client once the function returns.
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("sample_mflix").Collection("movies")
	var result bson.M
	opts := options.FindOne()
	opts.SetProjection(bson.M{"_id": 0, "title": 1, "imdb": 1})
	opts.SetSort(bson.D{{"imdb.rating", -1}})
	err = coll.FindOne(ctx, bson.M{"title": "The Room"}, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return
		}
		log.Fatal(err)
	}
	fmt.Println(result)
}
