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

type Config struct {
	MongoURL                    string `json:"mongo_url"`
	DBName                      string `json:"db_name"`
	AppliedMigrationsCollection string `json:"applied_migrations_collection"`
}

type Migration struct {
	FileName string
	Up       func(db *mongo.Database) error
}

func RunMigrations(registry map[string]MigrationFunc) error {
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := mongo.Connect(options.Client().
		ApplyURI(config.MongoURL))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(config.DBName)

	err = createAppliedMigrationsCollectionIfNotExists(db, config.AppliedMigrationsCollection)
	if err != nil {
		return fmt.Errorf("failed to create applied migrations collection: %w", err)
	}

	appliedMigrations, err := getAppliedMigrations(db, config.AppliedMigrationsCollection)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for migrationKey, migrationFunc := range registry {
		if !isMigrationApplied(appliedMigrations, migrationKey) {
			err := migrationFunc(db)
			if err != nil {
				return fmt.Errorf("error applying migration %s: %w", migrationKey, err)
			}

			err = markMigrationAsApplied(db, config.AppliedMigrationsCollection, migrationKey)
			if err != nil {
				return fmt.Errorf("failed to mark migration %s as applied: %w", migrationKey, err)
			}
		}
	}

	return nil
}

func loadConfig() (*Config, error) {
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
	collections, err := db.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		return fmt.Errorf("error listing collections: %w", err)
	}

	for _, collection := range collections {
		if collection == collectionName {
			return nil
		}
	}

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
	collection := db.Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), bson.D{
		{Key: "migration", Value: migrationFile}, {Key: "applied_at", Value: time.Now()},
	})
	if err != nil {
		return fmt.Errorf("error marking migration as applied: %w", err)
	}
	return nil
}
