package api

import (
	"github.com/gorilla/websocket"
)

// Event creation helpers for common scenarios

// SendRoomJoined sends a room joined event to a specific client
func SendRoomJoined(token, roomName string, ws *websocket.Conn) error {
	data := RoomJoinedEventData{Token: token, RoomName: roomName}
	return sendEvent(RoomJoined, data, ws)
}

// SendRoomLeft sends a room left event to a specific client
func SendRoomLeft(token, roomName string, ws *websocket.Conn) error {
	data := RoomLeftEventData{Token: token, RoomName: roomName}
	return sendEvent(RoomLeft, data, ws)
}

// SendRoomCreated sends a room created event to a specific client
func SendRoomCreated(token, roomName string, ws *websocket.Conn) error {
	data := RoomCreatedEventData{Token: token, RoomName: roomName}
	return sendEvent(RoomCreated, data, ws)
}

// SendRoomReconnected sends a room reconnected event to a specific client
func SendRoomReconnected(token, roomName, username string, ws *websocket.Conn) error {
	data := RoomReconnectedEventData{Token: token, RoomName: roomName, Username: username}
	return sendEvent(RoomReconnected, data, ws)
}

// SendInvalidToken sends an invalid token event to a specific client
func SendInvalidToken(token string, ws *websocket.Conn) error {
	data := TokenData{Token: token}
	return sendEvent(InvalidToken, data, ws)
}

// BroadcastMessage broadcasts a message to all users in a room
func BroadcastMessage(username, body, roomName string) error {
	data := MessageReceivedEventData{Username: username, Body: body}
	return sendEventToRoom(MessageReceived, data, roomName)
}

// SendRoomCreationEvents sends both room created and room joined events
func SendRoomCreationEvents(token, roomName string, ws *websocket.Conn) error {
	events := []ServerEvent{
		NewServerEvent(RoomCreated, RoomCreatedEventData{Token: token, RoomName: roomName}),
		NewServerEvent(RoomJoined, RoomJoinedEventData{Token: token, RoomName: roomName}),
	}
	return sendMultipleEvents(events, ws)
}
