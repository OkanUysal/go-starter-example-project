package handlers

import (
	"github.com/OkanUysal/go-response"
	"github.com/OkanUysal/go-starter-example-project/auth"
	"github.com/gin-gonic/gin"
)

var authService *auth.Service

// InitAuthService initializes the auth service
func InitAuthService() error {
	authService = auth.NewService()
	return nil
}

// GuestLogin handles guest login
// @Summary Guest login
// @Description Creates a new guest user or logs in existing guest with guest_id
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.GuestLoginRequest false "Guest login request (optional guest_id)"
// @Success 200 {object} auth.GuestLoginResponse
// @Failure 500 {object} map[string]string
// @Router /auth/guest-login [post]
func GuestLogin(c *gin.Context) {
	var req auth.GuestLoginRequest
	// Bind JSON but don't require it
	_ = c.ShouldBindJSON(&req)

	result, err := authService.GuestLogin(req.GuestID)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, result, "Guest login successful")
}

// RefreshToken handles token refresh
// @Summary Refresh token
// @Description Validates refresh token only and issues new tokens. Old refresh token will be blacklisted.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.RefreshTokenRequest true "Refresh token only"
// @Success 200 {object} auth.GuestLoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var req auth.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "INVALID_REQUEST", "Invalid request body")
		return
	}

	result, err := authService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.Success(c, result, "Token refreshed successfully")
}

// GetMe returns the current authenticated user
// @Summary Get current user
// @Description Returns the current authenticated user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Router /auth/me [get]
func GetMe(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := authService.GetUserByID(userID)
	if err != nil {
		response.NotFound(c, "User")
		return
	}

	response.Success(c, user)
}
