package routers

import (
	controllers "fib/controllers/auth"
	"fib/middleware"
	validators "fib/validators/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/user")

	userGroup.Post("/signup", validators.Signup(), controllers.Signup)
	userGroup.Post("/login", validators.Login(), controllers.Login)
	userGroup.Post("/send/otp", validators.SendOTP(), controllers.SendOTP)
	userGroup.Patch("/verify/otp", validators.VerifyOTP(), controllers.VerifyOTP)
	userGroup.Post("/forgot/password/send/otp", validators.SendOTP(), controllers.ForgotPasswordSendOTP)
	userGroup.Patch("/forgot/password/verify/otp", validators.VerifyOTP(), controllers.ForgotPasswordVerifyOTP)
	userGroup.Patch("/reset/password", validators.ResetPassword(), middleware.JWTMiddleware, controllers.ResetPassword)
}
