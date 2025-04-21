# gk-go-mongo-migrations

A lightweight Go-based migration and generator tool for MongoDB-backed projects. It allows you to easily run versioned migrations and manage your project structure.

---

## ✨ Features

- Built-in MongoDB migration runner with support for versioned migration files.
- Easy scaffolding via CLI commands like `gk init`, `gk generate <description...>`.
- Installable as a module or usable directly in any Go project.

---

## 🚀 Installation

You can install `gk` either globally or use it directly within your Go projects.

### Option 1: Install Globally (Recommended)

To install `gk` as a global command-line tool, run the following command:

```bash
go install github.com/ygbadamosi662/gk-go-mongo-migrations/cmd/gk@latest
```
This will install the gk tool and make it available globally on your system, allowing you to run the gk commands from anywhere.

Note: Ensure that $GOPATH/bin or $GOBIN is in your system's $PATH to use gk globally.

### Option 2: Install Using go get
Alternatively, you can use go get to install the tool. This will download the latest version and install it globally if $GOPATH/bin or $GOBIN is in your $PATH.

```bash
go get github.com/ygbadamosi662/gk-go-mongo-migrations/cmd/gk@latest
```
This command will fetch and install the latest version of gk.

Note: Ensure that $GOPATH/bin (or $GOBIN) is in your $PATH so that you can run gk globally.

### Option 3: Use in Go Projects
If you want to use gk as part of your Go project without installing it globally, you can import it directly into your code:

```go
import "github.com/ygbadamosi662/gk-go-mongo-migrations"
```
You can then use the gk commands within your Go code as needed.


## 🧱 Project Scaffolding
### To initialize a new project structure, run:

```bash
gk init
```
This will scaffold the following structure:
```
your-project/
├── gk_migrate.go
├── database/
│   ├── config.json                 # MongoDB connection + migration config
│   └── migrations/
│       └── registry.go            # Auto-registered migrations live here
```
### 🧬 Generate a Migration
To generate a new migration file, use the following command:

```bash
gk generate <description>
```
This creates a new file in the database/migrations/ directory like:

database/migrations/20250418142926.535_add_address.go
The contents of the generated file will be:

```go
package migrations

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	Registry["20250418142926_535_add_address.go"] = Up_0250418142926_535_add_address
}

func Up_0250418142926_535_add_address(db *mongo.Database) error {
	log.Println("Running 20250418142926_535_add_address.go migration")
	// Add your logic here
	return nil
}
```
You only need to implement the logic inside the Up() function to perform your migration.

## ⚡ Running Migrations
Your database/config.json file should look like this:

```json
{
  "mongo_url": "mongodb://user:password@127.0.0.1:27017?authSource=db_name",
  "db_name": "test_db",
  "applied_migrations_collection": "migrations"
}
```
Then run:

```bash
go run ./gk_migrate.go
```
This will:

Connect to MongoDB using the given mongo_url.

Check the applied_migrations_collection to track which migrations have been run.

Run all new migrations found in database/migrations/registry.go.

✅ Already applied migrations will be skipped.

🔁 Example Migration Logic
Here’s a sample migration logic you can test with:

```go
func Up(db *mongo.Database) error {
	log.Println("Creating test document in 'test_collection'")
	_, err := db.Collection("test_collection").InsertOne(context.Background(), map[string]interface{}{
		"created_at": time.Now(),
		"status":     "testing",
	})
	return err
}
```
## 🧩 How It Works
All generated migration files auto-register themselves inside the global Registry map.

registry.go imports these files, ensuring they are registered at runtime.

gk_migrate.go reads your config.json and runs all unapplied migrations in filename order.

## 📝 License
This project is licensed under the MIT License - see the LICENSE file for details.

## 🌐 GitHub Repository
Visit the GitHub repository for more information
