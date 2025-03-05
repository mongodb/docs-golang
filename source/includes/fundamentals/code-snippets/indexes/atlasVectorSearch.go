package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	// Defines the structs used for the index definition
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

	// Sets the index name and type to "vectorSearch"
	indexName := "vector_index"
	opts := options.SearchIndexes().SetName(indexName).SetType("vectorSearch")

	// Defines the index definition
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

	// Creates the index
	searchIndexName, err := coll.SearchIndexes().CreateOne(ctx, vectorSearchIndexModel)
	// end-create-vector-search

	// Creates an Atlas Search index
	// start-create-atlas-search
	// Sets the index name and type to "search"
	indexName := "atlas_search_index"
	opts := options.SearchIndexes().SetName(indexName).SetType("search")

	// Defines the index definition 
	indexModel := mongo.SearchIndexModel{
		Definition: bson.D{
			{Key: "mappings", Value: bson.D{
				{Key: "dynamic", Value: false},
				{Key: "fields", Value: bson.D{
					{Key: "plot", Value: bson.D{
						{Key: "type", Value: "string"},
					}},
				}},
			}},
		},
		Options: opts,
	}

	// Creates the index
	searchIndexName, err := coll.SearchIndexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Fatalf("failed to create the search index: %v", err)
	}
	// end-create-atlas-search

	// start-list-index
	// Specifies the options for the index to retrieve
	indexName := "<indexName>"
	opts := options.SearchIndexes().SetName(indexName)

	cursor, err := coll.SearchIndexes().List(ctx, opts)

	// Print the index details to the console as JSON
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		log.Fatalf("failed to unmarshal results to bson: %v", err)
	}
	res, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("failed to marshal results to json: %v", err)
	}
	fmt.Println(res)
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
	err := coll.SearchIndexes().UpdateOne(ctx, indexName, definition)

	if err != nil {
		log.Fatalf("failed to update the index: %v", err)
	}
	// end-update-index

	// start-delete-index
	err := coll.SearchIndexes().DropOne(ctx, "<indexName>")
	if err != nil {
		log.Fatalf("failed to delete the index: %v", err)
	}
	// end-delete-index
}
