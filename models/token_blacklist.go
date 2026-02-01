package models

import (
	"os"
	"time"
)

// TokenBlacklist represents a blacklisted token
type TokenBlacklist struct {
	JTI       string    `json:"jti" gorm:"primaryKey;type:varchar(255)"`
	UserID    string    `json:"user_id" gorm:"type:varchar(255);not null;index"`
	FamilyID  *string   `json:"family_id,omitempty" gorm:"type:varchar(255);index"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name from environment variable
func (TokenBlacklist) TableName() string {
	tableName := os.Getenv("TOKEN_BLACKLIST_TABLE")
	if tableName == "" {
		return "example_token_blacklist" // default fallback
	}
	return tableName
}
