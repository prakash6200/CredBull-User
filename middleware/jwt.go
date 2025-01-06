package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key") // Replace with your actual secret key

// GenerateJWT generates a JWT token for the user
func GenerateJWT(userID uint, firstName string) (string, error) {
	// Set claims
	claims := jwt.MapClaims{
		"sub":        userID,                                // Subject (User ID)
		"first_name": firstName,                             // First Name
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Expiry (24 hours)
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString(jwtSecret)
}

// JWTMiddleware is a middleware to check for valid JWT token in the request
func JWTMiddleware(c *fiber.Ctx) error {
	// Get the token from the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Missing or invalid Authorization header",
		})
	}

	// The token should be prefixed with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid Authorization header format",
		})
	}

	// Extract the token part
	tokenString := authHeader[len("Bearer "):]

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check if the token method is valid
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	// If there's an error parsing the token
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Invalid or expired token",
		})
	}

	// If valid, continue to the next handler
	return c.Next()
}
