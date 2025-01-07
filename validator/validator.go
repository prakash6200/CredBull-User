package validator

import (
	"fib/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// UserValidator validates the request body for models.User
func UserValidator(model *models.User) fiber.Handler {
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

		if len(strings.TrimSpace(model.Name)) < 5 {
			errors["first_name"] = "First name must be at least 5 characters long"
		}
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
