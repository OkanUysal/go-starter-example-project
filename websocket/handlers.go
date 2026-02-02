package websocket

import (
	"fmt"
	"time"

	"github.com/OkanUysal/go-logger"
	"github.com/OkanUysal/go-response"
	"github.com/OkanUysal/go-starter-example-project/auth"
	"github.com/OkanUysal/go-starter-example-project/config"
	gowebsocket "github.com/OkanUysal/go-websocket"
	"github.com/gin-gonic/gin"
)

// WebSocketConnect handles WebSocket connection
// @Summary Connect to WebSocket
// @Description Establish WebSocket connection for real-time communication. Requires authentication via Bearer token in header OR token query parameter.
// @Tags websocket
// @Security BearerAuth
// @Param room_id query string true "Room ID to join (use 'lobby' for public lobby)"
// @Param token query string false "JWT token (alternative to Authorization header for WebSocket connections)"
// @Success 101 "Switching Protocols"
// @Failure 401 {object} map[string]string "Unauthorized - Token required"
// @Failure 404 {object} map[string]string "Room not found"
// @Router /ws [get]
func WebSocketConnect(c *gin.Context) {
	roomID := c.Query("room_id")
	if roomID == "" {
		roomID = LobbyRoomID
	}

	// Get user from context
	userID, exists := auth.GetUserID(c)
	if !exists {
		response.Error(c, 401, "Unauthorized", nil)
		return
	}

	// Get user info for username
	username := userID // Default to userID
	authService := auth.NewService()
	if user, err := authService.GetUserByID(userID); err == nil {
		if user.GuestID != nil {
			username = "Guest_" + (*user.GuestID)[:8]
		}
	}

	manager := GetRoomManager()

	// Check if room exists
	_, err := manager.GetRoom(roomID)
	if err != nil {
		response.Error(c, 404, "Room not found", nil)
		return
	}

	// Handle WebSocket connection - go-websocket HandleConnection signature: (hub, w, r, userID)
	// Note: This upgrades the HTTP connection to WebSocket, no response should be sent after this
	// The client will join the room after connection by sending a join message
	err = gowebsocket.HandleConnection(manager.GetHub(), c.Writer, c.Request, userID)
	if err != nil {
		config.Logger.Error("WebSocket connection failed",
			logger.Err(err),
			logger.String("user_id", userID),
			logger.String("room_id", roomID))
		// Don't send response after WebSocket upgrade attempt
		return
	}

	// After connection is established, automatically join the requested room
	// We'll do this via a goroutine to avoid blocking
	go func() {
		defer func() {
			if r := recover(); r != nil {
				config.Logger.Error("Panic in auto-join",
					logger.String("error", fmt.Sprint(r)),
					logger.String("user_id", userID),
					logger.String("room_id", roomID))
			}
		}()

		time.Sleep(100 * time.Millisecond) // Give connection time to establish
		if joinErr := manager.JoinRoom(roomID, userID, username); joinErr != nil {
			config.Logger.Error("Failed to auto-join room after connection",
				logger.Err(joinErr),
				logger.String("user_id", userID),
				logger.String("room_id", roomID))
		}
	}()
}

// GetRooms returns all active rooms
// @Summary Get all rooms
// @Description Get list of all active WebSocket rooms
// @Tags websocket
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of rooms"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /ws/rooms [get]
func GetRooms(c *gin.Context) {
	manager := GetRoomManager()
	rooms := manager.GetAllRooms()

	response.Success(c, gin.H{
		"rooms": rooms,
		"total": len(rooms),
	}, "Rooms retrieved successfully")
}

// GetRoomInfo returns information about a specific room
// @Summary Get room information
// @Description Get detailed information about a specific room
// @Tags websocket
// @Security BearerAuth
// @Param room_id path string true "Room ID"
// @Success 200 {object} map[string]interface{} "Room information"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Room not found"
// @Router /ws/rooms/{room_id} [get]
func GetRoomInfo(c *gin.Context) {
	roomID := c.Param("room_id")

	manager := GetRoomManager()
	room, err := manager.GetRoom(roomID)
	if err != nil {
		response.Error(c, 404, "Room not found", nil)
		return
	}

	response.Success(c, gin.H{
		"room": room,
	}, "Room information retrieved successfully")
}

// CreateRoom creates a new game room (admin only)
// @Summary Create game room
// @Description Create a new game room (admin only)
// @Tags websocket
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateRoomRequest true "Room details"
// @Success 200 {object} map[string]interface{} "Room created successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Admin access required"
// @Router /ws/rooms [post]
func CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request", err)
		return
	}

	userID, _ := auth.GetUserID(c)

	manager := GetRoomManager()
	room, err := manager.CreateRoom(req.Name, userID, req.MaxPlayers)
	if err != nil {
		response.Error(c, 500, "Failed to create room", err)
		return
	}

	// Broadcast room creation to lobby
	manager.BroadcastToRoom(LobbyRoomID, &Message{
		Type:   MessageTypeRoomCreated,
		RoomID: LobbyRoomID,
		Data: map[string]interface{}{
			"room": room,
		},
	})

	response.Success(c, gin.H{
		"room": room,
	}, "Room created successfully")
}

// CloseRoom closes a game room (admin only)
// @Summary Close game room
// @Description Close an existing game room (admin only)
// @Tags websocket
// @Security BearerAuth
// @Param room_id path string true "Room ID"
// @Success 200 {object} map[string]string "Room closed successfully"
// @Failure 400 {object} map[string]string "Cannot close lobby"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Admin access required"
// @Failure 404 {object} map[string]string "Room not found"
// @Router /ws/rooms/{room_id} [delete]
func CloseRoom(c *gin.Context) {
	roomID := c.Param("room_id")

	manager := GetRoomManager()
	err := manager.CloseRoom(roomID)
	if err != nil {
		if err.Error() == "cannot close lobby room" {
			response.Error(c, 400, err.Error(), nil)
		} else if err.Error() == "room not found" {
			response.Error(c, 404, err.Error(), nil)
		} else {
			response.Error(c, 500, "Failed to close room", err)
		}
		return
	}

	// Broadcast room closure to lobby
	manager.BroadcastToRoom(LobbyRoomID, &Message{
		Type:   MessageTypeRoomClosed,
		RoomID: LobbyRoomID,
		Data: map[string]interface{}{
			"room_id": roomID,
		},
	})

	response.Success(c, gin.H{
		"room_id": roomID,
	}, "Room closed successfully")
}

// InviteToRoom invites users to a game room (admin only)
// @Summary Invite users to room
// @Description Invite specific users to a game room (admin only)
// @Tags websocket
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body InviteRequest true "Invitation details"
// @Success 200 {object} map[string]string "Invitations sent successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Admin access required"
// @Failure 404 {object} map[string]string "Room not found"
// @Router /ws/invite [post]
func InviteToRoom(c *gin.Context) {
	var req InviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "Invalid request", err)
		return
	}

	manager := GetRoomManager()

	// Check if room exists
	room, err := manager.GetRoom(req.RoomID)
	if err != nil {
		response.Error(c, 404, "Room not found", err)
		return
	}

	// Grant permission to join (if room auth is enabled)
	if err := manager.InviteToRoom(req.RoomID, req.UserIDs); err != nil {
		response.Error(c, 500, "Failed to invite users", err)
		return
	}

	// Send invitation to each user
	for _, targetUserID := range req.UserIDs {
		manager.SendToClient(targetUserID, &Message{
			Type:   MessageTypeInvite,
			RoomID: req.RoomID,
			Data: map[string]interface{}{
				"room":    room,
				"message": "You have been invited to join a game room",
			},
		})
	}

	config.Logger.Info("Room invitations sent",
		logger.String("room_id", req.RoomID),
		logger.Int("user_count", len(req.UserIDs)))

	response.Success(c, gin.H{
		"invited_count": len(req.UserIDs),
	}, "Invitations sent successfully")
}
