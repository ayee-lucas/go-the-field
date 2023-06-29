package handler

import (
	"strings"
	"time"

	"github.com/alopez-2018459/go-bank-system/internal/models"
	"github.com/alopez-2018459/go-bank-system/internal/repository"
	"github.com/alopez-2018459/go-bank-system/internal/utils/validations"
	"github.com/gofiber/fiber/v2"
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

type createUser struct {
	Username  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email,"`
	Password  string    `json:"password" bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func SaveUser(ctx *fiber.Ctx) error {
	body := new(createUser)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Failed to parse request body", "message": err.Error()})
	}

	err = validations.IsStringEmpty(body.Username)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Username is required"})
	}

	err = validations.IsStringEmpty(body.Email)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Email is required"})
	}

	err = validations.IsEmailValid(body.Email)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid email format", "message": err.Error()})
	}

	err = validations.IsStringEmpty(body.Password)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Password is required"})
	}

	err = validations.IsPasswordValid(body.Password)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid password", "message": err.Error()})
	}

	body.Password, err = validations.HashPassword(body.Password)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to hash password", "message": err.Error()})
	}

	body.Username = strings.ToLower(body.Username)
	body.Email = strings.ToLower(body.Email)

	usernameExist, _ := repository.GetByUsername(body.Username)
	if usernameExist != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Username already exists"})
	}

	emailExist, _ := repository.GetByEmail(body.Email)
	if emailExist != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Email already exists"})
	}

	user := &models.User{
		Username:  body.Username,
		Email:     body.Email,
		Password:  body.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := repository.SaveUser(user)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to save user", "message": err.Error()})
	}

	return ctx.Status(201).JSON(fiber.Map{
		"result": fiber.Map{
			"message": "User created successfully",
			"id":      id,
			"user":    user,
		},
	})
}

type loginUser struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

func Login(ctx *fiber.Ctx) error {

	body := new(loginUser)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Failed to parse request body", "message": err.Error()})
	}

	err = validations.IsStringEmpty(body.Username)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Username is required"})
	}

	err = validations.IsStringEmpty(body.Password)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Password is required"})
	}

	usernameExist, err := repository.GetByUsername(body.Username)
	if usernameExist == nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Username does not exist", "message": err.Error()})
	}

	err = validations.VerifyPassword(body.Password, usernameExist.Password)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid password", "message": err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "Logged in successfully"})
}
