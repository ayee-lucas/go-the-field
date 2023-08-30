package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alopez-2018459/go-the-field/internal/db"
	"github.com/alopez-2018459/go-the-field/internal/models"
)

func GetAllUsers() ([]models.User, error) {
	coll := db.GetDBCollection(db.COLL_USER)

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

func GetUserById(id string) (*models.User, error) {
	coll := db.GetDBCollection(db.COLL_USER)

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

func GetByEmail(email string) (*models.User, error) {
	coll := db.GetDBCollection(db.COLL_USER)

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
	coll := db.GetDBCollection(db.COLL_USER)

	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create indexes: %v", err)
	}

	result, err := coll.InsertOne(context.Background(), users)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func UpdateUser(id primitive.ObjectID, update bson.D) (*mongo.UpdateResult, error) {
	coll := db.GetDBCollection(db.COLL_USER)

	filter := bson.D{{Key: "_id", Value: id}}

	data := bson.D{{Key: "$set", Value: update}}

	result, err := coll.UpdateOne(context.Background(), filter, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetByUsername(username string) (*models.User, error) {
	coll := db.GetDBCollection(db.COLL_USER)

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
