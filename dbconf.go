package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collections = map[string]string{
	"user": "users",
}

var db *mongo.Database
var dbClient *mongo.Client

// Initializes a mongo.Database or returns the existing one.
func getDB() *mongo.Database {

	godotenv.Load()

	if db != nil {
		return db
	}

	dbURI := os.Getenv(MONGODB_URI)
	dbName := os.Getenv(DBNAME)

	if dbURI == "" {
		log.Fatal(ERR_EMPTY_DB_URI)
	}

	if dbName == "" {
		log.Fatal(ERR_EMPTY_DB_NAME)
	}

	dbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(ERR_DB_CONN)
	}

	db := dbClient.Database(dbName)
	return db
}

func closeDBClient() {
	if dbClient == nil {
		return
	}

	err := dbClient.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(ERR_DB_DISCONN)
	}

}

type dbController struct {
	collection *mongo.Collection
}
