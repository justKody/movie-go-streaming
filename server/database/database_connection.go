package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client = DBInstance()

func DBInstance() *mongo.Client {
	err := godotenv.Load()

	if err != nil {
		log.Println("Warning: Unable to find .env file")
	}

	MongoDb := os.Getenv("DATABASE_URI")

	if MongoDb == "" {
		log.Fatal("DATABASE_URI not set!")
	}

	fmt.Println("MongoDb URI: ", MongoDb)

	clientOptions := options.Client().ApplyURI(MongoDb)

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		return nil
	}

	return client
}

func OpenCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Unable to find .env file")
	}

	databaseName := os.Getenv("DATABASE_NAME")

	fmt.Println("DATABASE_NAME: ", databaseName)

	collection := Client.Database(databaseName).Collection(collectionName)

	if collection == nil {
		return nil
	}

	return collection
}
