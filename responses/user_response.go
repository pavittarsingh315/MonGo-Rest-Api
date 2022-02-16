package responses

import "github.com/gofiber/fiber/v2"

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"string"`
	Data    *fiber.Map `json:"data"`
}