package main

import (
	"fiberTodo/src"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading:", err)
	}

	app := src.SetApp()

	port := os.Getenv("PORT")
	log.Println("Server starting on Port:", port)
	app.Listen(port)
}
