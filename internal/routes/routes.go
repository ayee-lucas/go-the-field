package routes

import (
	"github.com/alopez-2018459/go-bank-system/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	apiGroup := app.Group("/api")
	{
		userGroup := apiGroup.Group("/users")
		userGroup.Get("/", handler.GetUsers)
		userGroup.Post("/", handler.SaveUser)
	}
}
