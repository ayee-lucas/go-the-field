package handler_user

import (
	"github.com/alopez-2018459/go-the-field/internal/repository"
	"github.com/alopez-2018459/go-the-field/internal/responses"
	"github.com/alopez-2018459/go-the-field/internal/tags"
	tags_user "github.com/alopez-2018459/go-the-field/internal/tags/user"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsers(ctx *fiber.Ctx) error {
	users, err := repository.GetAllUsers()
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: responses.DATA_RETRIEVAL + tags_user.TAG_USER, tags.MESSAGE: err.Error()})
	}
	if len(users) == 0 {
		return ctx.Status(404).
			JSON(fiber.Map{tags.ERROR: responses.DATA_NOT_FOUND + tags_user.TAG_USER, tags.MESSAGE: responses.DATA_NOT_FOUND_MESSAGE + tags_user.TAG_USER})
	}
	return ctx.Status(200).JSON(fiber.Map{tags.MESSAGE: responses.OK, tags.DATA: users})
}

func GetUserId(ctx *fiber.Ctx) error {
	param := ctx.Params("id")

	_, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.INVALID_ID})
	}

	user, err := repository.GetUserById(param)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.DATA_RETRIEVAL})
	}

	return ctx.Status(200).JSON(fiber.Map{tags.MESSAGE: responses.OK, tags.DATA: user})
}
