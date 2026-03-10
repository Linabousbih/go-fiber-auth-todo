package controllers

import (
	"fiberTodo/models"
	"fiberTodo/src/database"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTodo(c fiber.Ctx) error {

	userId := c.Locals("userId").(string)

	type body struct {
		Title       string `json:"title"`
		Description string `json:"desription"`
		Status      string `json:"status"`
	}

	var data body

	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// If we do not get a status from the user will default to incomplete
	if data.Status != string(models.StatusCompleted) && data.Status != string(models.StatusIncomplete) {
		data.Status = string(models.StatusIncomplete)
	}

	todo := bson.M{
		"_id":         primitive.NewObjectID(),
		"title":       data.Description,
		"description": data.Description,
		"status":      data.Status,
		"userId":      userId,
	}

	_, err := database.DB.Collection("todos").InsertOne(c.Context(), todo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create task",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"message": "Task created successfully",
			"todo":    todo,
		},
	)
}

func GetTasks(c fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	cursor, err := database.DB.Collection("todos").Find(c.Context(), bson.M{"userId": userId})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch task"})
	}

	var tasks []bson.M

	if err := cursor.All(c.Context(), &tasks); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse tasks"})
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"todos": tasks,
		},
	)
}
