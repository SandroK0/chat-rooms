package state

import (
	"chat-rooms-backend/entities"

	"github.com/gorilla/websocket"
)

type UserRoom struct {
	Username string
	RoomName string
}

func NewUserRoom(username, roomname string) *UserRoom {
	return &UserRoom{Username: username, RoomName: roomname}
}

// Global state variables
var (
	// Clients holds all active websocket connections
	Clients = make(map[*websocket.Conn]bool)

	// Messages is the channel for handling chat messages - using interface{} to accept api.Message
	Messages = make(chan interface{})

	// Rooms maps room names to Room entities
	Rooms = make(map[string]*entities.Room)

	// TokenToRooms maps user tokens to their current room information
	TokenToRooms = make(map[string]*UserRoom)
)
