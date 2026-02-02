package config

import (
	"os"
	"strconv"
)

// Feature flags
var (
	// RoomAuthEnabled enables room authorization feature
	// When enabled, users need explicit permission from admin to join game rooms
	RoomAuthEnabled bool
)

// LoadConfig loads configuration from environment variables
func LoadConfig() {
	// Load room authorization feature flag (default: false)
	RoomAuthEnabled = getEnvBool("ROOM_AUTH_ENABLED", false)
}

// getEnvBool gets boolean value from environment variable
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
