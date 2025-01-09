package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("asdfasqsdfgsdasdfasdfawqe") // Replace with your actual secret key

// GenerateJWT generates a JWT token for the user
func GenerateJWT(userID uint, name, role string) (string, error) {
	// Set claims
	claims := jwt.MapClaims{
		"sub":  userID,                                // Subject (User ID)
		"name": name,                                  // Name of the user
		"role": role,                                  // User role
		"iat":  time.Now().Unix(),                     // Issued at (current timestamp)
		"exp":  time.Now().Add(24 * time.Hour).Unix(), // Expiry (24 hours from now)
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

type CustomErrorHandler struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// Error implements the error interface for CustomErrorHandler
func (ceh *CustomErrorHandler) Error() string {
	return ceh.Message
}

// New creates a new instance of CustomErrorHandler
func (ceh *CustomErrorHandler) New(status bool, message string) *CustomErrorHandler {
	ceh.Status = status
	ceh.Message = message
	return ceh
}

// AlreadyExist creates a "409 Conflict" error
func (ceh *CustomErrorHandler) AlreadyExist(message string) *CustomErrorHandler {
	return &CustomErrorHandler{Status: false, Message: message}
}

// WrongCredentials creates a "401 Unauthorized" error
func (ceh *CustomErrorHandler) WrongCredentials(message string) *CustomErrorHandler {
	return &CustomErrorHandler{Status: false, Message: message}
}

// LowBalance creates a "402 Payment Required" error
func (ceh *CustomErrorHandler) LowBalance(message string) *CustomErrorHandler {
	return &CustomErrorHandler{Status: false, Message: message}
}

// UnAuthorized creates a "401 Unauthorized" error
func (ceh *CustomErrorHandler) UnAuthorized(message string) *CustomErrorHandler {
	return &CustomErrorHandler{Status: false, Message: message}
}

// NotAllowed creates a "403 Forbidden" error
func (ceh *CustomErrorHandler) NotAllowed(message string) *CustomErrorHandler {
	return &CustomErrorHandler{Status: false, Message: message}
}

// NotFound creates a "404 Not Found" error
func (ceh *CustomErrorHandler) NotFound(message string) *CustomErrorHandler {
	return &CustomErrorHandler{Status: false, Message: message}
}

// ServerError creates a "500 Internal Server Error" error
func (ceh *CustomErrorHandler) ServerError(message string) *CustomErrorHandler {
	return &CustomErrorHandler{Status: false, Message: message}
}

// HandleErrorResponse handles the error and returns the proper JSON response
func HandleErrorResponse(c *fiber.Ctx, err error) error {
	// If it's a CustomErrorHandler, use its message and status
	if customErr, ok := err.(*CustomErrorHandler); ok {
		return c.Status(getStatusCode(customErr.Message)).JSON(fiber.Map{
			"status":  customErr.Status,
			"message": customErr.Message,
			"data":    nil,
		})
	}

	// If it's not a CustomErrorHandler, return internal server error
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  false,
		"message": "Something went wrong",
		"data":    nil,
	})
}

// Helper function to get appropriate status code for the error message
func getStatusCode(message string) int {
	switch message {
	case "Already exists":
		return fiber.StatusConflict
	case "Wrong credentials":
		return fiber.StatusUnauthorized
	case "Low balance":
		return fiber.StatusPaymentRequired
	case "Unauthorized":
		return fiber.StatusUnauthorized
	case "Not allowed":
		return fiber.StatusForbidden
	case "Not found":
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
	}
}
