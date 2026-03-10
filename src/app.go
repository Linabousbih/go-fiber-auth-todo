package src

import (
	"fiberTodo/src/database"

	"github.com/gofiber/fiber/v3"
)

func SetApp() *fiber.App {
	app := fiber.New()

	database.ConnectDB()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Todo list with fiber")
	})

	return app
}
