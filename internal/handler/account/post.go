package handler_account

import (
	"fmt"
	"strings"
	"time"

	"github.com/alopez-2018459/go-the-field/internal/auth"
	"github.com/alopez-2018459/go-the-field/internal/models"
	"github.com/alopez-2018459/go-the-field/internal/repository"
	"github.com/alopez-2018459/go-the-field/internal/responses"
	"github.com/alopez-2018459/go-the-field/internal/tags"
	tags_user "github.com/alopez-2018459/go-the-field/internal/tags/user"
	"github.com/alopez-2018459/go-the-field/internal/utils/validations"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type signUpBody struct {
	Username  string    `json:"username"   bson:"username"`
	Email     string    `json:"email"      bson:"email,"`
	Password  string    `json:"password"   bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func SignUp(ctx *fiber.Ctx) error {
	body := new(signUpBody)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags.ACCOUNT, tags.MESSAGE: responses.PARSE_BODY_ERROR})
	}

	err = validations.IsStringEmpty(body.Username)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: responses.REQUIRED_FIELD + tags_user.TAG_USERNAME, tags.MESSAGE: responses.SIGN_UP_MESSAGE_ERROR})
	}

	err = validations.IsStringEmpty(body.Email)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_EMAIL, tags.MESSAGE: responses.SIGN_UP_MESSAGE_ERROR})
	}

	err = validations.IsEmailValid(body.Email)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_EMAIL, tags.MESSAGE: responses.SIGN_UP_MESSAGE_ERROR})
	}

	err = validations.IsStringEmpty(body.Password)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_PASS, tags.MESSAGE: responses.SIGN_UP_MESSAGE_ERROR})
	}

	err = validations.IsPasswordValid(body.Password)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: responses.INVALID_INPUT + tags_user.TAG_PASS, tags.MESSAGE: err.Error()})
	}

	body.Password, err = validations.HashPassword(body.Password)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_PASS, tags.MESSAGE: responses.SIGN_UP_MESSAGE_ERROR})
	}

	body.Username = strings.ToLower(body.Username)
	body.Email = strings.ToLower(body.Email)

	usernameExist, _ := repository.GetByUsername(body.Username)
	if usernameExist != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": "Error Saving User", "message": "Username already exists"})
	}

	emailExist, _ := repository.GetByEmail(body.Email)
	if emailExist != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": "Error Saving User", "message": "Email already exists"})
	}

	userProfile := &models.Profile{
		Name:           "",
		Bio:            "",
		PreferedSports: []string{},
		Online:         false,
		Finished:       false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	idProfile, err := repository.SaveProfile(userProfile)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": "Failed to save user", "message": err.Error()})
	}

	parsedProfileId, err := primitive.ObjectIDFromHex(idProfile)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": err.Error(), "message": "Fatal Error Parsing id"})
	}

	account := &models.User{
		Username:      body.Username,
		Email:         body.Email,
		EmailVerified: false,
		Password:      body.Password,
		Role:          "user",
		Verified:      false,
		Picture:       models.Picture{},
		ProfileID:     parsedProfileId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	id, err := repository.SaveUser(account)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": "Failed to save user", "message": err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"result": fiber.Map{
			"message": "User created successfully",
			"id":      id,
			"user":    account,
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
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to parse request body"})
	}

	err = validations.IsStringEmpty(body.Username)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "Username is required"})
	}

	err = validations.IsStringEmpty(body.Password)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "Password is required"})
	}

	usernameExist, err := repository.GetByUsername(body.Username)
	if usernameExist == nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "This username does not exist"})
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
		return ctx.Status(500).
			JSON(fiber.Map{"error": "Failed to generate session", "message": err.Error()})
	}

	ctx.Response().Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionId))

	return ctx.Status(200).
		JSON(fiber.Map{"message": "Logged in successfully", "session_id": sessionId})
}

func SignOut(ctx *fiber.Ctx) error {
	sessionHeader := ctx.Get("Authorization")

	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return ctx.Status(401).JSON(fiber.Map{"error": "Invalid header"})
	}

	sessionId := sessionHeader[7:]

	idDeleted, err := auth.SignOut(sessionId)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": "Failed to sign out", "message": err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "Signed out successfully", "id": idDeleted})
}

/**
*
* 	ADMIN LOGIN
*
 */

func SignInAdmin(ctx *fiber.Ctx) error {
	body := new(loginUser)

	var sessionId string

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to parse request body"})
	}

	err = validations.IsStringEmpty(body.Username)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "Username is required"})
	}

	err = validations.IsStringEmpty(body.Password)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "Password is required"})
	}

	usernameExist, err := repository.GetByUsername(body.Username)
	if usernameExist == nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "This username does not exist"})
	}

	err = validations.VerifyPassword(body.Password, usernameExist.Password)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Wrong Password"})
	}

	if usernameExist.Role != "ADMIN" {
		return ctx.Status(401).
			JSON(fiber.Map{"error": "Unauthorized", "message": "You are not an ADMIN"})
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
		return ctx.Status(500).
			JSON(fiber.Map{"error": "Failed to generate session", "message": err.Error()})
	}

	ctx.Response().Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionId))

	return ctx.Status(200).
		JSON(fiber.Map{"message": "Logged in successfully", "session_id": sessionId})
}
