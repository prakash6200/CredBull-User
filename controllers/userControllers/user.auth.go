package controllers

import (
	"crypto/rand"
	"fib/database"
	jwt "fib/middleware"
	"fib/models"
	"log"
	"math/big"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func generateReferralCode() string {
	const charSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	// Function to check if the referral code already exists in the database
	checkReferralCodeExists := func(code string) bool {
		var user models.User
		result := database.Database.Db.Where("referral_code = ?", code).First(&user)
		return result.RowsAffected > 0
	}

	for {
		// Create a slice to hold the generated characters
		code := make([]byte, length)

		// Generate random characters
		for i := 0; i < length; i++ {
			// Generate a random index within the charSet
			randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
			if err != nil {
				log.Fatal("Failed to generate random number:", err)
			}

			// Assign the random character to the code slice
			code[i] = charSet[randomIndex.Int64()]
		}

		// Convert byte slice to string
		referralCode := string(code)

		// Check if the generated referral code already exists in the database
		if !checkReferralCodeExists(referralCode) {
			return referralCode // Return the code if it doesn't exist
		}
	}
}

func Signup(c *fiber.Ctx) error {
	// Parse the request body
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid request body",
		})
	}

	// Check if email already exists
	existingUser := models.User{}
	result := database.Database.Db.Where("email = ?", user.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  false,
			"message": "Email is already registered",
		})
	}

	// Check if mobile already exists
	existingUserByMobile := models.User{}
	result = database.Database.Db.Where("mobile = ?", user.Mobile).First(&existingUserByMobile)
	if result.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  false,
			"message": "Mobile number is already registered",
		})
	}

	user.ReferralCode = generateReferralCode()

	// Hash the password using bcrypt with a cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to process your request",
		})
	}
	user.Password = string(hashedPassword)

	// Save the user to the database
	if err := database.Database.Db.Create(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to create user",
		})
	}

	// Remove the password from the response for security
	user.Password = ""
	user.UserKYC = models.UserKYC{}

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
	if loginRequest.Email == "" && loginRequest.Mobile == "" || loginRequest.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Email or Mobile and password are required",
		})
	}

	// Retrieve user from database by email or mobile
	var user models.User
	var result *gorm.DB
	if loginRequest.Email != "" {
		result = database.Database.Db.Where("email = ?", loginRequest.Email).First(&user)
	} else {
		result = database.Database.Db.Where("mobile = ?", loginRequest.Mobile).First(&user)
	}

	if result.Error != nil {
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
	token, err := jwt.GenerateJWT(user.ID, user.Name)
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

func Sdfa(c *fiber.Ctx) error {
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
