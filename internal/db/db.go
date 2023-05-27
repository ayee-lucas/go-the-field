package db

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func GetDBCollection(col string) *mongo.Collection {

	return db.Collection(col)
}

func InitDB() error {
	uri := os.Getenv("MONGO_URI")

	if uri == "" {
		return errors.New("MONGO_URI not set")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		return err
	}

	db = client.Database("go_bankdb")

	return nil
}

func CloseDB() error {
	return db.Client().Disconnect(context.Background())
}
