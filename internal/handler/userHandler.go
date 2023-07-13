package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/alopez-2018459/go-the-field/internal/auth"
	"github.com/alopez-2018459/go-the-field/internal/models"
	"github.com/alopez-2018459/go-the-field/internal/repository"
	"github.com/alopez-2018459/go-the-field/internal/utils/validations"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func SignUp(ctx *fiber.Ctx) error {
	body := new(createUser)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Failed to parse request body", "message": err.Error()})
	}

	err = validations.IsStringEmpty(body.Username)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Username is required"})
	}

	err = validations.IsStringEmpty(body.Email)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Email is required"})
	}

	err = validations.IsEmailValid(body.Email)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Invalid email format"})
	}

	err = validations.IsStringEmpty(body.Password)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Password is required"})
	}

	err = validations.IsPasswordValid(body.Password)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid password", "message": err.Error()})
	}

	body.Password, err = validations.HashPassword(body.Password)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error(), "message": "Server Error"})
	}

	body.Username = strings.ToLower(body.Username)
	body.Email = strings.ToLower(body.Email)

	usernameExist, _ := repository.GetByUsername(body.Username)
	if usernameExist != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Error Saving User", "message": "Username already exists"})
	}

	emailExist, _ := repository.GetByEmail(body.Email)
	if emailExist != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Error Saving User", "message": "Email already exists"})
	}

	user := &models.User{
		Username:  body.Username,
		Email:     body.Email,
		Password:  body.Password,
		Online:    false,
		Role:      "user",
		Bio:       "No bio yet",
		Likes:     []primitive.ObjectID{},
		Followers: []primitive.ObjectID{},
		Posts:     []primitive.ObjectID{},
		Picture:   "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := repository.SaveUser(user)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to save user", "message": err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{
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

func SignIn(ctx *fiber.Ctx) error {

	body := new(loginUser)

	var sessionId string

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Failed to parse request body"})
	}

	err = validations.IsStringEmpty(body.Username)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Username is required"})
	}

	err = validations.IsStringEmpty(body.Password)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Password is required"})
	}

	usernameExist, err := repository.GetByUsername(body.Username)
	if usernameExist == nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "This username does not exist"})
	}

	err = validations.VerifyPassword(body.Password, usernameExist.Password)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Wrong Password"})
	}

	userSession := &models.UserSession{
		Sub:      usernameExist.ID,
		Username: usernameExist.Username,
		Email:    usernameExist.Email,
		Role:     usernameExist.Role,
		Picture:  usernameExist.Picture,
	}

	sessionId, err = auth.GenerateSession(userSession)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to generate session", "message": err.Error()})
	}

	ctx.Response().Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionId))

	return ctx.Status(200).JSON(fiber.Map{"message": "Logged in successfully", "session_id": sessionId})
}

func SignOut(ctx *fiber.Ctx) error {

	sessionHeader := ctx.Get("Authorization")

	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return ctx.Status(401).JSON(fiber.Map{"error": "Invalid header"})
	}

	sessionId := sessionHeader[7:]

	idDeleted, err := auth.SignOut(sessionId)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to sign out", "message": err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "Signed out successfully", "id": idDeleted})

}

func SessionInfo(ctx *fiber.Ctx) error {

	sessionHeader := ctx.Get("Authorization")

	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return ctx.Status(401).JSON(fiber.Map{"error": "Invalid header"})
	}

	sessionId := sessionHeader[7:]

	user, err := auth.GetSession(sessionId)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to get session", "message": err.Error(), "status": "unauthenticated"})
	}

	return ctx.Status(200).JSON(fiber.Map{"status": "authenticated", "user": user})

}
