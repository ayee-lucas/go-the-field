package repository

import (
	"context"

	"github.com/alopez-2018459/go-the-field/internal/db"
	"github.com/alopez-2018459/go-the-field/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func SaveSession(session *models.UserSession) (string, error) {

	coll := db.GetDBCollection("session")

	index := mongo.IndexModel{
		Keys: bsonx.Doc{{Key: "created_at", Value: bsonx.Int32(1)}},
		//		Options: options.Index().SetExpireAfterSeconds(int32(time.Now().Add(time.Hour * 24).Unix())),
		Options: options.Index().SetExpireAfterSeconds(10),
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
