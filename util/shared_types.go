package util

import "go.mongodb.org/mongo-driver/v2/mongo"

type MigrationFunc func(*mongo.Database) error
