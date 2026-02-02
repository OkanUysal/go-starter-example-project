package websocket

import "time"

// RoomType defines the type of room
type RoomType string

const (
	// RoomTypeLobby is the public lobby room (always open)
	RoomTypeLobby RoomType = "lobby"

	// RoomTypeGame is a dynamic game room (created/closed by admin)
	RoomTypeGame RoomType = "game"
)

// Room authorization constants
const (
	// LobbyRoomID is the fixed ID for the public lobby
	LobbyRoomID = "lobby"
)

// MessageType defines the type of message
type MessageType string

const (
	// MessageTypeJoin when user joins a room
	MessageTypeJoin MessageType = "join"

	// MessageTypeLeave when user leaves a room
	MessageTypeLeave MessageType = "leave"

	// MessageTypeChat for chat messages
	MessageTypeChat MessageType = "chat"

	// MessageTypeGameEvent for game-specific events
	MessageTypeGameEvent MessageType = "game_event"

	// MessageTypeRoomCreated when admin creates a room
	MessageTypeRoomCreated MessageType = "room_created"

	// MessageTypeRoomClosed when admin closes a room
	MessageTypeRoomClosed MessageType = "room_closed"

	// MessageTypeInvite when user is invited to a room
	MessageTypeInvite MessageType = "invite"

	// MessageTypeError for error messages
	MessageTypeError MessageType = "error"
)

// Message represents a WebSocket message
type Message struct {
	Type      MessageType            `json:"type"`
	RoomID    string                 `json:"room_id"`
	UserID    string                 `json:"user_id"`
	Username  string                 `json:"username,omitempty"`
	Content   string                 `json:"content,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// RoomInfo represents room information
type RoomInfo struct {
	ID           string               `json:"id"`
	Type         RoomType             `json:"type"`
	Name         string               `json:"name"`
	CreatedBy    string               `json:"created_by,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	PlayerCount  int                  `json:"player_count"`
	MaxPlayers   int                  `json:"max_players,omitempty"`
	IsActive     bool                 `json:"is_active"`
	Users        map[string]*UserInfo `json:"users,omitempty"`   // Track users in room
	AllowedUsers map[string]bool      `json:"allowed_users,omitempty"` // Authorized users (if room auth enabled)
}

// UserInfo represents user information in a room
type UserInfo struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	JoinedAt time.Time `json:"joined_at"`
}

// CreateRoomRequest represents a request to create a room
type CreateRoomRequest struct {
	Name       string `json:"name" binding:"required"`
	MaxPlayers int    `json:"max_players"`
}

// InviteRequest represents a request to invite users to a room
type InviteRequest struct {
	RoomID  string   `json:"room_id" binding:"required"`
	UserIDs []string `json:"user_ids" binding:"required"`
}
