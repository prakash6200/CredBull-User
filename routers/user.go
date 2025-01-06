package routers

import (
	"fib/controllers"
	"fib/middleware"
	"fib/models"
	"fib/validator"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/users")

	// Attach the validator middleware before the controller
	userGroup.Post("/", validator.UserValidator(&models.User{}), controllers.CreateUser)

	userGroup.Post("/signup", controllers.SignUp)
	userGroup.Post("/login", controllers.Login)
	userGroup.Get("/profile", jwt.JWTMiddleware, controllers.GetUserProfile)
}
