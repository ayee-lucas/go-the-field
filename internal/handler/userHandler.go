package handler

import (
	"github.com/alopez-2018459/go-the-field/internal/models"
	"github.com/alopez-2018459/go-the-field/internal/repository"
	"github.com/alopez-2018459/go-the-field/internal/utils/validations"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsers(ctx *fiber.Ctx) error {
	users, err := repository.GetAllUsers()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to retrieve users", "message": err.Error()})
	}
	if len(users) == 0 {
		return ctx.Status(404).JSON(fiber.Map{"error": "No users found"})
	}
	return ctx.Status(200).JSON(users)
}

func GetUserId(ctx *fiber.Ctx) error {
	param := ctx.Params("id")

	_, err := primitive.ObjectIDFromHex(param)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Invalid Id"})
	}

	user, err := repository.GetUserById(param)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error(), "message": "Failed to get user"})
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "user found", "user": user})

}

type finishProfile struct {
	Name string `json:"name" bson:"name"`
	Bio  string `json:"bio" bson:"bio"`
}

func FinishProfile(ctx *fiber.Ctx) error {
	var result *mongo.UpdateResult
	var user *models.User

	param := ctx.Params("id")
	id, err := primitive.ObjectIDFromHex(param)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Invalid Id"})
	}

	user, err = repository.GetUserById(param)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error(), "message": "Failed to get user"})
	}

	if user.Finished {
		return ctx.Status(400).JSON(fiber.Map{"error": "User already finished profile", "message": "User already finished profile"})
	}

	body := new(finishProfile)

	err = ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Failed to parse request body", "message": err.Error()})
	}

	err = validations.IsStringEmpty(body.Name)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Name is required"})
	}

	err = validations.IsStringEmpty(body.Bio)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Bio is required"})
	}

	data := bson.D{{
		Key:   "name",
		Value: body.Name,
	}, {
		Key:   "bio",
		Value: body.Bio,
	}, {
		Key:   "finished",
		Value: true,
	}}

	result, err = repository.UpdateUser(id, data)

	return ctx.Status(200).JSON(fiber.Map{"message": "user found", "user": result})
}
