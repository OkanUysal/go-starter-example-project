package auth

import (
	"time"

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

	blacklist := models.TokenBlacklist{
		JTI:       familyID, // Using family ID as JTI for family-based blacklist
		UserID:    userID,
		FamilyID:  &familyID,
		ExpiresAt: expiresAt,
	}

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
	db := config.GetDB()

	var count int64
	db.Model(&models.TokenBlacklist{}).Where("jti = ?", jti).Count(&count)

	return count > 0
}

// IsTokenFamilyBlacklisted checks if a token's family is blacklisted
func IsTokenFamilyBlacklisted(familyID string) bool {
	db := config.GetDB()

	var count int64
	db.Model(&models.TokenBlacklist{}).Where("family_id = ?", familyID).Count(&count)

	return count > 0
} // CleanupExpiredTokens removes expired tokens from blacklist
func CleanupExpiredTokens() error {
	db := config.GetDB()
	return db.Where("expires_at < ?", time.Now()).Delete(&models.TokenBlacklist{}).Error
}
