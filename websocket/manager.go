package websocket

import (
	"fmt"
	"sync"
	"time"

	"github.com/OkanUysal/go-logger"
	"github.com/OkanUysal/go-starter-example-project/config"
	gowebsocket "github.com/OkanUysal/go-websocket"
	"github.com/google/uuid"
)

const (
	// LobbyRoomID is the ID for the public lobby
	LobbyRoomID = "lobby"
)

// RoomManager manages all WebSocket rooms
type RoomManager struct {
	hub   *gowebsocket.Hub
	rooms map[string]*RoomInfo
	mu    sync.RWMutex
}

var (
	manager     *RoomManager
	managerOnce sync.Once
)

// GetRoomManager returns the singleton room manager instance
func GetRoomManager() *RoomManager {
	managerOnce.Do(func() {
		manager = &RoomManager{
			hub:   gowebsocket.NewHub(nil), // Use default config
			rooms: make(map[string]*RoomInfo),
		}

		// Set up message handler
		manager.hub.SetOnMessage(manager.handleMessage)

		// Create lobby room (always open)
		manager.rooms[LobbyRoomID] = &RoomInfo{
			ID:        LobbyRoomID,
			Type:      RoomTypeLobby,
			Name:      "Public Lobby",
			CreatedAt: time.Now(),
			IsActive:  true,
			Users:     make(map[string]*UserInfo),
		}

		// Also create lobby in hub (for BroadcastToRoom to work)
		// We need to manually create the room in hub's internal map
		// Since go-websocket doesn't expose a way to create room with custom ID,
		// we'll create it when first user joins

		config.Logger.Info("WebSocket room manager initialized",
			logger.String("lobby_id", LobbyRoomID))
	})
	return manager
}

// Start starts the room manager hub
func (rm *RoomManager) Start() {
	go rm.hub.Run()

	// Give hub time to initialize
	time.Sleep(100 * time.Millisecond)

	// Create lobby room in hub with fixed ID
	err := rm.hub.CreateRoomWithID(LobbyRoomID, &gowebsocket.RoomConfig{
		Name:       "Public Lobby",
		MaxClients: 0, // unlimited
		IsPrivate:  false,
	})

	if err != nil {
		config.Logger.Error("Failed to create lobby room", logger.Err(err))
	} else {
		config.Logger.Info("Lobby room created in hub", logger.String("room_id", LobbyRoomID))
	}
}

// GetHub returns the WebSocket hub
func (rm *RoomManager) GetHub() *gowebsocket.Hub {
	return rm.hub
}

// CreateRoom creates a new game room (admin only)
func (rm *RoomManager) CreateRoom(name, createdBy string, maxPlayers int) (*RoomInfo, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	roomID := uuid.New().String()

	room := &RoomInfo{
		ID:         roomID,
		Type:       RoomTypeGame,
		Name:       name,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
		MaxPlayers: maxPlayers,
		IsActive:   true,
		Users:      make(map[string]*UserInfo),
	}

	rm.rooms[roomID] = room

	// Create room in hub as well
	err := rm.hub.CreateRoomWithID(roomID, &gowebsocket.RoomConfig{
		Name:       name,
		MaxClients: maxPlayers,
		IsPrivate:  false,
	})

	if err != nil {
		config.Logger.Error("Failed to create room in hub",
			logger.Err(err),
			logger.String("room_id", roomID))
		delete(rm.rooms, roomID)
		return nil, err
	}

	config.Logger.Info("Game room created",
		logger.String("room_id", roomID),
		logger.String("name", name),
		logger.String("created_by", createdBy),
		logger.Int("max_players", maxPlayers))

	return room, nil
}

// CloseRoom closes a game room (admin only)
func (rm *RoomManager) CloseRoom(roomID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if roomID == LobbyRoomID {
		return fmt.Errorf("cannot close lobby room")
	}

	room, exists := rm.rooms[roomID]
	if !exists {
		return fmt.Errorf("room not found")
	}

	if !room.IsActive {
		return fmt.Errorf("room already closed")
	}

	// Mark as inactive
	room.IsActive = false

	// Notify all users in the room
	rm.hub.BroadcastToRoom(roomID, gowebsocket.Message{
		Type: string(MessageTypeRoomClosed),
		Data: map[string]interface{}{
			"room_id": roomID,
			"message": "This room has been closed by admin",
		},
	})

	// Remove all clients from the room
	rm.hub.CloseRoom(roomID)

	config.Logger.Info("Game room closed",
		logger.String("room_id", roomID),
		logger.String("name", room.Name))

	return nil
}

// GetRoom returns room information
func (rm *RoomManager) GetRoom(roomID string) (*RoomInfo, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("room not found")
	}

	// Update player count
	room.PlayerCount = rm.hub.GetRoomClientCount(roomID)

	return room, nil
}

// GetAllRooms returns all active rooms
func (rm *RoomManager) GetAllRooms() []*RoomInfo {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	rooms := make([]*RoomInfo, 0, len(rm.rooms))
	for _, room := range rm.rooms {
		if room.IsActive {
			// Update player count
			room.PlayerCount = rm.hub.GetRoomClientCount(room.ID)
			rooms = append(rooms, room)
		}
	}

	return rooms
}

// JoinRoom adds a client to a room
func (rm *RoomManager) JoinRoom(roomID, userID, username string) error {
	rm.mu.RLock()
	room, exists := rm.rooms[roomID]
	rm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("room not found")
	}

	if !room.IsActive {
		return fmt.Errorf("room is not active")
	}

	// Check max players for game rooms
	if room.Type == RoomTypeGame && room.MaxPlayers > 0 {
		currentCount := rm.hub.GetRoomClientCount(roomID)
		if currentCount >= room.MaxPlayers {
			return fmt.Errorf("room is full")
		}
	}

	// Add user to room
	rm.mu.Lock()
	room.Users[userID] = &UserInfo{
		UserID:   userID,
		Username: username,
		JoinedAt: time.Now(),
	}
	rm.mu.Unlock()

	// Try to join the user to the room in the hub
	err := rm.hub.JoinRoom(userID, roomID)
	if err != nil && err.Error() == "room not found" {
		// Room doesn't exist in hub, create it with our ID
		createErr := rm.hub.CreateRoomWithID(roomID, &gowebsocket.RoomConfig{
			Name:       room.Name,
			MaxClients: room.MaxPlayers,
			IsPrivate:  false,
		})

		if createErr != nil {
			config.Logger.Error("Failed to create room in hub",
				logger.Err(createErr),
				logger.String("room_id", roomID))
			return createErr
		}

		config.Logger.Info("Created room in hub", logger.String("room_id", roomID))

		// Try joining again
		err = rm.hub.JoinRoom(userID, roomID)
	}

	if err != nil {
		config.Logger.Error("Failed to join hub room",
			logger.Err(err),
			logger.String("user_id", userID),
			logger.String("room_id", roomID))
		return err
	}

	// Broadcast join message to room
	rm.hub.BroadcastToRoom(roomID, gowebsocket.Message{
		Type: string(MessageTypeJoin),
		Data: map[string]interface{}{
			"room_id":  roomID,
			"user_id":  userID,
			"username": username,
			"message":  fmt.Sprintf("%s joined the room", username),
		},
	})

	config.Logger.Info("User joined room",
		logger.String("user_id", userID),
		logger.String("username", username),
		logger.String("room_id", roomID))

	return nil
}

// LeaveRoom removes a client from a room
func (rm *RoomManager) LeaveRoom(roomID, userID, username string) {
	// Remove user from room
	rm.mu.Lock()
	if room, exists := rm.rooms[roomID]; exists {
		delete(room.Users, userID)
	}
	rm.mu.Unlock()

	// Leave the room in the hub
	rm.hub.LeaveRoom(userID, roomID)

	// Broadcast leave message to room
	rm.hub.BroadcastToRoom(roomID, gowebsocket.Message{
		Type: string(MessageTypeLeave),
		Data: map[string]interface{}{
			"room_id":  roomID,
			"user_id":  userID,
			"username": username,
			"message":  fmt.Sprintf("%s left the room", username),
		},
	})

	config.Logger.Info("User left room",
		logger.String("user_id", userID),
		logger.String("username", username),
		logger.String("room_id", roomID))
}

// BroadcastToRoom sends a message to all clients in a room
func (rm *RoomManager) BroadcastToRoom(roomID string, message *Message) {
	message.Timestamp = time.Now()

	// Build data map
	data := map[string]interface{}{
		"room_id":   message.RoomID,
		"user_id":   message.UserID,
		"username":  message.Username,
		"content":   message.Content,
		"timestamp": message.Timestamp,
	}

	// Merge message.Data if exists
	if message.Data != nil {
		for k, v := range message.Data {
			data[k] = v
		}
	}

	rm.hub.BroadcastToRoom(roomID, gowebsocket.Message{
		Type: string(message.Type),
		Data: data,
	})
}

// SendToClient sends a message to a specific client
func (rm *RoomManager) SendToClient(clientID string, message *Message) {
	message.Timestamp = time.Now()

	// Build data map
	data := map[string]interface{}{
		"room_id":   message.RoomID,
		"user_id":   message.UserID,
		"username":  message.Username,
		"content":   message.Content,
		"timestamp": message.Timestamp,
	}

	// Merge message.Data if exists
	if message.Data != nil {
		for k, v := range message.Data {
			data[k] = v
		}
	}

	rm.hub.SendToUser(clientID, gowebsocket.Message{
		Type: string(message.Type),
		Data: data,
	})
}

// handleMessage processes incoming WebSocket messages
func (rm *RoomManager) handleMessage(client *gowebsocket.Client, msg gowebsocket.Message) {
	config.Logger.Info("WebSocket message received",
		logger.String("user_id", client.UserID),
		logger.String("type", msg.Type))

	// Extract message data
	data := msg.Data

	switch msg.Type {
	case "chat":
		// Handle chat message
		roomID, _ := data["room_id"].(string)
		content, _ := data["content"].(string)

		if roomID == "" || content == "" {
			rm.SendToClient(client.UserID, &Message{
				Type: MessageTypeError,
				Data: map[string]interface{}{
					"message": "Invalid chat message format",
				},
			})
			return
		}

		// Get user info from room
		rm.mu.RLock()
		room, exists := rm.rooms[roomID]
		rm.mu.RUnlock()

		if !exists {
			rm.SendToClient(client.UserID, &Message{
				Type: MessageTypeError,
				Data: map[string]interface{}{
					"message": "Room not found",
				},
			})
			return
		}

		// Get username
		username := client.UserID
		for userID, user := range room.Users {
			if userID == client.UserID {
				username = user.Username
				break
			}
		}

		// Broadcast chat message to room
		rm.BroadcastToRoom(roomID, &Message{
			Type:     MessageTypeChat,
			RoomID:   roomID,
			UserID:   client.UserID,
			Username: username,
			Content:  content,
		})

	case "create_room":
		// Handle room creation (admin only)
		roomName, _ := data["name"].(string)
		if roomName == "" {
			roomName = "Game Room"
		}

		room, err := rm.CreateRoom(roomName, client.UserID, 10)
		if err != nil {
			rm.SendToClient(client.UserID, &Message{
				Type: MessageTypeError,
				Data: map[string]interface{}{
					"message": err.Error(),
				},
			})
			return
		}

		// Notify about room creation
		rm.BroadcastToRoom(LobbyRoomID, &Message{
			Type: MessageTypeRoomCreated,
			Data: map[string]interface{}{
				"room_id": room.ID,
				"name":    room.Name,
			},
		})

	case "close_room":
		// Handle room closure (admin only)
		roomID, _ := data["room_id"].(string)

		if err := rm.CloseRoom(roomID); err != nil {
			rm.SendToClient(client.UserID, &Message{
				Type: MessageTypeError,
				Data: map[string]interface{}{
					"message": err.Error(),
				},
			})
		}

	default:
		config.Logger.Warn("Unknown message type",
			logger.String("type", msg.Type),
			logger.String("user_id", client.UserID))
	}
}
