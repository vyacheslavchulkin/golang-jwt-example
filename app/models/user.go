package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type User struct {
	ID           string    `bson:"_id"`
	CreatedAt    time.Time `bson:"created_at"`
	Name         string    `bson:"name"`
	Admin        bool      `bson:"admin"`
	Password     string    `bson:"password"`
	RefreshToken string    `bson:"refresh_uuid"`
}

var usersCollection *mongo.Collection

func CreateUser(user User) error {
	_, err := usersCollection.InsertOne(dbContext, user)
	return err
}

func FindOneUser(filter interface{}) (*User, error) {
	var user *User
	err := usersCollection.FindOne(dbContext, filter).Decode(&user)
	return user, err
}

func FindOneUserAndUpdate(filter interface{}, update interface{}) (*User, error) {
	var user *User
	err := usersCollection.FindOneAndUpdate(dbContext, filter, bson.D{{"$set", update}}).Decode(&user)
	return user, err
}
