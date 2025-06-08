package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/solrac97gr/telegram-followers-checker/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Index struct {
	Field      string `json:"field"`
	Collection string `json:"collection"`
	Type       string `json:"type"` // e.g., "text", "hashed", etc.
}

type Databases struct {
	Name        string   `json:"name"`
	Collections []string `json:"collections"`
	Indexes     []Index  `json:"indexes"`
}

func NewDatabaseConfig(config *config.Config) ([]Databases, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	return []Databases{
		{
			Name: config.InfluencersDBName,
			Collections: []string{
				InfluencersCollectionName,
			},
			Indexes: []Index{
				{Field: "link", Collection: InfluencersCollectionName, Type: "text"},
				{Field: "user_id", Collection: InfluencersCollectionName, Type: "hashed"},
			},
		},
		{
			Name: config.UsersDBName,
			Collections: []string{
				UserCollectionName,
				UserProfileCollectionName,
				UserTokenCollectionName,
			},
			Indexes: []Index{
				{Field: "email", Collection: UserCollectionName, Type: "text"},
				{Field: "user_id", Collection: UserTokenCollectionName, Type: "hashed"},
				{Field: "user_id", Collection: UserProfileCollectionName, Type: "hashed"},
			},
		},
	}, nil
}

func InitializeDatabase(db *mongo.Client, config *config.Config) error {
	databases, err := NewDatabaseConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create database configuration: %w", err)
	}

	if err := isDatabaseRunning(db); err != nil {
		return fmt.Errorf("database is not up to date: %w", err)
	}

	// Iterate through each database configuration
	for _, database := range databases {
		log.Printf("Initializing database: %s", database.Name)

		// Log that collections will be created automatically on first write
		for _, collection := range database.Collections {
			log.Printf("Collection %s in database %s will be created automatically on first write", collection, database.Name)
		}

		// Create indexes for this database
		for _, index := range database.Indexes {
			if err := createIndexIfNotExists(db, database.Name, index.Collection, index.Field, index.Type); err != nil {
				// Log the error but don't return it if index already exists
				if !isAlreadyExistsError(err) {
					return fmt.Errorf("failed to create index %s on collection %s in database %s: %w", index.Field, index.Collection, database.Name, err)
				}
				log.Printf("Index %s on collection %s in database %s already exists, skipping creation", index.Field, index.Collection, database.Name)
			} else {
				log.Printf("Successfully created index %s on collection %s in database %s", index.Field, index.Collection, database.Name)
			}
		}
	}

	log.Printf("Database initialization completed successfully")
	return nil
}

func createIndexIfNotExists(db *mongo.Client, databaseName, collectionName, field, indexType string) error {
	ctx := context.Background()
	collection := db.Database(databaseName).Collection(collectionName)

	var indexModel mongo.IndexModel

	// Create index based on type
	switch indexType {
	case "text":
		indexModel = mongo.IndexModel{
			Keys:    bson.M{field: "text"},
			Options: options.Index().SetName(field + "_text_index"),
		}
	case "hashed":
		indexModel = mongo.IndexModel{
			Keys:    bson.M{field: "hashed"},
			Options: options.Index().SetName(field + "_hashed_index"),
		}
	default:
		// Default to ascending index
		indexModel = mongo.IndexModel{
			Keys:    bson.M{field: 1},
			Options: options.Index().SetName(field + "_index").SetUnique(false),
		}
	}

	if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		return err
	}

	return nil
}

// Helper function to check if an error indicates that something already exists
func isAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	errorMessage := err.Error()
	return strings.Contains(errorMessage, "already exists") ||
		strings.Contains(errorMessage, "E11000") || // MongoDB duplicate key error
		strings.Contains(errorMessage, "IndexOptionsConflict") ||
		strings.Contains(errorMessage, "NamespaceExists")
}

func isDatabaseRunning(db *mongo.Client) error {
	// ping to the database to check if it is up and running
	ctx := context.Background()
	if err := db.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}
