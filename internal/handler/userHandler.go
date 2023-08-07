package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alopez-2018459/go-the-field/internal/auth"
	"github.com/alopez-2018459/go-the-field/internal/models"
	"github.com/alopez-2018459/go-the-field/internal/repository"
	"github.com/alopez-2018459/go-the-field/internal/utils/validations"
)

func GetUsers(ctx *fiber.Ctx) error {
	users, err := repository.GetAllUsers()
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": "Failed to retrieve users", "message": err.Error()})
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
		return ctx.Status(500).
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to get user"})
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "user found", "user": user})
}

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
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Invalid Id"})
	}

	user, err = repository.GetUserById(param)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to get user"})
	}

	profile, err := repository.GetUserProfileById(user.ProfileID.Hex())

	if profile.Finished {
		return ctx.Status(400).
			JSON(fiber.Map{"error": "User already finished profile", "message": "User already finished profile"})
	}

	body := new(finishProfile)

	err = ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": "Failed to parse request body", "message": err.Error()})
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

	result, err = repository.UpdateProfile(user.ProfileID, data)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": err.Error(), "message": "Error updating profile"})
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "success", "Profile": result})
}

type updatePicture struct {
	Picture *models.Picture `json:"picture" bson:"picture"`
}

func UpdatePicture(ctx *fiber.Ctx) error {
	var result *mongo.UpdateResult

	param := ctx.Params("id")

	body := new(updatePicture)

	sessionHeader := ctx.Get("Authorization")

	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return ctx.Status(401).JSON(fiber.Map{"error": "Invalid header"})
	}

	sessionId := sessionHeader[7:]

	_, err := auth.GetSession(sessionId)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": "Failed to get session", "message": err.Error(), "status": "unauthenticated"})
	}

	id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Invalid Id"})
	}

	_, err = repository.GetUserById(param)

	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to get user"})
	}
	err = ctx.BodyParser(body)

	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": "Failed to parse request body", "message": err.Error()})
	}

	err = validations.IsStringEmpty(body.Picture.PictureKey)

	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "PictureKey is required"})
	}

	err = validations.IsStringEmpty(body.Picture.PictureURL)

	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "PictureURL is required"})
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
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to update user"})
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
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to update session"})
	}

	return ctx.Status(200).JSON(fiber.Map{"message": "success", "user": result})
}

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
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Invalid Id"})
	}

	user, err := repository.GetUserById(param)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to get user"})
	}

	err = ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "messsage": "Failed to parse request body"})
	}

	err = validations.IsStringEmpty(body.Country)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "Country is required"})
	}

	err = validations.IsStringEmpty(body.City)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "City is required"})
	}

	err = validations.IsEmailValid(body.Email)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "Invalid email format"})
	}

	if len(body.Sport) <= 0 {
		return ctx.Status(400).
			JSON(fiber.Map{"error": "Error with array lenght", "message": "Invalid sports quantity"})
	}

	body.Email = strings.ToLower(body.Email)

	parsedId := user.ProfileID.Hex()

	profile, err := repository.GetUserProfileById(parsedId)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to get profile"})
	}

	if !profile.Type.IsZero() {
		return ctx.Status(500).
			JSON(fiber.Map{"error": "Error adding org to user", "message": "This user already has an org or athlete attached"})
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
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to save org"})
	}

	parsedTeamId, err := primitive.ObjectIDFromHex(teamId)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to delete org"})
	}

	profileUpdateData := bson.D{{Key: "type_id", Value: parsedTeamId}}

	_, err = repository.UpdateProfile(user.ProfileID, profileUpdateData)

	if err != nil {

		resDeleteOrg, err := repository.DeleteOrgById(teamId)
		if err != nil {
			return ctx.Status(500).
				JSON(fiber.Map{"error": err.Error(), "message": "Failed to delete org"})
		}
		return ctx.Status(500).
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to add", "deleted": resDeleteOrg})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"result": fiber.Map{
			"message": "Team added successfully",
			"id":      teamId,
			"team":    team,
		},
	})
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
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Invalid Id"})
	}

	user, err := repository.GetUserById(param)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to get user"})
	}

	err = ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "messsage": "Failed to parse request body"})
	}

	err = validations.IsStringEmpty(body.Nationality)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "Nationality is required"})
	}

	err = validations.IsStringEmpty(body.Gender)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "Gender is required"})
	}

	err = validations.IsStringEmpty(body.Sport)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "Sport is required"})
	}

	if body.Height < 100 || body.Height > 200 {
		return ctx.Status(400).
			JSON(fiber.Map{"error": "Height values received not valid", "message": "Your height is not valid"})
	}

	if body.Weight < 100 || body.Weight > 400 {
		return ctx.Status(400).
			JSON(fiber.Map{"error": "Weight values received not valid", "message": "Your weight is not valid"})
	}

	body.Contact = strings.ToLower(body.Contact)

	err = validations.IsEmailValid(body.Contact)
	if err != nil {
		return ctx.Status(400).
			JSON(fiber.Map{"error": err.Error(), "message": "Your email is not valid"})
	}

	parsedId := user.ProfileID.Hex()
	profile, err := repository.GetUserProfileById(parsedId)
	if err != nil {
		return ctx.Status(500).
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to get profile"})
	}

	if !profile.Type.IsZero() {
		return ctx.Status(500).
			JSON(fiber.Map{"error": "Error adding info to user", "message": "This user already has an org or athlete attached"})
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
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to save athlete"})
	}

	profileUpdateData := bson.D{{Key: "athlete", Value: athleteId}}

	_, err = repository.UpdateProfile(profile.ID, profileUpdateData)

	if err != nil {

		resultAthlete, err := repository.DeleteAthleteById(athleteId)
		if err != nil {
			return ctx.Status(500).
				JSON(fiber.Map{"error": err.Error(), "message": "Failed to delete athlete"})
		}
		return ctx.Status(500).
			JSON(fiber.Map{"error": err.Error(), "message": "Failed to add", "deleted": resultAthlete})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"result": fiber.Map{
			"messsage": "Athlete added successfully",
			"id":       athleteId,
			"athlete":  athlete,
		},
	})
}
