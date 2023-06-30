package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/alopez-2018459/go-the-field/internal/db"
	"github.com/alopez-2018459/go-the-field/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindSession(id string) (*models.UserSession, error) {
	coll := db.GetDBCollection("session")

	session := &models.UserSession{}

	err := coll.FindOne(context.Background(), bson.M{"_id": id}).Decode(session)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}

	return session, nil

}

func SaveSession(session *models.UserSession) (string, error) {

	coll := db.GetDBCollection("session")

	index := mongo.IndexModel{
		Keys: bson.M{"created_at": 1},
		//Options: options.Index().SetExpireAfterSeconds(int32(time.Now().Add(time.Hour * 24).Unix())),
		Options: options.Index().SetExpireAfterSeconds(86400),
	}

	_, err := coll.Indexes().CreateOne(context.Background(), index)

	if err != nil {
		return "", err
	}
	_, err = coll.InsertOne(context.Background(), session)

	if err != nil {
		return "", err
	}

	return session.ID, nil
}

func DeleteSession(id string) (string, error) {

	coll := db.GetDBCollection("session")

	filter := bson.M{"_id": id}

	res, err := coll.DeleteOne(context.Background(), filter)

	if err != nil {
		return "", err
	}

	deletedCount := res.DeletedCount

	if deletedCount == 0 {
		return "", errors.New("session not found")

	}

	return id, nil

}
