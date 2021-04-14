package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var db *mongo.Database
var dbContext = context.TODO()

func Connect() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(dbContext, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(dbContext, nil)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database("mongo")
	usersCollection = db.Collection("users")
}
