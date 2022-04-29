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
		bson.D{{"type", "Masala"}, {"rating", 10}, {"visits", 24}},
		bson.D{{"type", "Earl Grey"}, {"rating", 5}, {"visits", 7}},
		bson.D{{"type", "Masala"}, {"rating", 7}, {"visits", 10}},
		bson.D{{"type", "Earl Grey"}, {"rating", 9}, {"visits", 12}},
		bson.D{{"type", "Earl Grey"}, {"rating", 5}, {"visits", 4}},
		bson.D{{"type", "Masala"}, {"rating", 9}, {"visits", 18}},
		bson.D{{"type", "Earl Grey"}, {"rating", 8}, {"visits", 15}},
		bson.D{{"type", "Masala"}, {"rating", 9}, {"visits", 14}},
		bson.D{{"type", "Masala"}, {"rating", 10}, {"visits", 24}},
		bson.D{{"type", "Earl Grey"}, {"rating", 10}, {"visits", 19}},
		bson.D{{"type", "Masala"}, {"rating", 7}, {"visits", 13}},
		bson.D{{"type", "Masala"}, {"rating", 5}, {"visits", 8}},
		bson.D{{"type", "Masala"}, {"rating", 9}, {"visits", 21}},
		bson.D{{"type", "Earl Grey"}, {"rating", 10}, {"visits", 17}},
		bson.D{{"type", "Earl Grey"}, {"rating", 5}, {"visits", 5}},
		bson.D{{"type", "Masala"}, {"rating", 9}, {"visits", 19}},
		bson.D{{"type", "Earl Grey"}, {"rating", 7}, {"visits", 14}},
		bson.D{{"type", "Masala"}, {"rating", 8}, {"visits", 17}},
		bson.D{{"type", "Masala"}, {"rating", 9}, {"visits", 20}},
	}

	result, err := coll.InsertMany(context.TODO(), docs)
	if err != nil {
		panic(err)
	}
	// end insert docs
	fmt.Printf("Number of documents inserted: %d\n", len(result.InsertedIDs))

	fmt.Println("Average:")
	{
		// create the stage
		groupStage := bson.D{
			{"$group", bson.D{
				{"_id", "$type"},
				{"average", bson.D{
					{"$avg", "$rating"},
				}},
				{"count", bson.D{
					{"$sum", 1},
				}},
			}}}

		// pass the stage into a pipeline
		// pass the pipeline as the second paramter in the Aggregate() method
		cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{groupStage})
		if err != nil {
			panic(err)
		}

		// display the results
		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		for _, result := range results {
			fmt.Printf("%v has an average rating of %v \n", result["_id"], result["average"])
			fmt.Printf("%v Count: %v \n", result["_id"], result["count"])
		}
	}

	fmt.Println("Unset:")
	{
		// create the stages
		matchStage := bson.D{{"$match", bson.D{
			{"rating", bson.D{
				{"$gt", 8}},
			}},
		}}
		unsetStage := bson.D{{"$unset",bson.A{"_id", "rating"},}}
		sortStage := bson.D{{"$sort", bson.D{
			{"visits", -1},
			{"type", 1}},
		}}
		limitStage := bson.D{{"$limit", 5}}

		// pass the stage into a pipeline
		// pass the pipeline as the second paramter in the Aggregate() method
		cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, unsetStage, sortStage, limitStage})
		if err != nil {
			panic(err)
		}

		// display the results
		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		for _, result := range results {
			fmt.Println(result)
		}
	}
}
