package handler

import (
	"time"

	"github.com/alopez-2018459/go-the-field/internal/models"
	"github.com/alopez-2018459/go-the-field/internal/repository"
	"github.com/alopez-2018459/go-the-field/internal/utils/validations"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createPost struct {
	Author    primitive.ObjectID `json:"author,omitempty" bson:"author,omitempty"`
	Content   models.PostContent `json:"content" bson:"content"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func UploadPost(ctx *fiber.Ctx) error {
	body := new(createPost)

	err := ctx.BodyParser(body)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "Failed to parse body"})
	}

	_, err = repository.GetUserById(body.Author.Hex())

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error(), "message": "User not found"})
	}

	err = validations.ValidatePostContent(&body.Content)

	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Failed to validate post content", "message": err.Error()})
	}

	post := &models.Post{
		Author:    body.Author,
		Content:   body.Content,
		Repost:    []primitive.ObjectID{},
		Starred:   []primitive.ObjectID{},
		Likes:     []primitive.ObjectID{},
		Comments:  []primitive.ObjectID{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := repository.SavePost(post)

	if err != nil {
		ctx.Status(500).JSON(fiber.Map{"error": err.Error(), "message": "Failed to save post"})
	}

	post.ID, _ = primitive.ObjectIDFromHex(id)

	return ctx.Status(200).JSON(fiber.Map{"result": fiber.Map{
		"message": "Post created successfully",
		"id":      id,
		"post":    post,
	}})

}
