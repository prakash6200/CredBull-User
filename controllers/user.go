package controllers

import (
	"fib/database"
	"fib/middleware"
	"fib/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	// Parse incoming request body
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	// Validate required fields
	if user.FirstName == "" || user.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "First name and last name are required",
		})
	}

	// Save user to database
	if result := database.Database.Db.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Return the created user
	return c.Status(fiber.StatusCreated).JSON(user)
}

func SignUp(c *fiber.Ctx) error {
	// Parse incoming request body
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to parse request",
		})
	}

	// Validate required fields
	if user.FirstName == "" || user.LastName == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "First name, last name, and password are required",
		})
	}

	// Hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Error hashing password",
		})
	}
	user.Password = string(hashedPassword)

	// Save user to database
	if result := database.Database.Db.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to create user",
		})
	}

	// Return the created user
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  true,
		"message": "User created successfully",
		"user":    user,
	})
}

func Login(c *fiber.Ctx) error {
	// Parse incoming request body
	loginRequest := new(models.User)
	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to parse request",
		})
	}

	// Validate required fields
	if loginRequest.FirstName == "" || loginRequest.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "First name and password are required",
		})
	}

	// Retrieve user from database
	var user models.User
	if result := database.Database.Db.Where("first_name = ?", loginRequest.FirstName).First(&user); result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid credentials",
		})
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid credentials",
		})
	}

	// Generate JWT token
	token, err := jwt.GenerateJWT(user.ID, user.FirstName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Error generating JWT token",
		})
	}

	// Successful login, return the token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}

func GetUserProfile(c *fiber.Ctx) error {
	// Get the user ID from the JWT claims (extracted in the JWT middleware)
	userID := c.Locals("user").(map[string]interface{})["sub"].(float64)

	// Retrieve user from the database by user ID
	var user models.User
	if result := database.Database.Db.First(&user, userID); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "User not found",
		})
	}

	// Return the user profile
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"user":   user,
	})
}
