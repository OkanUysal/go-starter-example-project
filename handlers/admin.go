package handlers

import (
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

	db := config.GetDB()

	// Get total user count
	var totalUsers int64
	db.Model(&models.User{}).Count(&totalUsers)

	// Get admin count
	var adminCount int64
	db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&adminCount)

	// Get guest count
	var guestCount int64
	db.Model(&models.User{}).Where("is_guest = ?", true).Count(&guestCount)

	response.Success(c, gin.H{
		"admin": gin.H{
			"user_id": userID,
			"role":    role,
		},
		"statistics": gin.H{
			"total_users": totalUsers,
			"admin_count": adminCount,
			"guest_count": guestCount,
		},
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
