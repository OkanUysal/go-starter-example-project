package models

import (
	"os"
)

// GetUserTableName returns the user table name from environment variable
func GetUserTableName() string {
	tableName := os.Getenv("USER_TABLE")
	if tableName == "" {
		return "example_user" // default fallback
	}
	return tableName
}
