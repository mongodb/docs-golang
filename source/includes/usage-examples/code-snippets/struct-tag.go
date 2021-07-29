package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace the uri string with your MongoDB deployment's connection string.
const uri = "mongodb+srv://<username>:<password>@<cluster-address>/test?w=majority"

type BlogPost struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title,omitempty"`
	Author    string             `bson:"author,omitempty"`
	WordCount int                `bson:"word_count,omitempty"`
	Tags      []string           `bson:"tags,omitempty"`
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// begin create and insert
	myCollection := client.Database("sample_training").Collection("blogPosts")

	post1 := BlogPost{
		Title:     "Caring for your Monstera plant",
		WordCount: 478,
		Tags:      []string{"plant care", "gardening", "housekeeping"},
	}

	post2 := BlogPost{
		Title:     "Annuals vs. Perennials?",
		Author:    "Sam Lee",
		WordCount: 682,
		Tags:      []string{"flowering plants", "gardening"},
	}

	docs := []interface{}{post1, post2}

	insertResult, err := myCollection.InsertMany(context.TODO(), docs)
	// end create and insert

	_ = insertResult

	if err != nil {
		panic(err)
	}

	cursor, err := myCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
