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

	// begin replace
	coll := client.Database("sample_mflix").Collection("movies")
	filter := bson.D{{"title", "Shrek"}}
	replacement := bson.D{{"title", "Shrek"}, {"plot", "After his swamp is filled with magical creatures, an ogre agrees to rescue a princess for a villainous lord in order to get his land back."}}

	result, err := coll.ReplaceOne(context.TODO(), filter, replacement)
	// end replace

	if err != nil {
		panic(err)
	}

	// When you run this file for the first time, it should print: "Matched and replaced an existing document."

	if result.MatchedCount != 0 {
		fmt.Println("Matched and replaced an existing document.")
		return
	}
}
