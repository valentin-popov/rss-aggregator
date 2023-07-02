package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collections = map[string]string{
	"user":       "users",
	"feed":       "feeds",
	"feedfollow": "feed_follows",
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
		log.Fatal(ERR_MSG_DB_EMPTY_URI)
	}

	if dbName == "" {
		log.Fatal(ERR_MSG_DB_EMPTY_NAME)
	}

	dbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(ERR_MSG_DB_CONN)
	}

	db = dbClient.Database(dbName)
	return db
}

func CloseDBClient() {
	if dbClient == nil {
		return
	}

	err := dbClient.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(ERR_MSG_DB_DISCONN)
	}

}
