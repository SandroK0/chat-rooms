package api

import (
	"fmt"

	"github.com/SandroK0/chat-rooms/backend/entities"
	"github.com/gorilla/websocket"
)

// Helper functions for DRY event handling

// sendEvent creates and sends a server event to a specific client
func sendEvent(eventType ServerEventType, data any, ws *websocket.Conn) error {
	event := NewServerEvent(eventType, data)
	msg, err := NewClientMessage(ws, event)
	if err != nil {
		return fmt.Errorf("creating client message: %w", err)
	}
	Messages <- msg
	return nil
}

// sendEventToRoom creates and sends a server event to all clients in a room
func sendEventToRoom(eventType ServerEventType, data any, roomName string) error {
	event := NewServerEvent(eventType, data)
	msg, err := NewRoomMessage(roomName, event)
	if err != nil {
		return fmt.Errorf("creating room message: %w", err)
	}
	Messages <- msg
	return nil
}

// validateEventData checks if required fields are not empty
func validateEventData(fields map[string]string) error {
	for fieldName, value := range fields {
		if value == "" {
			return fmt.Errorf("missing %s", fieldName)
		}
	}
	return nil
}

// assertEventData performs type assertion and returns typed data
func assertEventData[T any](clientEventData any, eventType ClientEventType) (T, error) {
	var zero T
	data, ok := clientEventData.(T)
	if !ok {
		return zero, fmt.Errorf("invalid event data type for %s", eventType)
	}
	return data, nil
}

func getRoom(roomName string) (*entities.Room, error) {
	room, exists := Rooms[roomName]
	if !exists {
		return nil, fmt.Errorf("room not found: %s", roomName)
	}
	return room, nil
}

func SendError(Code, Message string, ws *websocket.Conn) {
	errorData := ErrorEventData{Code: Code, Message: Message}
	if err := sendEvent(Error, errorData, ws); err != nil {
		HandleEventError(err, "sending error event")
	}
}

func createRoom(clientEventData ClientEventData, ws *websocket.Conn) {
	data, err := assertEventData[*CreateRoomEventData](clientEventData, CreateRoom)
	if err != nil {
		HandleEventError(err, "type assertion")
		return
	}

	if err := validateEventData(map[string]string{
		"roomName": data.RoomName,
		"username": data.Username,
	}); err != nil {
		HandleEventError(err, "create_room event")
		return
	}

	if _, exists := Rooms[data.RoomName]; exists {
		SendError("RoomAlreadyExists", "Room name is taken", ws)
		return
	}

	room := entities.NewRoom(data.RoomName)
	Rooms[data.RoomName] = room

	user := entities.NewUser(data.Username, ws)
	room.AddUser(user)

	TokenToRooms[user.Token] = NewUserRoom(user.Name, room.Name)

	if err := sendEvent(RoomCreated, RoomCreatedEventData{Token: user.Token, RoomName: room.Name}, ws); err != nil {
		HandleEventError(err, "sending room created event")
		return
	}

	if err := sendEvent(RoomJoined, RoomJoinedEventData{Token: user.Token, RoomName: room.Name}, ws); err != nil {
		HandleEventError(err, "sending room joined event")
		return
	}
}

func joinRoom(clientEventData ClientEventData, ws *websocket.Conn) {
	data, err := assertEventData[*JoinRoomEventData](clientEventData, JoinRoom)
	if err != nil {
		HandleEventError(err, "type assertion")
		return
	}

	if err := validateEventData(map[string]string{
		"roomName": data.RoomName,
		"username": data.Username,
	}); err != nil {
		HandleEventError(err, "join event")
		return
	}

	room, err := getRoom(data.RoomName)
	if err != nil {
		HandleEventError(err, "join event")
		return
	}

	user := entities.NewUser(data.Username, ws)

	err = room.AddUser(user)
	if err != nil {
		SendError("UsernameTaken", "User with that name already exists in that room", ws)
		return
	}

	TokenToRooms[user.Token] = NewUserRoom(user.Name, room.Name)

	joinedData := RoomJoinedEventData{Token: user.Token, RoomName: room.Name}
	if err := sendEvent(RoomJoined, joinedData, ws); err != nil {
		HandleEventError(err, "sending room joined event")
		return
	}
}

func leaveRoom(clientEventData ClientEventData, ws *websocket.Conn) {
	data, err := assertEventData[*LeaveRoomEventData](clientEventData, LeaveRoom)
	if err != nil {
		HandleEventError(err, "type assertion")
		return
	}

	if err := validateEventData(map[string]string{
		"roomName": data.RoomName,
		"username": data.Username,
		"token":    data.Token,
	}); err != nil {
		HandleEventError(err, "leave event")
		return
	}

	room, err := getRoom(data.RoomName)
	if err != nil {
		HandleEventError(err, "leave event")
		return
	}

	room.RemoveUser(data.Token)

	leftData := RoomLeftEventData{Token: data.Token, RoomName: data.RoomName}
	if err := sendEvent(RoomLeft, leftData, ws); err != nil {
		HandleEventError(err, "sending room left event")
		return
	}
}

func reconnectRoom(clientEventData ClientEventData, ws *websocket.Conn) {
	data, err := assertEventData[*ReconnectRoomEventData](clientEventData, ReconnectRoom)
	if err != nil {
		HandleEventError(err, "type assertion")
		return
	}

	userRoom, ok := TokenToRooms[data.Token]
	if !ok {
		tokenData := TokenData{Token: data.Token}
		if err := sendEvent(InvalidToken, tokenData, ws); err != nil {
			HandleEventError(err, "sending invalid token event")
		}
		return
	}

	room, err := getRoom(userRoom.RoomName)
	if err != nil {
		HandleEventError(err, "reconnect_room event")
		return
	}

	user := room.GetUserByToken(data.Token)
	if user == nil {
		tokenData := TokenData{Token: data.Token}
		if err := sendEvent(InvalidToken, tokenData, ws); err != nil {
			HandleEventError(err, "sending invalid token event")
		}
		return
	}

	user.Conn = ws
	reconnectData := RoomReconnectedEventData{Token: user.Token, RoomName: room.Name, Username: user.Name}
	if err := sendEvent(RoomReconnected, reconnectData, ws); err != nil {
		HandleEventError(err, "sending room reconnected event")
		return
	}
}

func sendMessage(clientEventData ClientEventData) {
	data, err := assertEventData[*SendMessageEventData](clientEventData, SendMessage)
	if err != nil {
		HandleEventError(err, "type assertion")
		return
	}

	if err := validateEventData(map[string]string{
		"roomName": data.RoomName,
		"username": data.Username,
		"body":     data.Body,
	}); err != nil {
		HandleEventError(err, "message event")
		return
	}

	room, err := getRoom(data.RoomName)
	if err != nil {
		HandleEventError(err, "message event")
		return
	}

	messageData := MessageReceivedEventData{Username: data.Username, Body: data.Body}
	if err := sendEventToRoom(MessageReceived, messageData, room.Name); err != nil {
		HandleEventError(err, "broadcasting message")
		return
	}
}

func HandleEvents(event ClientEvent, ws *websocket.Conn) {

	clientEventData, err := UnmarshalClientEventData(event)
	if err != nil {
		HandleEventError(err, "unmarshaling create_room event")
		return
	}

	switch event.EventType {
	case CreateRoom:
		createRoom(clientEventData, ws)
	case JoinRoom:
		joinRoom(clientEventData, ws)
	case LeaveRoom:
		leaveRoom(clientEventData, ws)
	case ReconnectRoom:
		reconnectRoom(clientEventData, ws)
	case SendMessage:
		sendMessage(clientEventData)
	default:
		HandleEventError(fmt.Errorf("unknown event type: %s", event.EventType), "handling event")

	}
}
