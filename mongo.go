package main

import (
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var users mongo.Collection
var appointments mongo.Collection
var pictures mongo.Collection

func getMongoUrlFromArgs() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return ""
}

func getMongoUrl() string {
	urlsFromArgs := getMongoUrlFromArgs()

	if urlsFromArgs != "" {
		return urlsFromArgs
	}

	var DatabaseHost = "localhost"
	var DatabasePort = "27017"
	var DatabaseUser = "diligner"
	var DatabasePassword = "vermont2042Inmassacusset"
	var DatabaseName = "matrix"
	var env = "DEV"

	if env != "PROD" {
		return fmt.Sprintf("mongodb://%s:%s", DatabaseHost, DatabasePort)
	}

	return fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", DatabaseUser, DatabasePassword, DatabaseHost, DatabaseName)
}

func connectTodataBase() {
	DatabaseName := "matrix"
	MongoUrl := getMongoUrl()
	clientOptions := options.Client().ApplyURI(MongoUrl)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error while connecting to the database, error: %v", err)
	}

	log.Println("Successfully connect to the database")

	users = *client.Database(DatabaseName).Collection("users")
	appointments = *client.Database(DatabaseName).Collection("appointments")
	pictures = *client.Database(DatabaseName).Collection("pictures")
}
