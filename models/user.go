package models

import (
	"time"
)

// UserRole represents the role of a user
type UserRole string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)

// User represents a user in the system
type User struct {
	ID          string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	GuestID     *string   `json:"guest_id,omitempty" gorm:"type:varchar(255);uniqueIndex"`
	GoogleID    *string   `json:"google_id,omitempty" gorm:"type:varchar(255);uniqueIndex"`
	DisplayName string    `json:"display_name" gorm:"type:varchar(255);not null"`
	Role        UserRole  `json:"role" gorm:"type:varchar(50);not null;default:'USER'"`
	IsGuest     bool      `json:"is_guest" gorm:"not null;default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for User model based on environment variable
func (User) TableName() string {
	return GetUserTableName()
}
