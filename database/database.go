package database

import (
	"log"
	"os"

	"fib/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {

	dbName := "test.db"

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database '%s': %v", dbName, err)
		os.Exit(2)
	}

	log.Println("Connected Successfully to Database")

	// Set GORM logger
	db.Logger = logger.Default.LogMode(logger.Info)

	// Run migrations
	runMigrations(db)

	// Save database instance
	Database = DbInstance{
		Db: db,
	}
}

func runMigrations(db *gorm.DB) {
	log.Println("Running Migrations")
	if err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
