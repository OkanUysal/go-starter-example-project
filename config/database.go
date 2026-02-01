package config

import (
	"fmt"
	"log"
	"os"

	"github.com/OkanUysal/go-logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	DB     *gorm.DB
	Logger *logger.Logger
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		// Use standard log here as Logger is not yet initialized
		log.Println("No .env file found, using environment variables")
	}
}

// ConnectDatabase establishes a connection to the PostgreSQL database
func ConnectDatabase() error {
	// Get DATABASE_URL_LOCAL from environment
	dsn := os.Getenv("DATABASE_URL_LOCAL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL_LOCAL environment variable is not set")
	}

	// Connect to database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	Logger.Info("Database connection established successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// GetEnv returns environment variable value or default
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
