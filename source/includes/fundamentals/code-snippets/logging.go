package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/bombsimon/logrusr/v4"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
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
	//standardLogging(uri)
	//customLogging(uri)
	thirdPartyLogging(uri)
}

func standardLogging(uri string) {
	// start-standard-logger
	loggerOptions := options.
		Logger().
		SetMaxDocumentLength(25).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	clientOptions := options.
		Client().
		ApplyURI(uri).
		SetLoggerOptions(loggerOptions)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	// end-standard-logger
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
	}

	defer client.Disconnect(context.TODO())

	// start-insert
	coll := client.Database("testDB").Collection("testColl")
	_, err = coll.InsertOne(context.TODO(), bson.D{{"item", "grapefruit"}})
	// end-insert

	if err != nil {
		log.Fatalf("InsertOne failed: %v", err)
	}
}

// start-customlogger-struct
type CustomLogger struct {
	io.Writer
	mu sync.Mutex
}
// end-customlogger-struct

// start-customlogger-funcs
func (logger *CustomLogger) Info(level int, msg string, _ ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	if level == 1 {
		fmt.Fprintf(logger, "level: %d DEBUG, message: %s\n", level, msg)
	} else {
		fmt.Fprintf(logger, "level: %d INFO, message: %s\n", level, msg)
	}
}

func (logger *CustomLogger) Error(err error, msg string, _ ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	fmt.Fprintf(logger, "error: %v, message: %s\n", err, msg)
}
// end-customlogger-funcs

func customLogging(uri string) {
	// start-set-customlogger
	buf := bytes.NewBuffer(nil)
	sink := &CustomLogger{Writer: buf}

	loggerOptions := options.
		Logger().
		SetSink(sink).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug).
		SetComponentLevel(options.LogComponentConnection, options.LogLevelDebug)

	clientOptions := options.
		Client().
		ApplyURI(uri).
		SetLoggerOptions(loggerOptions)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	// end-set-customlogger

	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
	}

	defer client.Disconnect(context.TODO())

	coll := client.Database("testDB").Collection("testColl")
	_, err = coll.InsertOne(context.TODO(), bson.D{{"item", "grapefruit"}})
	if err != nil {
		log.Fatalf("InsertOne failed: %v", err)
	}
	fmt.Println(buf.String())
}

func thirdPartyLogging(uri string) {
    // start-make-logrus
	myLogger := &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg% <%commandName%>\n",
		},
	}
    // end-make-logrus

    // start-set-thirdparty-logger
	sink := logrusr.New(myLogger).GetSink()

	loggerOptions := options.
		Logger().
		SetSink(sink).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	clientOptions := options.
		Client().
		ApplyURI(uri).
		SetLoggerOptions(loggerOptions)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	// end-set-thirdparty-logger
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
	}

	defer client.Disconnect(context.TODO())

	// start-log-operations
	coll := client.Database("testDB").Collection("testColl")
	docs := []interface{}{
		bson.D{{"item", "starfruit"}},
		bson.D{{"item", "kiwi"}},
		bson.D{{"item", "cantaloupe"}},
	}
	_, err = coll.InsertMany(context.TODO(), docs)
	_, err = coll.DeleteOne(context.TODO(), bson.D{{"item", "kiwi"}})
	_, err = coll.UpdateOne(
		context.TODO(),
		bson.D{{"item", "cantaloupe"}},
		bson.D{{"$set", bson.D{{"qty", 3}}}},
	)
	// end-log-operations
	if err != nil {
		log.Fatalf("Operation failed: %v", err)
	}
}
