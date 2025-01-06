package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Config struct to hold all the environment variables
type Config struct {
	Port             string
	GRPCUrl          string
	GRPCUserUrl      string
	GRPCWalletUrl    string
	RedisUrl         string
	DatabaseUri      string
	FrontUrl         string
	BaseUrl          string
	AdminUrl         string
	BinanceApiKey    string
	BinanceSecretKey string
	CoinDCXApiKey    string
	CoinDCXSecretKey string
	JwtSecretKey     string
}

// Global variable to hold the configuration
var AppConfig *Config

// LoadConfig loads environment variables from the .env file
func LoadConfig() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize AppConfig with values from environment variables
	AppConfig = &Config{
		Port:             os.Getenv("PORT"),
		GRPCUrl:          os.Getenv("GRPC_URL"),
		GRPCUserUrl:      os.Getenv("GRPC_USER_URL"),
		GRPCWalletUrl:    os.Getenv("GRPC_WALLET_URL"),
		RedisUrl:         os.Getenv("REDIS_URL"),
		DatabaseUri:      os.Getenv("DATABASE_URI"),
		FrontUrl:         os.Getenv("FRONT_URL"),
		BaseUrl:          os.Getenv("BASE_URL"),
		AdminUrl:         os.Getenv("ADMIN_URL"),
		BinanceApiKey:    os.Getenv("BINANCE_API_KEY"),
		BinanceSecretKey: os.Getenv("BINANCE_SECRET_KEY"),
		CoinDCXApiKey:    os.Getenv("COINDCX_API_KEY"),
		CoinDCXSecretKey: os.Getenv("COINDCX_SECRET_KEY"),
		JwtSecretKey:     os.Getenv("JWT_SECRET_KEY"),
	}
}
