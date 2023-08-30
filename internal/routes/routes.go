package routes

import (
	handler_account "github.com/alopez-2018459/go-the-field/internal/handler/account"
	handler_user "github.com/alopez-2018459/go-the-field/internal/handler/user"
	"github.com/alopez-2018459/go-the-field/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	apiGroup := app.Group("/api")
	{
		userGroup := apiGroup.Group("/users", middleware.EnsureAuth)
		userGroup.Get("/", handler_user.GetUsers)
		userGroup.Get("/:id", handler_user.GetUserId)
		userGroup.Post("/request/team/:id", handler_user.RequestTeam)
		userGroup.Post("/request/athl/:id", handler_user.RequestAthlete)
		userGroup.Put("/finish/:id", handler_user.FinishProfile)
		userGroup.Put("/picture/:id", handler_user.UpdatePicture)
	}
	{
		accountGroup := apiGroup.Group("/account")
		accountGroup.Post("/register", handler_account.SignUp)
		accountGroup.Post("/login", handler_account.SignIn)
		accountGroup.Post("/logout", handler_account.SignOut)
		accountGroup.Get("/me", handler_account.SessionInfo)

	}

	/**
	 * 	ADMIN PORTAL ROUTES
	 * 	Routes for the admin portal angular app
	 **/
	// adminGroup := app.Group("/admin")
	// {
	// 	accountAdmin := adminGroup.Group("/account")
	// 	accountAdmin.Post("login")
	//
	// }
}
