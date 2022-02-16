package main

import (
	"fiber-mongo-api/configs"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	configs.ConnectDb()

	app.Listen(":8000")
}