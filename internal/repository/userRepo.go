package repository

import (
	"context"
	"fmt"

	"github.com/alopez-2018459/go-bank-system/internal/db"
	"github.com/alopez-2018459/go-bank-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllUsers() ([]models.User, error) {

	coll := db.GetDBCollection("users")

	users := make([]models.User, 0)

	cursor, err := coll.Find(context.Background(), bson.M{})

	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		user := models.User{}
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)

	}
	return users, nil

}

func GetById(id string) (*models.User, error) {

	coll := db.GetDBCollection("users")

	user := &models.User{}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = coll.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("User not found")
		}
		return nil, err
	}

	return user, nil
}

func GetByUsername(username string) (*models.User, error) {
	coll := db.GetDBCollection("users")

	user := &models.User{}

	err := coll.FindOne(context.Background(), bson.M{"username": username}).Decode(user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("User not found")
		}
		return nil, err
	}
	return user, nil
}

func GetByEmail(email string) (*models.User, error) {
	coll := db.GetDBCollection("users")

	user := &models.User{}

	err := coll.FindOne(context.Background(), bson.M{"email": email}).Decode(user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("User not found")
		}
		return nil, err
	}

	return user, nil

}

func SaveUser(users *models.User) (string, error) {
	coll := db.GetDBCollection("users")

	result, err := coll.InsertOne(context.Background(), users)

	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}
