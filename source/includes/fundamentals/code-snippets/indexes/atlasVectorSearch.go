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
	// Defines the Atlas Vector Search index definition
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

	// Creates the index
	searchIndexName, err := coll.SearchIndexes().CreateOne(ctx, vectorSearchIndexModel)
	// end-create-vector-search

	// Creates an Atlas Search index
	// start-create-atlas-search
	// Defines the Atlas Search index
	indexName := "atlas_search_index"
	opts := options.SearchIndexes().SetName(indexName).SetType("atlasSearch")

	indexModel := mongo.SearchIndexModel{
		Definition: map[string]interface{}{
			"mappings": map[string]interface{}{
				"dynamic": false,
				"fields": map[string]interface{}{
					"plot": map[string]string{
						"type": "string",
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
