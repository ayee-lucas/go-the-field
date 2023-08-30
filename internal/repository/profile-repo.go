package repository

import (
	"context"
	"fmt"

	"github.com/alopez-2018459/go-the-field/internal/db"
	"github.com/alopez-2018459/go-the-field/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserProfileById(id string) (*models.Profile, error) {
	coll := db.GetDBCollection(db.COLL_PROFILE)

	profile := &models.Profile{}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = coll.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(profile)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("Profile not found")
		}
		return nil, err
	}

	return profile, nil
}

func UpdateProfile(id primitive.ObjectID, update bson.D) (*mongo.UpdateResult, error) {
	coll := db.GetDBCollection(db.COLL_PROFILE)

	filter := bson.D{{Key: "_id", Value: id}}

	data := bson.D{{Key: "$set", Value: update}}

	result, err := coll.UpdateOne(context.Background(), filter, data)
	if err != nil {
		return nil, err
	}

	return result, err
}

func SaveProfile(profile *models.Profile) (string, error) {
	coll := db.GetDBCollection(db.COLL_PROFILE)

	result, err := coll.InsertOne(context.Background(), profile)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}
