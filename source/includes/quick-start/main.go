package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	uri = os.Getenv("MONGODB_URI")
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	coll := client.Database("sample_mflix").Collection("movies")
	title := "Back to the Future"

	findResult := coll.FindOne(ctx, bson.D{{"title", title}})

	var doc bson.D
	if err = findResult.Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the title %s\n", title)
			return
		} else {
			log.Panic(findResult.Err().Error())
		}
	}

	jsonData, err := json.MarshalIndent(doc, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%s\n", jsonData)
}
