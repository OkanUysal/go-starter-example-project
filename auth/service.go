package auth

import (
	"fmt"
	"math/rand"

	"github.com/OkanUysal/go-starter-example-project/config"
	"github.com/OkanUysal/go-starter-example-project/models"
	"github.com/google/uuid"
)

// Service handles authentication operations
type Service struct{}

// NewService creates a new auth service
func NewService() *Service {
	return &Service{}
}

// GuestLoginResponse represents the response for guest login
type GuestLoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         models.User `json:"user"`
}

// RefreshTokenRequest represents the request for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// GuestLoginRequest represents the request for guest login
type GuestLoginRequest struct {
	GuestID *string `json:"guest_id,omitempty"`
}

// GuestLogin creates a new guest user or logs in existing guest and returns tokens
func (s *Service) GuestLogin(guestID *string) (*GuestLoginResponse, error) {
	db := config.GetDB()
	var user models.User

	// If guest_id is provided, try to find existing user
	if guestID != nil && *guestID != "" {
		if err := db.Where("guest_id = ?", *guestID).First(&user).Error; err == nil {
			// User found, generate new token pair
			tokenPair, err := GenerateTokenPair(user.ID, string(user.Role))
			if err != nil {
				return nil, fmt.Errorf("failed to generate tokens: %w", err)
			}

			return &GuestLoginResponse{
				AccessToken:  tokenPair.AccessToken,
				RefreshToken: tokenPair.RefreshToken,
				User:         user,
			}, nil
		}
		// If not found, continue to create new user
	}

	// Generate unique guest ID if not provided
	newGuestID := uuid.New().String()
	if guestID != nil && *guestID != "" {
		newGuestID = *guestID
	}

	userID := uuid.New().String()

	// Generate random display name
	randomNum := rand.Intn(9000) + 1000
	displayName := fmt.Sprintf("Guest%d", randomNum)

	// Create new user
	user = models.User{
		ID:          userID,
		GuestID:     &newGuestID,
		DisplayName: displayName,
		Role:        models.RoleUser,
		IsGuest:     true,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate token pair
	tokenPair, err := GenerateTokenPair(user.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &GuestLoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User:         user,
	}, nil
}

// RefreshToken validates refresh token only and issues new tokens, blacklisting the old token family
func (s *Service) RefreshToken(token string) (*GuestLoginResponse, error) {
	// Validate refresh token specifically
	claims, err := ValidateRefreshToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if token or its family is blacklisted
	if IsTokenBlacklisted(claims.ID) || IsTokenFamilyBlacklisted(claims.FamilyID) {
		return nil, fmt.Errorf("token has been revoked")
	}

	// Get user from database
	db := config.GetDB()
	var user models.User
	if err := db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate new token pair
	tokenPair, err := GenerateTokenPair(user.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Blacklist the entire old token family (both access and refresh tokens)
	if claims.ExpiresAt != nil {
		if err := BlacklistByFamilyID(claims.FamilyID, user.ID, claims.ExpiresAt.Time); err != nil {
			return nil, fmt.Errorf("failed to blacklist token family: %w", err)
		}
	}

	return &GuestLoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User:         user,
	}, nil
}

// GetUserByID returns a user by ID
func (s *Service) GetUserByID(userID string) (*models.User, error) {
	db := config.GetDB()
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}
