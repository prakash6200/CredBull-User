package controllers

import (
	"crypto/rand"
	"fib/database"
	"fib/middleware"
	"fib/models"
	"log"
	"math/big"
	"time"

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
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return middleware.JsonResponse(c, fiber.StatusBadRequest, false, "Invalid request body!", nil)
	}

	// Check if email already exists
	existingUser := models.User{}
	result := database.Database.Db.Where("email = ?", user.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		return middleware.JsonResponse(c, fiber.StatusConflict, false, "Email is already registered!", nil)
	}

	// Check if mobile already exists
	existingUserByMobile := models.User{}
	result = database.Database.Db.Where("mobile = ?", user.Mobile).First(&existingUserByMobile)
	if result.RowsAffected > 0 {
		return middleware.JsonResponse(c, fiber.StatusConflict, false, "Mobile number is already registered!", nil)
	}

	user.ReferralCode = generateReferralCode()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return middleware.JsonResponse(c, fiber.StatusInternalServerError, false, "Failed to process your request!", nil)
	}
	user.Password = string(hashedPassword)

	if err := database.Database.Db.Create(user).Error; err != nil {
		log.Printf("Error saving user to database: %v", err)
		return middleware.JsonResponse(c, fiber.StatusInternalServerError, false, "Failed to Signup user!", nil)
	}

	user.Password = ""
	user.UserKYC = models.UserKYC{}

	return middleware.JsonResponse(c, fiber.StatusCreated, true, "User registered successfully.", user)
}

func Login(c *fiber.Ctx) error {
	reqData := new(struct {
		Mobile   string `json:"mobile"`
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	if err := c.BodyParser(reqData); err != nil {
		return middleware.JsonResponse(c, fiber.StatusBadRequest, false, "Failed to parse request body!", nil)
	}

	if reqData.Password == "" {
		return middleware.JsonResponse(c, fiber.StatusBadRequest, false, "Password is required!", nil)
	}

	if reqData.Email == "" && reqData.Mobile == "" {
		return middleware.JsonResponse(c, fiber.StatusBadRequest, false, "Either email or mobile number is required!", nil)
	}

	var user models.User
	var result *gorm.DB

	// Retrieve user by email or mobile
	if reqData.Email != "" {
		result = database.Database.Db.Where("email = ? AND is_deleted = ?", reqData.Email, false).First(&user)
	} else {
		result = database.Database.Db.Where("mobile = ? AND is_deleted = ?", reqData.Mobile, false).First(&user)
	}

	if result.Error != nil {
		return middleware.JsonResponse(c, fiber.StatusUnauthorized, false, "Invalid credentials!", nil)
	}

	if !user.IsEmailVerified {
		return middleware.JsonResponse(c, fiber.StatusUnauthorized, false, "Email not verified!", nil)
	}

	if !user.IsMobileVerified {
		return middleware.JsonResponse(c, fiber.StatusUnauthorized, false, "Mobile not verified!", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqData.Password)); err != nil {
		return middleware.JsonResponse(c, fiber.StatusUnauthorized, false, "Wrong password!", nil)
	}

	// Update last login time
	user.LastLogin = time.Now()
	if err := database.Database.Db.Save(&user).Error; err != nil {
		log.Printf("Error saving last login time: %v", err)
	}

	// Sanitize user data (remove sensitive fields)
	user.Password = ""
	user.ProfileImage = ""

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Name, user.Role)
	if err != nil {
		return middleware.JsonResponse(c, fiber.StatusInternalServerError, false, "Error generating JWT token!", nil)
	}

	return middleware.JsonResponse(c, fiber.StatusOK, true, "Login successful.", fiber.Map{
		"user":  user,
		"token": token,
	})
}