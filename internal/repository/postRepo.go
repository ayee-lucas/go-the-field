package repository

import (
	"context"

	"github.com/alopez-2018459/go-the-field/internal/db"
	"github.com/alopez-2018459/go-the-field/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SavePost(post *models.Post) (string, error) {
	coll := db.GetDBCollection("posts")

	result, err := coll.InsertOne(context.Background(), post)

	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil

}
