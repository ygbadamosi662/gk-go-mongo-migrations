package util

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Config holds the database connection and migration collection information
type Config struct {
	MongoURL                    string `json:"mongo_url"`
	DBName                      string `json:"db_name"`
	AppliedMigrationsCollection string `json:"applied_migrations_collection"`
}

// Migration represents the structure of each migration file
type Migration struct {
	FileName string
	Up       func(db *mongo.Database) error
}

func RunMigrations(registry map[string]MigrationFunc) error {
	// Load database connection information from config.json
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Connect to MongoDB
	client, err := mongo.Connect(options.Client().
		ApplyURI(config.MongoURL))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(config.DBName)

	// Create applied_migrations collection if it doesn't exist
	err = createAppliedMigrationsCollectionIfNotExists(db, config.AppliedMigrationsCollection)
	if err != nil {
		return fmt.Errorf("failed to create applied migrations collection: %w", err)
	}

	// Get the applied migrations from the database
	appliedMigrations, err := getAppliedMigrations(db, config.AppliedMigrationsCollection)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Iterate over migration files and apply any pending ones
	for migrationKey, migrationFunc := range registry {
		// If the migration file has not been applied yet
		if !isMigrationApplied(appliedMigrations, migrationKey) {
			// Run the migration function for each registered migration
			err := migrationFunc(db)
			if err != nil {
				return fmt.Errorf("error applying migration %s: %w", migrationKey, err)
			}

			// Mark the migration as applied
			err = markMigrationAsApplied(db, config.AppliedMigrationsCollection, migrationKey)
			if err != nil {
				return fmt.Errorf("failed to mark migration %s as applied: %w", migrationKey, err)
			}
		}
	}

	return nil // Success
}

func loadConfig() (*Config, error) {
	// Read the configuration from database/config.json
	configFile, err := os.Open("database/config.json")
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}
	return &config, nil
}

func createAppliedMigrationsCollectionIfNotExists(db *mongo.Database, collectionName string) error {
	// Check if the applied_migrations collection exists
	collections, err := db.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		return fmt.Errorf("error listing collections: %w", err)
	}

	// If the applied_migrations collection does not exist, create it
	for _, collection := range collections {
		if collection == collectionName {
			return nil // Collection already exists
		}
	}

	// Create the applied_migrations collection (optionally with a unique index on the filename)
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"migration", "applied_at"},
			"properties": bson.M{
				"migration": bson.M{
					"bsonType":    "string",
					"description": "must be a string and is required",
				},
				"applied_at": bson.M{
					"bsonType":    "date",
					"description": "must be a date and is required",
				},
			},
		},
	}
	cmd := bson.D{
		{Key: "create", Value: collectionName},
		{Key: "validator", Value: validator},
	}

	err = db.RunCommand(context.Background(), cmd).Err()
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("failed to create the applied_migrations collection with schema: %w", err)
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "migration", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = db.Collection(collectionName).Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return fmt.Errorf("failed to create unique index on migration: %w", err)
	}

	return nil
}

func getAppliedMigrations(db *mongo.Database, collectionName string) (map[string]struct{}, error) {
	// Get the applied migrations from the applied_migrations collection
	collection := db.Collection(collectionName)
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("error finding applied migrations: %w", err)
	}
	defer cursor.Close(context.Background())

	appliedMigrations := make(map[string]struct{})
	for cursor.Next(context.Background()) {
		var result struct {
			Migration string `bson:"migration"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("error decoding applied migration: %w", err)
		}
		appliedMigrations[result.Migration] = struct{}{}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return appliedMigrations, nil
}

func isMigrationApplied(appliedMigrations map[string]struct{}, migrationFile string) bool {
	_, applied := appliedMigrations[migrationFile]
	return applied
}

func markMigrationAsApplied(db *mongo.Database, collectionName string, migrationFile string) error {
	// Mark the migration as applied by inserting its name into the applied_migrations collection
	collection := db.Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), bson.D{
		{Key: "migration", Value: migrationFile}, {Key: "applied_at", Value: time.Now()},
	})
	if err != nil {
		return fmt.Errorf("error marking migration as applied: %w", err)
	}
	return nil
}
