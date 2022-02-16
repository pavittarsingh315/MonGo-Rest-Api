package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	// validate the request body
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "Error",
				Data: &fiber.Map{
					"data": err.Error(),
				},
			},
		)
	}

	// use the validator library to validate required fields
	validationErr := validate.Struct(&user)
	if validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Error",
			Data: &fiber.Map{
				"data": validationErr.Error(),
			},
		})
	}

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			responses.UserResponse{
				Status:  fiber.StatusBadRequest,
				Message: "Error",
				Data: &fiber.Map{
					"data": err.Error(),
				},
			},
		)
	}

	return c.Status(fiber.StatusCreated).JSON(responses.UserResponse{Status: fiber.StatusCreated, Message: "Success", Data: &fiber.Map{"data": result}})
}
