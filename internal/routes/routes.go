package routes

import (
	"github.com/alopez-2018459/go-the-field/internal/handler"
	"github.com/alopez-2018459/go-the-field/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	apiGroup := app.Group("/api")
	{
		userGroup := apiGroup.Group("/users", middleware.EnsureAuth)
		userGroup.Get("/", handler.GetUsers)
	}
	{
		accountGroup := apiGroup.Group("/account")
		accountGroup.Post("/register", handler.SignUp)
		accountGroup.Post("/login", handler.SignIn)
		accountGroup.Post("/logout", handler.SignOut)
		accountGroup.Get("/me", handler.SessionInfo)

	}
}
