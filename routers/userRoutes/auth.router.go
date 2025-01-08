package routers

import (
	controllers "fib/controllers/userControllers"
	"fib/models"
	validator "fib/validator/userValidator"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/user")

	// Attach the validator middleware before the controller
	userGroup.Post("/signup", validator.Singup(&models.User{}), controllers.Signup)

	// userGroup.Post("/signup", controllers.SignUp)
	// userGroup.Post("/login", controllers.Login)
	// userGroup.Get("/profile", jwt.JWTMiddleware, controllers.GetUserProfile)
}
