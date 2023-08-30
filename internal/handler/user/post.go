package handler_user

import (
	"strings"

	"github.com/alopez-2018459/go-the-field/internal/models"
	"github.com/alopez-2018459/go-the-field/internal/repository"
	"github.com/alopez-2018459/go-the-field/internal/responses"
	"github.com/alopez-2018459/go-the-field/internal/tags"
	tags_user "github.com/alopez-2018459/go-the-field/internal/tags/user"
	"github.com/alopez-2018459/go-the-field/internal/utils/validations"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type requestTeam struct {
	Country  string   `json:"country"  bson:"country"`
	Email    string   `json:"email"    bson:"email"`
	City     string   `json:"city"     bson:"city"`
	Links    []string `json:"links"    bson:"links"`
	Sport    string   `json:"sport"    bson:"sport"`
	Sponsors []string `json:"sponsors" bson:"sponsor"`
}

func RequestTeam(ctx *fiber.Ctx) error {
	body := new(requestTeam)

	param := ctx.Params("id")

	_, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.INVALID_ID})
	}

	user, err := repository.GetUserById(param)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.DATA_RETRIEVAL})
	}

	err = ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.PARSE_BODY_ERROR})
	}

	err = validations.IsStringEmpty(body.Country)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.REQUIRED_FIELD + tags_user.TAG_COUNTRY})
	}

	err = validations.IsStringEmpty(body.City)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.REQUIRED_FIELD + tags_user.TAG_CITY})
	}

	err = validations.IsEmailValid(body.Email)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.INVALID_EMAIL_FORMAT + tags_user.TAG_USER})
	}

	if len(body.Sport) <= 0 {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: responses.INVALID_INPUT + tags_user.TAG_SPORTS, tags.MESSAGE: responses.P_UPDATE_ERROR})
	}

	body.Email = strings.ToLower(body.Email)

	parsedId := user.ProfileID.Hex()

	profile, err := repository.GetUserProfileById(parsedId)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.DATA_RETRIEVAL})
	}

	if !profile.Type.IsZero() {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: responses.UPDATE_DATA_ERROR + tags_user.TAG_USER, tags.MESSAGE: responses.U_ORG_MESSAGE_ERROR})
	}

	team := &models.Team{
		Official: false,
		Country:  body.Country,
		Email:    body.Email,
		City:     body.City,
		Links:    body.Links,
		Sport:    body.Sport,
		Sponsors: body.Sponsors,
	}

	teamId, err := repository.SaveTeam(team)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.SAVE_DATA_MESSAGE_ERROR + tags_user.TAG_ORG})
	}

	parsedTeamId, err := primitive.ObjectIDFromHex(teamId)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_ORG, tags.MESSAGE: responses.INVALID_ID})
	}

	profileUpdateData := bson.D{{Key: "type_id", Value: parsedTeamId}}

	_, err = repository.UpdateProfile(user.ProfileID, profileUpdateData)

	if err != nil {
		_, err := repository.DeleteOrgById(teamId)
		if err != nil {
			return ctx.Status(500).
				JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.DELETE_DATA_MESSAGE_ERROR + tags_user.TAG_ORG})
		}
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.P_UPDATE_ERROR})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"result": fiber.Map{
			"message": "Team added successfully",
			"id":      teamId,
			"team":    team,
		},
	})
}
