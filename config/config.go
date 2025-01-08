package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port   string
	DBName string
	JWTKey string
}

// AppConfig is a global variable to access configuration
var AppConfig *Config

// LoadConfig initializes configuration from environment variables or defaults
func LoadConfig() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using system environment variables.")
	}

	// Initialize AppConfig with values from environment variables
	AppConfig = &Config{
		Port:   getEnv("PORT", "3000"),
		DBName: getEnv("DB_NAME", "credUser.db"),
		JWTKey: getEnv("JWT_SECRET_KEY", "defaultSecret"),
	}

	// Validate critical configuration
	if AppConfig.JWTKey == "defaultSecret" {
		log.Println("Warning: Using default JWT_SECRET_KEY. Update it in your environment.")
	}
	if AppConfig.DBName == "credUser.db" {
		log.Println("Warning: Using default DBName. Update it in your environment.")
	}

}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
