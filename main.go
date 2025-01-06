package main

import (
	"fib/database"
	"fib/routers"
	"log"

	"github.com/gofiber/fiber/v2"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to an Awesome API")
}

func main() {
	database.ConnectDb()

	app := fiber.New()
	app.Get("/", welcome)

	// Setup routes
	routers.SetupUserRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
