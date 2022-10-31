package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func main() {

	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	database := client.Database("myDB")
	coll := database.Collection("myColl")

	// start-session
	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(context.TODO())

	err = mongo.WithSession(context.TODO(), session, func(ctx mongo.SessionContext) error {
		if err = session.StartTransaction(txnOptions); err != nil {
			return err
		}

		docs := []interface{}{
			bson.D{{"title", "The Bluest Eye"}, {"author", "Toni Morrison"}},
			bson.D{{"title", "Sula"}, {"author", "Toni Morrison"}},
			bson.D{{"title", "Song of Solomon"}, {"author", "Toni Morrison"}},
		}
		result, err := coll.InsertMany(ctx, docs)
		if err != nil {
			return err
		}

		if err = session.CommitTransaction(ctx); err != nil {
			return err
		}

		fmt.Println(result.InsertedIDs)
		return nil
	})
	if err != nil {
		if err := session.AbortTransaction(context.TODO()); err != nil {
			panic(err)
		}
		panic(err)
	}
	// end-session
}
