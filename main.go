package main

import (
	"os"

	"github.com/bmdavis419/fiber-mongo-example/common"
	"github.com/bmdavis419/fiber-mongo-example/router"
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

	err := common.LoadEnv()

	if err != nil {
		return err
	}

	//	defer common.CloseDB()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(recover.New())

	router.AddBookGroup(app)

	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	app.Listen(":" + port)

	return nil

}
