package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ygbadamosi662/gk-go-mongo-migrations/util"
)

func Generate(description string) {
	timestamp := time.Now().Format("20060102150405_000")
	migrationUpName := fmt.Sprintf("%s_%s", timestamp, description)
	migrationFileName := fmt.Sprintf("%s_%s.go", timestamp, description)
	registrationKey := fmt.Sprintf("%s_%s.go", timestamp, description)

	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting project directory: %v", err)
	}

	migrationsDir := util.JoinPaths(projectRoot, "database", "migrations")

	migrationFilePath := util.JoinPaths(migrationsDir, migrationFileName)
	file, err := os.Create(migrationFilePath)
	if err != nil {
		log.Fatalf("Failed to create migration file: %v", err)
	}
	defer file.Close()

	migrationTemplate := fmt.Sprintf(`package migrations

import (
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	Registry["%s"] = Up_%s
}

func Up_%s(db *mongo.Database) error {
	log.Println("Running %s migration")

	// Write your migration logic here
	return nil
}
`, registrationKey, migrationUpName, migrationUpName, registrationKey)

	if _, err := file.Write([]byte(migrationTemplate)); err != nil {
		log.Fatalf("Failed to write to migration file: %v", err)
	}

	fmt.Printf("Generated new migration file: %s\n", migrationFileName)
}
