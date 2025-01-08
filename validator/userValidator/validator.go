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

func Signup() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the request body into the User model
		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Invalid request body",
			})
		}

		// Validate the fields
		errors := make(map[string]string)

		// Validate Name
		if len(strings.TrimSpace(user.Name)) < 5 {
			errors["name"] = "Name must be at least 5 characters long"
		}

		// Validate Email
		if user.Email == "" || !isValidEmail(user.Email) {
			errors["email"] = "Invalid email format"
		}

		// Validate Mobile
		if user.Mobile == "" || !isValidMobile(user.Mobile) {
			errors["mobile"] = "Invalid mobile number format"
		}

		// Validate Password
		if len(strings.TrimSpace(user.Password)) < 8 {
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

		// Store the validated user in Locals for use in the controller
		c.Locals("validatedUser", user)

		// Continue to the next handler
		return c.Next()
	}
}

func Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the request body into the User model
		loginRequest := new(models.User)
		if err := c.BodyParser(loginRequest); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Invalid request body",
			})
		}

		// Validate the fields
		errors := make(map[string]string)

		// Validate Email or Mobile
		if loginRequest.Email == "" && loginRequest.Mobile == "" {
			errors["credentials"] = "Email or Mobile is required"
		} else {
			if loginRequest.Email != "" && !isValidEmail(loginRequest.Email) {
				errors["email"] = "Invalid email format"
			}
			if loginRequest.Mobile != "" && !isValidMobile(loginRequest.Mobile) {
				errors["mobile"] = "Invalid mobile number format"
			}
		}

		// Validate Password
		if len(strings.TrimSpace(loginRequest.Password)) < 8 {
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
		c.Locals("validatedUser", loginRequest)

		// Continue to the next handler
		return c.Next()
	}
}
