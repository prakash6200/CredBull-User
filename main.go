package main

import (
	"fib/config"
	"fib/database"
	"fib/routers"
	"log"

	"github.com/gofiber/fiber/v2"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to an Awesome API")
}

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to the database
	database.ConnectDb()

	// Initialize Fiber app
	app := fiber.New()

	// Setup routes
	app.Get("/", welcome)
	routers.SetupUserRoutes(app)

	// Start the server
	log.Printf("Server is running on port %s", config.AppConfig.Port)
	log.Fatal(app.Listen(":" + config.AppConfig.Port))
}
