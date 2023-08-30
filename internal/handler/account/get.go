package handler_account

import (
	"github.com/alopez-2018459/go-the-field/internal/auth"
	"github.com/alopez-2018459/go-the-field/internal/responses"
	"github.com/alopez-2018459/go-the-field/internal/tags"
	tags_user "github.com/alopez-2018459/go-the-field/internal/tags/user"
	"github.com/gofiber/fiber/v2"
)

func SessionInfo(ctx *fiber.Ctx) error {
	sessionHeader := ctx.Get("Authorization")

	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return ctx.Status(401).
			JSON(fiber.Map{tags.ERROR: responses.INVALID_HEADER_ERROR, tags.MESSAGE: responses.UNAUTHORIZED_MESSAGE})
	}

	sessionId := sessionHeader[7:]

	user, err := auth.GetSession(sessionId)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: responses.GET_SESSION_ERROR + tags_user.TAG_SESSION, tags.MESSAGE: err.Error(), tags.STATUS: tags.UNAUTHENTICATED})
	}

	return ctx.Status(200).JSON(fiber.Map{tags.STATUS: tags.AUTHENTICATED, tags.USER: user})
}
