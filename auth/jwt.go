package auth

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrExpiredToken    = errors.New("token has expired")
	ErrNotRefreshToken = errors.New("token is not a refresh token")
)

// TokenType represents the type of token
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// Claims represents JWT claims
type Claims struct {
	UserID    string    `json:"user_id"`
	Role      string    `json:"role,omitempty"`
	TokenType TokenType `json:"token_type"`
	FamilyID  string    `json:"family_id"` // Links access and refresh tokens together
	jwt.RegisteredClaims
}

// getTokenDuration returns the token duration from environment or default
func getTokenDuration(envKey string, defaultHours int) time.Duration {
	durationStr := os.Getenv(envKey)
	if durationStr == "" {
		return time.Duration(defaultHours) * time.Hour
	}

	hours, err := strconv.Atoi(durationStr)
	if err != nil {
		return time.Duration(defaultHours) * time.Hour
	}

	return time.Duration(hours) * time.Hour
}

// TokenPair represents an access and refresh token pair
type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessTokenJTI   string
	RefreshTokenJTI  string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

// GenerateTokenPair generates a new access and refresh token pair with the same family ID
func GenerateTokenPair(userID, role string) (*TokenPair, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key"
	}

	familyID := uuid.New().String()

	// Generate access token
	accessDuration := getTokenDuration("ACCESS_TOKEN_DURATION", 24)
	accessExpirationTime := time.Now().Add(accessDuration)
	accessJTI := uuid.New().String()

	accessClaims := &Claims{
		UserID:    userID,
		Role:      role,
		TokenType: TokenTypeAccess,
		FamilyID:  familyID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        accessJTI,
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshDuration := getTokenDuration("REFRESH_TOKEN_DURATION", 168)
	refreshExpirationTime := time.Now().Add(refreshDuration)
	refreshJTI := uuid.New().String()

	refreshClaims := &Claims{
		UserID:    userID,
		TokenType: TokenTypeRefresh,
		FamilyID:  familyID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        refreshJTI,
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      accessTokenString,
		RefreshToken:     refreshTokenString,
		AccessTokenJTI:   accessJTI,
		RefreshTokenJTI:  refreshJTI,
		AccessExpiresAt:  accessExpirationTime,
		RefreshExpiresAt: refreshExpirationTime,
	}, nil
}

// Legacy functions for backward compatibility
// GenerateAccessToken generates a new access token
func GenerateAccessToken(userID, role string) (string, error) {
	pair, err := GenerateTokenPair(userID, role)
	if err != nil {
		return "", err
	}
	return pair.AccessToken, nil
}

// GenerateRefreshToken generates a new refresh token
func GenerateRefreshToken(userID string) (string, error) {
	pair, err := GenerateTokenPair(userID, "")
	if err != nil {
		return "", err
	}
	return pair.RefreshToken, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key"
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token specifically
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != TokenTypeRefresh {
		return nil, ErrNotRefreshToken
	}

	return claims, nil
}
