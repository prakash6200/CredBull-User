package validator

import (
	"fib/middleware"
	"fib/models"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Helper to validate email format
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Helper to validate mobile number format
func isValidMobile(mobile string) bool {
	re := regexp.MustCompile(`^\d{10}$`)
	return re.MatchString(mobile)
}

// Signup validator middleware
func Signup() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return middleware.JsonResponse(c, fiber.StatusBadRequest, false, "Invalid request body!", nil)
		}

		errors := make(map[string]string)

		// Validate Name
		if len(strings.TrimSpace(user.Name)) < 5 {
			errors["name"] = "Name must be at least 5 characters long!"
		}

		// Validate Email
		if user.Email == "" || !isValidEmail(user.Email) {
			errors["email"] = "Invalid email!"
		}

		// Validate Mobile
		if user.Mobile == "" || !isValidMobile(user.Mobile) {
			errors["mobile"] = "Invalid mobile number!"
		}

		// Validate Password
		if len(strings.TrimSpace(user.Password)) < 8 {
			errors["password"] = "Password must be at least 8 characters long!"
		}

		// Respond with errors if any exist
		if len(errors) > 0 {
			return middleware.ValidationErrorResponse(c, errors)
		}

		// Pass validated user to the next middleware
		c.Locals("validatedUser", user)
		return c.Next()
	}
}

// Login validator middleware
func Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		loginRequest := new(models.User)
		if err := c.BodyParser(loginRequest); err != nil {
			return middleware.JsonResponse(c, fiber.StatusBadRequest, false, "Invalid request body!", nil)
		}

		errors := make(map[string]string)

		// Validate credentials
		if loginRequest.Email == "" && loginRequest.Mobile == "" {
			errors["credentials"] = "Either email or mobile number is required!"
		} else {
			if loginRequest.Email != "" && !isValidEmail(loginRequest.Email) {
				errors["email"] = "Invalid email!"
			}
			if loginRequest.Mobile != "" && !isValidMobile(loginRequest.Mobile) {
				errors["mobile"] = "Invalid mobile number!"
			}
		}

		// Validate Password
		if len(strings.TrimSpace(loginRequest.Password)) < 8 {
			errors["password"] = "Password must be at least 8 characters long!"
		}

		// Respond with errors if any exist
		if len(errors) > 0 {
			return middleware.ValidationErrorResponse(c, errors)
		}

		// Pass validated login request to the next middleware
		c.Locals("validatedUser", loginRequest)
		return c.Next()
	}
}
