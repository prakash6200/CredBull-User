package routers

import (
	controllers "fib/controllers/auth"
	validators "fib/validators/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/user")

	userGroup.Post("/signup", validators.Signup(), controllers.Signup)
	userGroup.Post("/login", validators.Login(), controllers.Login)
	userGroup.Post("/send/otp", validators.SendOTP(), controllers.SendOTP)
}
