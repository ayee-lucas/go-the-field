package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alopez-2018459/go-the-field/internal/db"
	"github.com/alopez-2018459/go-the-field/internal/models"
)

func GetAllUsers() ([]models.User, error) {
	coll := db.GetDBCollection("User")

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
	coll := db.GetDBCollection("User")

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

func GetUserProfileById(id string) (*models.Profile, error) {
	coll := db.GetDBCollection("Profile")

	profile := &models.Profile{}

	objectID, err := primitive.ObjectIDFromHex(id)

	err = coll.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(profile)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("Profile not found")
		}
		return nil, err
	}

	return profile, nil
}

func GetByUsername(username string) (*models.User, error) {
	coll := db.GetDBCollection("User")

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
	coll := db.GetDBCollection("User")

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
	coll := db.GetDBCollection("User")

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

func SaveProfile(profile *models.Profile) (string, error) {
	coll := db.GetDBCollection("Profile")

	result, err := coll.InsertOne(context.Background(), profile)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func UpdateUser(id primitive.ObjectID, update bson.D) (*mongo.UpdateResult, error) {
	coll := db.GetDBCollection("User")

	filter := bson.D{{Key: "_id", Value: id}}

	data := bson.D{{Key: "$set", Value: update}}

	result, err := coll.UpdateOne(context.Background(), filter, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateProfile(id primitive.ObjectID, update bson.D) (*mongo.UpdateResult, error) {
	coll := db.GetDBCollection("Profile")

	filter := bson.D{{Key: "_id", Value: id}}

	data := bson.D{{Key: "$set", Value: update}}

	result, err := coll.UpdateOne(context.Background(), filter, data)
	if err != nil {
		return nil, err
	}

	return result, err
}

func SaveTeam(org *models.Team) (string, error) {
	coll := db.GetDBCollection("Team")

	result, err := coll.InsertOne(context.Background(), org)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func GetOrgById(id string) (*models.Team, error) {
	coll := db.GetDBCollection("Team")

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
	coll := db.GetDBCollection("Team")

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
	coll := db.GetDBCollection("Team")

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

func SaveAthlete(athlete *models.Athlete) (string, error) {
	coll := db.GetDBCollection("Athlete")

	result, err := coll.InsertOne(context.Background(), athlete)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func DeleteAthleteById(id string) (*mongo.DeleteResult, error) {
	coll := db.GetDBCollection("Athlete")

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
