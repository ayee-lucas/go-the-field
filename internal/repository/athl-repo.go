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

func SaveAthlete(athlete *models.Athlete) (string, error) {
	coll := db.GetDBCollection(db.COLL_ATHL)

	result, err := coll.InsertOne(context.Background(), athlete)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func DeleteAthleteById(id string) (*mongo.DeleteResult, error) {
	coll := db.GetDBCollection(db.COLL_ATHL)

	filter := bson.M{"_id": id}

	res, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	deletedCount := res.DeletedCount

	if deletedCount == 0 {
		return nil, fmt.Errorf("athlete info not found")
	}

	return res, nil
}
