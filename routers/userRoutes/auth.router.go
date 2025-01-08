package routers

import (
	controllers "fib/controllers/userControllers"
	validator "fib/validator/userValidator"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/user")

	userGroup.Post("/signup", validator.Signup(), controllers.Signup)
	userGroup.Post("/login", validator.Login(), controllers.Login)
}
