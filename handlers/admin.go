package handlers

import (
	"context"

	"github.com/OkanUysal/go-logger"
	"github.com/OkanUysal/go-response"
	"github.com/OkanUysal/go-starter-example-project/auth"
	"github.com/OkanUysal/go-starter-example-project/config"
	"github.com/OkanUysal/go-starter-example-project/models"
	"github.com/gin-gonic/gin"
)

// AdminDashboard godoc
// @Summary Get admin dashboard data
// @Description Returns admin dashboard information - only accessible by admins
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Admin dashboard data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Admin access required"
// @Router /admin/dashboard [get]
func AdminDashboard(c *gin.Context) {
	userID, _ := auth.GetUserID(c)
	role, _ := auth.GetRole(c)

	cache := config.GetCache()
	ctx := context.Background()
	cacheKey := "admin:dashboard:stats"

	// Define stats structure
	type DashboardStats struct {
		TotalUsers int64 `json:"total_users"`
		AdminCount int64 `json:"admin_count"`
		GuestCount int64 `json:"guest_count"`
	}

	var stats DashboardStats

	// Try to get from cache first
	if err := cache.GetJSON(ctx, cacheKey, &stats); err != nil {
		// Cache miss - get from database
		config.Logger.Info("Cache miss: admin dashboard stats")
		db := config.GetDB()
		db.Model(&models.User{}).Count(&stats.TotalUsers)
		db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&stats.AdminCount)
		db.Model(&models.User{}).Where("is_guest = ?", true).Count(&stats.GuestCount)

		// Cache for default TTL (5 minutes)
		cache.SetJSON(ctx, cacheKey, stats)
	} else {
		config.Logger.Info("Cache hit: admin dashboard stats",
			logger.Int64("total_users", stats.TotalUsers),
			logger.Int64("admin_count", stats.AdminCount),
			logger.Int64("guest_count", stats.GuestCount))
	}

	response.Success(c, gin.H{
		"admin": gin.H{
			"user_id": userID,
			"role":    role,
		},
		"statistics": stats,
	}, "Admin dashboard data")
}

// ListUsers godoc
// @Summary List all users
// @Description Returns a list of all users in the system - only accessible by admins
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of users"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Admin access required"
// @Router /admin/users [get]
func ListUsers(c *gin.Context) {
	db := config.GetDB()

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"users": users,
		"count": len(users),
	})
}
