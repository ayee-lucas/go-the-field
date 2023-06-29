package main

import (
	"os"

	"github.com/alopez-2018459/go-the-field/internal/db"
	"github.com/alopez-2018459/go-the-field/internal/routes"
	"github.com/alopez-2018459/go-the-field/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	err := run()

	if err != nil {
		panic(err)
	}
}

func run() error {

	err := utils.LoadEnv()
	if err != nil {
		return err
	}

	err = db.InitDB()

	if err != nil {
		return err
	}

	defer db.CloseDB()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(recover.New())

	routes.SetupRoutes(app)

	var port string

	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	app.Listen(":" + port)

	return nil

}
