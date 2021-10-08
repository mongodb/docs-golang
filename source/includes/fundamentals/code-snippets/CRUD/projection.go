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
	if uri = os.Getenv("DRIVER_REF_URI"); uri == "" {
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

	// begin insertDocs
	coll := client.Database("tea").Collection("ratings")
	docs := []interface{}{
		bson.D{{"type", "Masala"}, {"rating", 10}},
		bson.D{{"type", "Assam"}, {"rating", 5}},
		bson.D{{"type", "Oolong"}, {"rating", 7}},
		bson.D{{"type", "Earl Grey"}, {"rating", 8}},
		bson.D{{"type", "English Breakfast"}, {"rating", 5}},
	}

	result, insertErr := coll.InsertMany(context.TODO(), docs)
	if insertErr != nil {
		panic(insertErr)
	}
	//end insertDocs
	fmt.Printf("Number of documents inserted: %d\n", len(result.InsertedIDs))

	fmt.Println("Exclude Projection:")
	//begin exclude projection
	excludeProjection := bson.D{{"rating", 0}}
	excludeOptions := options.Find().SetProjection(excludeProjection)

	excludeCursor, excludeErr := coll.Find(context.TODO(), bson.D{}, excludeOptions)

	var excludeResults []bson.D
	if excludeErr = excludeCursor.All(context.TODO(), &excludeResults); excludeErr != nil {
		panic(excludeErr)
	}
	for _, result := range excludeResults {
		fmt.Println(result)
	}
	//end exclude projection

	fmt.Println("Include Projection:")
	//begin include projection
	includeProjection := bson.D{{"type", 1}, {"rating", 1}, {"_id", 0}}
	includeOptions := options.Find().SetProjection(includeProjection)

	includeCursor, includeErr := coll.Find(context.TODO(), bson.D{}, includeOptions)

	var includeResults []bson.D
	if includeErr = includeCursor.All(context.TODO(), &includeResults); includeErr != nil {
		panic(includeErr)
	}
	for _, result := range includeResults {
		fmt.Println(result)
	}
	//end include projection

	fmt.Println("Aggregation Projection:")
	// begin aggregate projection
	projectStage := bson.D{{"$project", bson.D{{"type", 1}, {"rating", 1}, {"_id", 0}}}}

	aggCursor, aggErr := coll.Aggregate(context.TODO(), mongo.Pipeline{projectStage})
	if aggErr != nil {
		panic(aggErr)
	}

	var aggResults []bson.D
	if aggErr = aggCursor.All(context.TODO(), &aggResults); aggErr != nil {
		panic(aggErr)
	}
	for _, result := range aggResults {
		fmt.Println(result)
	}
	// end aggregate projection
}
