package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	// start-logger
	loggerOptions := options.
		Logger().
		SetMaxDocumentLength(25).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	clientOptions := options.
		Client().
		ApplyURI(uri).
		SetLoggerOptions(loggerOptions)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	// end-logger
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
	}

	defer client.Disconnect(context.TODO())

	// start-log-insert
	coll := client.Database("testDB").Collection("testColl")
	_, err = coll.InsertOne(context.TODO(), bson.D{{"item", "grapefruit"}, {"qty", 4}})
	// end-log-insert

	if err != nil {
		log.Fatalf("InsertOne failed: %v", err)
	}
}