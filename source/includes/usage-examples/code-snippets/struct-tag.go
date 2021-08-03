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

// begin struct
type BlogPost struct {
	Title     string
	Author    string
	WordCount int `bson:"word_count"`
	Tags      []string
}

// end struct

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

	// begin create and insert
	myCollection := client.Database("sample_training").Collection("posts")

	post := BlogPost{
		Title:     "Annuals vs. Perennials?",
		Author:    "Sam Lee",
		WordCount: 682,
		Tags:      []string{"seasons", "gardening", "flower"},
	}

	_, err = myCollection.InsertOne(context.TODO(), post)
	// end create and insert

	if err != nil {
		panic(err)
	}

	var result bson.D
	err = myCollection.FindOne(context.TODO(), bson.D{{"author", "Sam Lee"}}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		panic(err)
	}

	fmt.Printf("Found document: %v", result)
}
