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
			JSON(fiber.Map{tags.ERROR: responses.SAVE_DATA_MESSAGE_ERROR + tags_user.TAG_SESSION, tags.MESSAGE: responses.USERNAME_EXISTS})
	}

	emailExist, _ := repository.GetByEmail(body.Email)
	if emailExist != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: responses.SAVE_DATA_MESSAGE_ERROR + tags_user.TAG_SESSION, tags.MESSAGE: responses.EMAIL_EXISTS})
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
			JSON(fiber.Map{tags.ERROR: responses.SAVE_DATA_MESSAGE_ERROR + tags_user.TAG_SESSION, tags.MESSAGE: err.Error()})
	}

	parsedProfileId, err := primitive.ObjectIDFromHex(idProfile)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.INVALID_ID})
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
			JSON(fiber.Map{tags.ERROR: responses.SAVE_DATA_MESSAGE_ERROR + tags_user.TAG_SESSION, tags.MESSAGE: err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"result": fiber.Map{
			tags.MESSAGE: "User created successfully",
			"id":         id,
			tags.USER:    account,
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
			JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.BODY_PARSE_MESSAGE})
	}

	err = validations.IsStringEmpty(body.Username)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.REQUIRED_FIELD + tags_user.TAG_USERNAME})
	}

	err = validations.IsStringEmpty(body.Password)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.REQUIRED_FIELD + tags_user.TAG_PASS})
	}

	usernameExist, err := repository.GetByUsername(body.Username)
	if usernameExist == nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.USERNAME_NOTFOUND})
	}

	err = validations.VerifyPassword(body.Password, usernameExist.Password)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.WRONG_PASSWORD})
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
			JSON(fiber.Map{tags.ERROR: "Failed to generate session", tags.MESSAGE: err.Error()})
	}

	ctx.Response().Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionId))

	return ctx.Status(200).
		JSON(fiber.Map{tags.MESSAGE: "Logged in successfully", "session_id": sessionId})
}

func SignOut(ctx *fiber.Ctx) error {
	sessionHeader := ctx.Get("Authorization")

	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return ctx.Status(401).JSON(fiber.Map{tags.ERROR: "Invalid header"})
	}

	sessionId := sessionHeader[7:]

	idDeleted, err := auth.SignOut(sessionId)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: "Failed to sign out", tags.MESSAGE: err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{tags.MESSAGE: "Signed out successfully", "id": idDeleted})
}
