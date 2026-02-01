package handlers

import (
	"github.com/OkanUysal/go-response"
	"github.com/gin-gonic/gin"
)

// HelloResponse represents the response for hello endpoint
type HelloResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// HelloHandler handles the hello endpoint
// @Summary Hello endpoint
// @Description Returns a hello message
// @Tags general
// @Accept json
// @Produce json
// @Success 200 {object} HelloResponse
// @Router /hello [get]
func HelloHandler(c *gin.Context) {
	data := HelloResponse{
		Message: "Hello, World!",
		Status:  "success",
	}
	response.Success(c, data)
}
