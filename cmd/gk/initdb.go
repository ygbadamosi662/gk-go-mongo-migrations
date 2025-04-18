package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ygbadamosi662/gk-go-mongo-migrations/util"
)

type Config struct {
	MongoURL                    string `json:"mongo_url"`
	DBName                      string `json:"db_name"`
	AppliedMigrationsCollection string `json:"applied_migrations_collection"`
}

func InitDB() {
	// Get the root project directory
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting project directory: %v", err)
	}

	// Define the migrations directory and config file path
	migrationsDir := util.JoinPaths(projectRoot, "database", "migrations")
	configFilePath := util.JoinPaths(projectRoot, "database", "config.json")

	// Create the migrations directory if it doesn't exist
	if err := util.CreateDirIfNotExist(migrationsDir); err != nil {
		log.Fatalf("Error creating migrations directory: %v", err)
	}

	// Check if the config file already exists
	if util.FileExists(configFilePath) {
		log.Println("Config file already exists, skipping creation.")
		return
	}

	// Create a new config structure with default values
	config := Config{
		MongoURL:                    "mongodb://username:password@localhost:27017", // Replace with default or user input
		DBName:                      "testdb",                                      // Replace with default or user input
		AppliedMigrationsCollection: "migrations",                                  // Default collection name
	}

	// Create the config file
	file, err := os.Create(configFilePath)
	if err != nil {
		log.Fatalf("Error creating config file: %v", err)
	}
	defer file.Close()

	// Marshal the config struct to JSON and write to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // For pretty printing
	if err := encoder.Encode(config); err != nil {
		log.Fatalf("Error writing to config file: %v", err)
	}

	fmt.Println("Migration config initialized successfully.")

	createRegistry()
	createMigrator(projectRoot)
}

func createRegistry() {
	registryFile := "database/migrations/registry.go"

	// Create registry.go if it doesn't exist
	if _, err := os.Stat(registryFile); os.IsNotExist(err) {
		registryContent := `package migrations

import "github.com/ygbadamosi662/gk-go-mongo-migrations/util"

var Registry = make(map[string]util.MigrationFunc)
`
		err := os.WriteFile(registryFile, []byte(registryContent), 0644)
		if err != nil {
			log.Fatalf("Error writing registry.go: %v", err)
		}
		log.Println("Created database/migrations/registry.go")
	} else {
		log.Println("database/migrations/registry.go already exists")
	}
}

func createMigrator(projectRoot string) {
	migratorFile := filepath.Join(projectRoot, "gk_migrate.go")

	if _, err := os.Stat(migratorFile); os.IsNotExist(err) {
		migratorContent := fmt.Sprintf(`
package main

import (
	"fmt"
	"log"

	"github.com/ygbadamosi662/gk-go-mongo-migrations/util"
	"%s/database/migrations"
)

func main() {
	// Call the migration runner from your module
	err := util.RunMigrations(migrations.Registry)
	if err != nil {
		log.Fatalf("Migration failed: %%v", err)
	} else {
		fmt.Println("Migrations applied successfully!")
	}
}`, getModuleName(projectRoot))
		err := os.WriteFile(migratorFile, []byte(migratorContent), 0644)
		if err != nil {
			log.Fatalf("Error writing gk_migrate.go: %v", err)
		}
		log.Println("Created gk_migrate.go")
	} else {
		log.Println("gk_migrate.go already exists")
	}
}

func getModuleName(projectRoot string) string {
	data, _ := os.ReadFile(filepath.Join(projectRoot, "go.mod"))
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}
	return ""
}
