package api

import (
	"fmt"

	"chat-rooms-backend/entities"
	"chat-rooms-backend/state"

	"github.com/gorilla/websocket"
)

func handleCreateRoom(clientEventData ClientEventData, ws *websocket.Conn) {
	data, err := AssertEventData[*CreateRoomEventData](clientEventData, CreateRoom)
	if err != nil {
		HandleEventError(err, "type assertion")
		return
	}

	if err := ValidateEventData(map[string]string{
		"roomName": data.RoomName,
		"username": data.Username,
	}); err != nil {
		HandleEventError(err, "create_room event")
		return
	}

	if _, exists := state.Rooms[data.RoomName]; exists {
		SendError("RoomAlreadyExists", "Room name is taken", ws)
		return
	}

	room := entities.NewRoom(data.RoomName)
	state.Rooms[data.RoomName] = room

	user := entities.NewUser(data.Username, ws)
	room.AddUser(user)

	state.TokenToRooms[user.Token] = state.NewUserRoom(user.Name, room.Name)

	if err := SendEvent(RoomCreated, RoomCreatedEventData{Token: user.Token, RoomName: room.Name}, ws); err != nil {
		HandleEventError(err, "sending room created event")
		return
	}

	if err := SendEvent(RoomJoined, RoomJoinedEventData{Token: user.Token, RoomName: room.Name}, ws); err != nil {
		HandleEventError(err, "sending room joined event")
		return
	}
}

func handleJoinRoom(clientEventData ClientEventData, ws *websocket.Conn) {
	data, err := AssertEventData[*JoinRoomEventData](clientEventData, JoinRoom)
	if err != nil {
		HandleEventError(err, "type assertion")
		return
	}

	if err := ValidateEventData(map[string]string{
		"roomName": data.RoomName,
		"username": data.Username,
	}); err != nil {
		HandleEventError(err, "join event")
		return
	}

	room, err := GetRoom(data.RoomName)
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

	state.TokenToRooms[user.Token] = state.NewUserRoom(user.Name, room.Name)

	joinedData := RoomJoinedEventData{Token: user.Token, RoomName: room.Name}
	if err := SendEvent(RoomJoined, joinedData, ws); err != nil {
		HandleEventError(err, "sending room joined event")
		return
	}
}

func handleLeaveRoom(clientEventData ClientEventData, ws *websocket.Conn) {
	data, err := AssertEventData[*LeaveRoomEventData](clientEventData, LeaveRoom)
	if err != nil {
		HandleEventError(err, "type assertion")
		return
	}

	if err := ValidateEventData(map[string]string{
		"roomName": data.RoomName,
		"username": data.Username,
		"token":    data.Token,
	}); err != nil {
		HandleEventError(err, "leave event")
		return
	}

	room, err := GetRoom(data.RoomName)
	if err != nil {
		HandleEventError(err, "leave event")
		return
	}

	room.RemoveUser(data.Token)

	leftData := RoomLeftEventData{Token: data.Token, RoomName: data.RoomName}
	if err := SendEvent(RoomLeft, leftData, ws); err != nil {
		HandleEventError(err, "sending room left event")
		return
	}
}

func handleReconnectRoom(clientEventData ClientEventData, ws *websocket.Conn) {
	data, err := AssertEventData[*ReconnectRoomEventData](clientEventData, ReconnectRoom)
	if err != nil {
		HandleEventError(err, "type assertion")
		return
	}

	userRoom, ok := state.TokenToRooms[data.Token]
	if !ok {
		tokenData := TokenData{Token: data.Token}
		if err := SendEvent(InvalidToken, tokenData, ws); err != nil {
			HandleEventError(err, "sending invalid token event")
		}
		return
	}

	room, err := GetRoom(userRoom.RoomName)
	if err != nil {
		HandleEventError(err, "reconnect_room event")
		return
	}

	user := room.GetUserByToken(data.Token)
	if user == nil {
		tokenData := TokenData{Token: data.Token}
		if err := SendEvent(InvalidToken, tokenData, ws); err != nil {
			HandleEventError(err, "sending invalid token event")
		}
		return
	}

	user.Conn = ws
	reconnectData := RoomReconnectedEventData{Token: user.Token, RoomName: room.Name, Username: user.Name}
	if err := SendEvent(RoomReconnected, reconnectData, ws); err != nil {
		HandleEventError(err, "sending room reconnected event")
		return
	}
}

func handleSendChatMessage(clientEventData ClientEventData) {
	data, err := AssertEventData[*SendChatMessageEventData](clientEventData, SendChatMessage)
	if err != nil {
		HandleEventError(err, "type assertion")
		return
	}

	if err := ValidateEventData(map[string]string{
		"roomName": data.RoomName,
		"username": data.Username,
		"body":     data.Body,
	}); err != nil {
		HandleEventError(err, "chat message event")
		return
	}

	room, err := GetRoom(data.RoomName)
	if err != nil {
		HandleEventError(err, "chat message event")
		return
	}

	messageData := ChatMessageReceivedEventData{Username: data.Username, Body: data.Body}
	if err := SendEventToRoom(ChatMessageReceived, messageData, room.Name); err != nil {
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
		handleCreateRoom(clientEventData, ws)
	case JoinRoom:
		handleJoinRoom(clientEventData, ws)
	case LeaveRoom:
		handleLeaveRoom(clientEventData, ws)
	case ReconnectRoom:
		handleReconnectRoom(clientEventData, ws)
	case SendChatMessage:
		handleSendChatMessage(clientEventData)
	default:
		HandleEventError(fmt.Errorf("unknown event type: %s", event.EventType), "handling event")

	}
}
