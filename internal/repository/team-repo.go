package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alopez-2018459/go-the-field/internal/db"
	"github.com/alopez-2018459/go-the-field/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveTeam(org *models.Team) (string, error) {
	coll := db.GetDBCollection(db.COLL_TEAM)

	result, err := coll.InsertOne(context.Background(), org)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func GetOrgById(id string) (*models.Team, error) {
	coll := db.GetDBCollection(db.COLL_TEAM)

	org := &models.Team{}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = coll.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(org)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("Team not found")
		}
		return nil, err
	}

	return org, nil
}

func GetOrgByEmail(email string) (*models.Team, error) {
	coll := db.GetDBCollection(db.COLL_TEAM)

	org := &models.Team{}

	err := coll.FindOne(context.Background(), bson.M{"email": email}).Decode(org)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("Team not found")
		}
		return nil, err
	}
	return org, nil
}

func DeleteOrgById(id string) (*mongo.DeleteResult, error) {
	coll := db.GetDBCollection(db.COLL_TEAM)

	filter := bson.M{"_id": id}

	res, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	deletedCount := res.DeletedCount

	if deletedCount == 0 {
		return nil, errors.New("Team not found")
	}

	return res, nil
}
