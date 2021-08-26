package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
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

	// begin insertMany
	coll := client.Database("insertDB").Collection("haikus")
	docs := []interface{}{
		bson.D{{"title", "Record of a Shriveled Datum"}, {"text", "No bytes, no problem. Just insert a document, in MongoDB"}},
		bson.D{{"title", "Showcasing a Blossoming Binary"}, {"text", "Binary data, safely stored with GridFS. Bucket the data"}},
	}

	result, err := coll.InsertMany(context.TODO(), docs)
	// end insertMany

	if err != nil {
		panic(err)
	}

	// When you run this file it should print "2 documents inserted with IDs: ...
	fmt.Printf("%d documents inserted with IDs:\n", len(result.InsertedIDs))
	for _, id := range result.InsertedIDs {
		fmt.Printf("\t%s\n", id)
	}
}
