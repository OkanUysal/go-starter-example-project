package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/OkanUysal/go-logger"
	"github.com/OkanUysal/go-starter-example-project/config"
	"github.com/OkanUysal/go-starter-example-project/models"
)

// BlacklistToken adds a token to the blacklist
func BlacklistToken(jti, userID string, expiresAt time.Time) error {
	db := config.GetDB()

	blacklist := models.TokenBlacklist{
		JTI:       jti,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	return db.Create(&blacklist).Error
}

// BlacklistByFamilyID blacklists all tokens in a family
func BlacklistByFamilyID(familyID, userID string, expiresAt time.Time) error {
	db := config.GetDB()
	cache := config.GetCache()
	ctx := context.Background()

	blacklist := models.TokenBlacklist{
		JTI:       familyID, // Using family ID as JTI for family-based blacklist
		UserID:    userID,
		FamilyID:  &familyID,
		ExpiresAt: expiresAt,
	}

	// Invalidate cache
	cache.Delete(ctx, fmt.Sprintf("blacklist:family:%s", familyID))

	return db.Create(&blacklist).Error
}

// BlacklistTokenPair adds both access and refresh tokens to the blacklist
func BlacklistTokenPair(accessJTI, refreshJTI, userID string, accessExpiresAt, refreshExpiresAt time.Time) error {
	db := config.GetDB()

	blacklists := []models.TokenBlacklist{
		{
			JTI:       accessJTI,
			UserID:    userID,
			ExpiresAt: accessExpiresAt,
		},
		{
			JTI:       refreshJTI,
			UserID:    userID,
			ExpiresAt: refreshExpiresAt,
		},
	}

	return db.Create(&blacklists).Error
}

// IsTokenBlacklisted checks if a token is blacklisted (by JTI or family ID)
func IsTokenBlacklisted(jti string) bool {
	cache := config.GetCache()
	ctx := context.Background()
	cacheKey := fmt.Sprintf("blacklist:jti:%s", jti)

	// Check cache first
	var cached bool
	if err := cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		config.Logger.Info("Cache hit: token blacklist check", logger.String("jti", jti), logger.Bool("is_blacklisted", cached))
		return cached
	}

	config.Logger.Info("Cache miss: token blacklist check", logger.String("jti", jti))

	// Check database
	db := config.GetDB()
	var count int64
	db.Model(&models.TokenBlacklist{}).Where("jti = ?", jti).Count(&count)

	isBlacklisted := count > 0

	// Cache the result (uses default TTL: 5 minutes)
	cache.SetJSON(ctx, cacheKey, isBlacklisted)

	return isBlacklisted
}

// IsTokenFamilyBlacklisted checks if a token's family is blacklisted
func IsTokenFamilyBlacklisted(familyID string) bool {
	cache := config.GetCache()
	ctx := context.Background()
	cacheKey := fmt.Sprintf("blacklist:family:%s", familyID)

	// Check cache first
	var cached bool
	if err := cache.GetJSON(ctx, cacheKey, &cached); err == nil {
		config.Logger.Info("Cache hit: token family blacklist check", logger.String("family_id", familyID), logger.Bool("is_blacklisted", cached))
		return cached
	}

	config.Logger.Info("Cache miss: token family blacklist check", logger.String("family_id", familyID))

	// Check database
	db := config.GetDB()
	var count int64
	db.Model(&models.TokenBlacklist{}).Where("family_id = ?", familyID).Count(&count)

	isBlacklisted := count > 0

	// Cache the result (uses default TTL: 5 minutes)
	cache.SetJSON(ctx, cacheKey, isBlacklisted)

	return isBlacklisted
} // CleanupExpiredTokens removes expired tokens from blacklist
func CleanupExpiredTokens() error {
	db := config.GetDB()
	return db.Where("expires_at < ?", time.Now()).Delete(&models.TokenBlacklist{}).Error
}
