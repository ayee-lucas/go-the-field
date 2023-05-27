package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/alopez-2018459/go-bank-system/internal/db"
	"github.com/alopez-2018459/go-bank-system/internal/models"
	"github.com/alopez-2018459/go-bank-system/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUsers(ctx *fiber.Ctx) error {

	coll := db.GetDBCollection("users")

	users := make([]models.User, 0)

	cursor, err := coll.Find(ctx.Context(), bson.M{})

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	for cursor.Next(ctx.Context()) {
		user := models.User{}
		err := cursor.Decode(&user)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return ctx.Status(404).JSON(fiber.Map{"error": "No users found"})

	}

	return ctx.Status(200).JSON(users)

}

type createUser struct {
	USERNAMWE string    `json:"username" bson:"username"`
	EMAIL     string    `json:"email" bson:"email,"`
	PASSWORD  string    `json:"password" bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func SaveUser(ctx *fiber.Ctx) error {
	body := new(createUser)

	err := ctx.BodyParser(&body)

	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid Body"})
	}

	if body.EMAIL == "" || body.PASSWORD == "" || body.USERNAMWE == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Missing fields"})
	}
	body.PASSWORD, err = utils.HashPassword(body.PASSWORD)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Error hashing password"})
	}

	body.USERNAMWE = strings.ToLower(body.USERNAMWE)
	body.EMAIL = strings.ToLower(body.EMAIL)
	body.CreatedAt = time.Now()
	body.UpdatedAt = time.Now()

	fmt.Println(body)

	coll := db.GetDBCollection("users")
	result, err := coll.InsertOne(ctx.Context(), body)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Error saving user", "message": err.Error()})
	}

	return ctx.Status(201).JSON(fiber.Map{
		"result": fiber.Map{
			"id":   result,
			"user": body,
		},
	})

}
