package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "golangAPI", "users")
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
	_ = result // since we're not using result
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

	return c.Status(fiber.StatusCreated).JSON(responses.UserResponse{Status: fiber.StatusCreated, Message: "Success", Data: &fiber.Map{"data": newUser}})
}

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.UserResponse{Status: fiber.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
}

func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	validationErr := validate.Struct(&user)
	if validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.UserResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(fiber.StatusOK).JSON(responses.UserResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": updatedUser}})
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(fiber.StatusNotFound).JSON(
			responses.UserResponse{Status: fiber.StatusNotFound, Message: "Error", Data: &fiber.Map{"data": "User with specified ID not found!"}},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.UserResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		err := results.Decode(&singleUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.UserResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.UserResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": users}},
	)
}
