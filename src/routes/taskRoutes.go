package routes

import (
	"fiberTodo/src/controllers"
	"fiberTodo/src/middleware"

	"github.com/gofiber/fiber/v3"
)

func TodoRoutes(app *fiber.App) {
	task := app.Group("/todo", middleware.AuthMiddleware)

	task.Post("/", controllers.CreateTodo)
	task.Get("/", controllers.GetTasks)
}
