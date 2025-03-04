package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	// Replace the placeholder with your Atlas connection string
	const uri = "mongodb+srv://user:123@atlascluster.spm1ztf.mongodb.net/?retryWrites=true&w=majority&appName=AtlasCluster"

	// Connect to your Atlas cluster
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("failed to connect to the server: %v", err)
	}
	defer func() { _ = client.Disconnect(ctx) }()

	// Set the namespace
	coll := client.Database("sample_mflix").Collection("embedded_movies")

	// start-create-vector-search
	// Define the Atlas Vector Search index definition
	type vectorDefinitionField struct {
		Type          string `bson:"type"`
		Path          string `bson:"path"`
		NumDimensions int    `bson:"numDimensions"`
		Similarity    string `bson:"similarity"`
		Quantization  string `bson:"quantization"`
	}

	type vectorDefinition struct {
		Fields []vectorDefinitionField `bson:"fields"`
	}

	indexName := "vector_index"
	opts := options.SearchIndexes().SetName(indexName).SetType("vectorSearch")

	vectorSearchIndexModel := mongo.SearchIndexModel{
		Definition: vectorDefinition{
			Fields: []vectorDefinitionField{{
				Type:          "vector",
				Path:          "plot_embedding",
				NumDimensions: 1536,
				Similarity:    "dotProduct",
				Quantization:  "scalar"}},
		},
		Options: opts,
	}

	// Create the index
	searchIndexName := coll.SearchIndexes().CreateOne(ctx, vectorSearchIndexModel)
	// end-create-vector-search

	if err != nil {
		log.Fatalf("failed to create the search index: %v", err)
	}
	log.Println("New search index named " + searchIndexName + " is building.")

	// Await the creation of the index.
	log.Println("Polling to check if the index is ready. This may take up to a minute.")
	searchIndexes := coll.SearchIndexes()
	var doc bson.Raw
	for doc == nil {
		cursor, err := searchIndexes.List(ctx, options.SearchIndexes().SetName(searchIndexName))
		if err != nil {
			fmt.Errorf("failed to list search indexes: %w", err)
		}

		if !cursor.Next(ctx) {
			break
		}

		name := cursor.Current.Lookup("name").StringValue()
		queryable := cursor.Current.Lookup("queryable").Boolean()
		if name == searchIndexName && queryable {
			doc = cursor.Current
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	log.Println(searchIndexName + " is ready for querying.")

	// Creates an Atlas Search index
	// start-create-atlas-search
	// Define the Atlas Search index
	indexName := "atlas_search_index"
	opts := options.SearchIndexes().SetName(indexName).SetType("atlasSearch")

	indexModel := mongo.SearchIndexModel{
		Definition: map[string]interface{}{
			"mappings": map[string]interface{}{
				"dynamic": false,
				"fields": map[string]interface{}{
					"<fieldName>": map[string]<fieldType>{
						"type": "<fieldType>",
					},
				},
			},
		},
		Options: opts,
	}

	// Creates the index
	searchIndexName, err := coll.SearchIndexes().CreateOne(ctx, indexModel)
	
	// end-create-atlas-search

	// start-list-index
	// Specifies the options for the index to retrieve
	indexName := "<indexName>"
	opts := options.SearchIndexes().SetName(indexName)

	coll.SearchIndexes().List(ctx, opts)
	// end-list-index

	// start-update-index
	indexName := "<indexName>"

	type vectorDefinitionField struct {
		Type          string `bson:"type"`
		Path          string `bson:"path"`
		NumDimensions int    `bson:"numDimensions"`
		Similarity    string `bson:"similarity"`
	}
	type vectorDefinition struct {
		Fields []vectorDefinitionField `bson:"fields"`
	}
	definition := vectorDefinition{
		Fields: []vectorDefinitionField{{
			Type:          "vector",
			Path:          "<fieldToIndex>",
			NumDimensions: <numberOfDimensions>,
			Similarity:    "dotProduct"}},
	}
	coll.SearchIndexes().UpdateOne(ctx, indexName, definition)
	// end-update-index

	// start-delete-index
	coll.SearchIndexes().DropOne(ctx, "<indexName>")
	// end-delete-index
}
