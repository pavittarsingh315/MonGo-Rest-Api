package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	// all routes related to users
	app.Post("/user", controllers.CreateUser)
}
