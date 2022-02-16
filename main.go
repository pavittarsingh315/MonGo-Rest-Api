package main

import (
	"fiber-mongo-api/configs"
	"fiber-mongo-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(helmet.New())

	configs.ConnectDb()

	routes.UserRoute(app)

	app.Listen(":8000")
}
