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

	fmt.Println("FindOneAndDelete:")
	//begin FindOneAndDelete
	deleteFilter := bson.D{{"type", "Assam"}}

	var deleteResult bson.D
	deleteErr := coll.FindOneAndDelete(context.TODO(), deleteFilter).Decode(&deleteResult)
	if deleteErr != nil {
		panic(deleteErr)
	}
	
	fmt.Println(deleteResult)
	//end FindOneAndDelete

	fmt.Println("FindOneAndReplace:")
	//begin FindOneAndReplace
	replaceFilter := bson.D{{"type", "English Breakfast"}}
	replaceDocument := bson.D{{"type", "Ceylon"}, {"rating", 6}}
	replaceOptions := options.FindOneAndReplace().SetReturnDocument(options.After)

	var replaceResult bson.D
	replaceErr := coll.FindOneAndReplace(context.TODO(), replaceFilter, replaceDocument, replaceOptions).Decode(&replaceResult)
	if replaceErr != nil {
		panic(replaceErr)
	}
	
	fmt.Println(replaceResult)
	//end FindOneAndReplace

	fmt.Println("FindOneAndUpdate:")
	//begin FindOneAndUpdate
	updateFilter := bson.D{{"type", "Oolong"}}
	updateDocument := bson.D{{"$set", bson.D{{"rating", 9}}}}
	updateOptions := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updateResult bson.D
	updateErr := coll.FindOneAndUpdate(context.TODO(), updateFilter, updateDocument, updateOptions).Decode(&updateResult)
	if updateErr != nil {
		panic(updateErr)
	}
	
	fmt.Println(updateResult)
	//end FindOneAndUpdate
}
