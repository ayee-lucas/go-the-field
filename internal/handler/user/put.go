package handler_user

import (
	"strings"

	"github.com/alopez-2018459/go-the-field/internal/auth"
	"github.com/alopez-2018459/go-the-field/internal/models"
	"github.com/alopez-2018459/go-the-field/internal/repository"
	"github.com/alopez-2018459/go-the-field/internal/responses"
	"github.com/alopez-2018459/go-the-field/internal/tags"
	tags_user "github.com/alopez-2018459/go-the-field/internal/tags/user"
	"github.com/alopez-2018459/go-the-field/internal/utils/validations"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type finishProfile struct {
	Name string `json:"name" bson:"name"`
	Bio  string `json:"bio"  bson:"bio"`
}

func FinishProfile(ctx *fiber.Ctx) error {
	var result *mongo.UpdateResult
	var user *models.User

	param := ctx.Params("id")
	_, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.INVALID_ID})
	}

	user, err = repository.GetUserById(param)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.DATA_RETRIEVAL})
	}

	profile, err := repository.GetUserProfileById(user.ProfileID.Hex())
	if err != nil {
		return ctx.Status(404).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.DATA_NOT_FOUND})
	}

	if profile.Finished {
		return ctx.Status(409).
			JSON(fiber.Map{tags.ERROR: responses.PROFILE_FINISHED_ERROR + tags_user.TAG_USER, tags.MESSAGE: responses.P_FINISHED_MESSAGE_ERROR})
	}

	body := new(finishProfile)

	err = ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: responses.PARSE_BODY_ERROR + tags_user.TAG_USER, tags.MESSAGE: responses.BODY_PARSE_MESSAGE})
	}

	err = validations.IsStringEmpty(body.Name)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.REQUIRED_FIELD + tags_user.TAG_NAME})
	}

	err = validations.IsStringEmpty(body.Bio)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.REQUIRED_FIELD + tags_user.TAG_BIO})
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

	result, err = repository.UpdateProfile(user.ProfileID, data)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: responses.UPDATE_DATA_ERROR + tags_user.TAG_USER, tags.MESSAGE: responses.P_UPDATE_ERROR})
	}

	return ctx.Status(200).JSON(fiber.Map{tags.MESSAGE: responses.OK, tags.DATA: result})
}

// UPDATE PICTURE

type updatePicture struct {
	Picture *models.Picture `json:"picture" bson:"picture"`
}

func UpdatePicture(ctx *fiber.Ctx) error {
	var result *mongo.UpdateResult

	param := ctx.Params("id")

	body := new(updatePicture)

	sessionHeader := ctx.Get("Authorization")

	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return ctx.Status(401).
			JSON(fiber.Map{tags.ERROR: responses.INVALID_HEADER_ERROR + tags_user.TAG_USER, tags.MESSAGE: responses.UNAUTHORIZED_MESSAGE})
	}

	sessionId := sessionHeader[7:]

	_, err := auth.GetSession(sessionId)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: responses.GET_SESSION_ERROR + tags_user.TAG_USER, tags.MESSAGE: err.Error(), tags.STATUS: tags.UNAUTHENTICATED})
	}

	id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.INVALID_ID})
	}

	_, err = repository.GetUserById(param)

	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.DATA_NOT_FOUND_MESSAGE})
	}
	err = ctx.BodyParser(body)

	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: responses.PARSE_BODY_ERROR + tags_user.TAG_USER, tags.MESSAGE: responses.BODY_PARSE_MESSAGE})
	}

	err = validations.IsStringEmpty(body.Picture.PictureKey)

	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_PICTURE, tags.MESSAGE: responses.REQUIRED_FIELD + tags_user.TAG_PICTURE_KEY})
	}
	err = validations.IsStringEmpty(body.Picture.PictureURL)

	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_PICTURE, tags.MESSAGE: responses.REQUIRED_FIELD + tags_user.TAG_PICTURE_URL})
	}

	data := bson.D{
		{
			Key: "picture",
			Value: bson.D{
				{Key: "pictureKey", Value: body.Picture.PictureKey},
				{Key: "pictureURL", Value: body.Picture.PictureURL},
			},
		},
	}

	result, err = repository.UpdateUser(id, data)

	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.P_UPDATE_ERROR})
	}

	sessionData := bson.D{
		{
			Key: "picture",
			Value: bson.D{
				{Key: "pictureKey", Value: body.Picture.PictureKey},
				{Key: "pictureURL", Value: body.Picture.PictureURL},
			},
		},
	}

	_, err = repository.UpdateSession(sessionId, sessionData)

	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_SESSION, tags.MESSAGE: responses.P_UPDATE_ERROR + tags_user.TAG_SESSION})
	}

	return ctx.Status(200).JSON(fiber.Map{tags.MESSAGE: responses.OK, tags.DATA: result})
}

type requestAthlete struct {
	Nationality  string   `json:"nationality"  bson:"nationality"`
	Gender       string   `json:"gender"       bson:"gender"`
	Sport        string   `json:"sport"        bson:"sport"`
	Sponsors     []string `json:"sponsors"     bson:"sponsors"`
	CurrentTeam  string   `json:"current_team" bson:"current_team"`
	Height       int      `json:"height"       bson:"height"`
	Weight       int      `json:"weight"       bson:"weight"`
	Achievements string   `json:"achievements" bson:"achievements"`
	Contact      string   `json:"contact"      bson:"contact"`
}

func RequestAthlete(ctx *fiber.Ctx) error {
	body := new(requestAthlete)

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

	err = ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.BODY_PARSE_MESSAGE})
	}

	err = validations.IsStringEmpty(body.Nationality)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.REQUIRED_FIELD + tags_user.TAG_NATIONALITY})
	}

	err = validations.IsStringEmpty(body.Gender)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.REQUIRED_FIELD + tags_user.TAG_GENDER})
	}

	err = validations.IsStringEmpty(body.Sport)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error(), tags.MESSAGE: responses.REQUIRED_FIELD + tags_user.TAG_SPORTS})
	}

	if body.Height < 100 || body.Height > 200 {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: responses.INVALID_INPUT + tags_user.TAG_HEIGHT, tags.MESSAGE: responses.SAVE_DATA_MESSAGE_ERROR})
	}

	if body.Weight < 100 || body.Weight > 400 {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: responses.INVALID_INPUT + tags_user.TAG_WEIGHT, tags.MESSAGE: responses.SAVE_DATA_MESSAGE_ERROR})
	}

	body.Contact = strings.ToLower(body.Contact)

	err = validations.IsEmailValid(body.Contact)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.SAVE_DATA_MESSAGE_ERROR})
	}

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

	athlete := &models.Athlete{
		Nationality:  body.Nationality,
		Gender:       body.Gender,
		Sport:        body.Sport,
		Sponsors:     body.Sponsors,
		CurrentTeam:  body.CurrentTeam,
		Height:       body.Height,
		Weight:       body.Weight,
		Achievements: body.Achievements,
		Contact:      body.Contact,
	}

	athleteId, err := repository.SaveAthlete(athlete)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.SAVE_DATA_MESSAGE_ERROR})
	}

	profileUpdateData := bson.D{{Key: "athlete", Value: athleteId}}

	_, err = repository.UpdateProfile(profile.ID, profileUpdateData)

	if err != nil {

		_, err := repository.DeleteAthleteById(athleteId)
		if err != nil {
			return ctx.Status(500).
				JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.DELETE_DATA_MESSAGE_ERROR})
		}
		return ctx.Status(500).
			JSON(fiber.Map{tags.ERROR: err.Error() + tags_user.TAG_USER, tags.MESSAGE: responses.SAVE_DATA_MESSAGE_ERROR})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"result": fiber.Map{
			"messsage": "Athlete added successfully",
			"id":       athleteId,
			"athlete":  athlete,
		},
	})
}
