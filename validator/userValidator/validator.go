package validator

import (
	"fib/models"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func isValidMobile(mobile string) bool {
	re := regexp.MustCompile(`^\d{10}$`)
	return re.MatchString(mobile)
}

func Singup(model *models.User) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the request body into the User model
		if err := c.BodyParser(model); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Invalid request body",
			})
		}

		// Validate the fields
		errors := make(map[string]string)

		// Validate Name
		if len(strings.TrimSpace(model.Name)) < 5 {
			errors["name"] = "Name must be at least 5 characters long"
		}

		// Validate Email
		if model.Email == "" || !isValidEmail(model.Email) {
			errors["email"] = "Invalid email format"
		}

		// Validate Mobile
		if model.Mobile == "" || !isValidMobile(model.Mobile) {
			errors["mobile"] = "Invalid mobile number format"
		}

		// Validate Password
		if len(strings.TrimSpace(model.Password)) < 8 {
			errors["password"] = "Password must be at least 8 characters long"
		}

		// If there are validation errors, return them
		if len(errors) > 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  false,
				"message": "Validation failed",
				"errors":  errors,
			})
		}

		// Store the validated model in Locals for use in the controller
		c.Locals("validatedUser", model)

		// Continue to the next handler
		return c.Next()
	}
}
