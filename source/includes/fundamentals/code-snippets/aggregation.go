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
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/")
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

	client.Database("tea").Collection("ratings").Drop(context.TODO())

	// begin insert docs
	coll := client.Database("tea").Collection("ratings")
	docs := []interface{}{
		bson.D{{"type", "Masala"}, {"rating", 10}},
		bson.D{{"type", "Earl Grey"}, {"rating", 5}},
		bson.D{{"type", "Masala"}, {"rating", 7}},
		bson.D{{"type", "Earl Grey"}, {"rating", 9}},
		bson.D{{"type", "Earl Grey"}, {"rating", 5}},
		bson.D{{"type", "Masala"}, {"rating", 9}},
		bson.D{{"type", "Earl Grey"}, {"rating", 8}},
		bson.D{{"type", "Masala"}, {"rating", 9}},
		bson.D{{"type", "Masala"}, {"rating", 10}},
		bson.D{{"type", "Earl Grey"}, {"rating", 10}},
		bson.D{{"type", "Masala"}, {"rating", 7}},
		bson.D{{"type", "Masala"}, {"rating", 5}},
		bson.D{{"type", "Masala"}, {"rating", 9}},
		bson.D{{"type", "Earl Grey"}, {"rating", 10}},
		bson.D{{"type", "Earl Grey"}, {"rating", 5}},
		bson.D{{"type", "Masala"}, {"rating", 9}},
		bson.D{{"type", "Earl Grey"}, {"rating", 7}},
		bson.D{{"type", "Masala"}, {"rating", 8}},
		bson.D{{"type", "Masala"}, {"rating", 9}},
	}

	result, err := coll.InsertMany(context.TODO(), docs)
	if err != nil {
		panic(err)
	}
	// end insert docs
	fmt.Printf("Number of documents inserted: %d\n", len(result.InsertedIDs))

	fmt.Println("Average:")
	{
		groupStage := bson.D{
			{"$group", bson.D{
				{"_id", "$type"},
				{"average", bson.D{
					{"$avg", "$rating"},
				}},
			}}}

		cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{groupStage})
		if err != nil {
			panic(err)
		}

		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		for _, result := range results {
			fmt.Printf("%v has an average rating of %v \n", result["_id"], result["average"])
		}
	}

	fmt.Println("Count:")
	{
		matchStage := bson.D{{"$match", bson.D{
			{"rating", bson.D{
				{"$gt", 8}},
			}},
		}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", "$type"},
			{"count", bson.D{
				{"$sum", 1},
			}},
		}}}

		cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage})
		if err != nil {
			panic(err)
		}

		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		for _, result := range results {
			fmt.Printf("%v Count: %v \n", result["_id"], result["count"])
		}
	}
}
