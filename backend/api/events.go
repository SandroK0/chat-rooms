package api

import (
	"encoding/json"
	"fmt"
)

type ClientEventType string

const (
	CreateRoom    ClientEventType = "create_room"
	JoinRoom      ClientEventType = "join_room"
	ReconnectRoom ClientEventType = "reconnect_room"
	LeaveRoom     ClientEventType = "leave_room"
	SendMessage   ClientEventType = "send_message"
)

type ServerEventType string

const (
	RoomCreated     ServerEventType = "room_created"
	RoomJoined      ServerEventType = "room_joined"
	RoomLeft        ServerEventType = "room_left"
	RoomReconnected ServerEventType = "room_reconnected"
	InvalidToken    ServerEventType = "invalid_token"
	MessageReceived ServerEventType = "message_received"
	Error           ServerEventType = "error"
)

type ClientEvent struct {
	EventType ClientEventType `json:"eventType"`
	Data      json.RawMessage `json:"data"`
}

func NewClientEvent(msg []byte) (*ClientEvent, error) {
	var event ClientEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

type JoinRoomEventData struct {
	RoomName string `json:"roomName"`
	Username string `json:"username"`
}

type ReconnectRoomEventData struct {
	Token string `json:"token"`
}

type LeaveRoomEventData struct {
	RoomName string `json:"roomName"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type CreateRoomEventData struct {
	RoomName string `json:"roomName"`
	Username string `json:"username"`
}

type SendMessageEventData struct {
	RoomName string `json:"roomName"`
	Username string `json:"username"`
	Body     string `json:"body"`
}

type ServerEvent struct {
	EventType ServerEventType `json:"eventType"`
	Data      any             `json:"data"`
}

func NewServerEvent(eventType ServerEventType, data any) ServerEvent {
	return ServerEvent{
		EventType: eventType,
		Data:      data,
	}
}

type TokenData struct {
	Token string `json:"token"`
}

type ErrorEventData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type RoomJoinedEventData struct {
	Token    string `json:"token"`
	RoomName string `json:"roomName"`
}

type RoomLeftEventData struct {
	Token    string `json:"token"`
	RoomName string `json:"roomName"`
}

type MessageReceivedEventData struct {
	Username string `json:"username"`
	Body     string `json:"body"`
}

type RoomCreatedEventData struct {
	Token    string `json:"token"`
	RoomName string `json:"roomName"`
}

type RoomReconnectedEventData struct {
	Token    string `json:"token"`
	RoomName string `json:"roomName"`
	Username string `json:"username"`
}

type ClientEventData interface{}

type ClientEventDataFactory func() ClientEventData

var clientEventFactories = map[ClientEventType]ClientEventDataFactory{
	CreateRoom:    func() ClientEventData { return &CreateRoomEventData{} },
	JoinRoom:      func() ClientEventData { return &JoinRoomEventData{} },
	LeaveRoom:     func() ClientEventData { return &LeaveRoomEventData{} },
	ReconnectRoom: func() ClientEventData { return &ReconnectRoomEventData{} },
	SendMessage:   func() ClientEventData { return &SendMessageEventData{} },
}

func UnmarshalClientEventData(event ClientEvent) (ClientEventData, error) {
	factory, exists := clientEventFactories[event.EventType]
	if !exists {
		return nil, fmt.Errorf("unknown event type: %s", event.EventType)
	}
	eventData := factory()
	if err := json.Unmarshal(event.Data, eventData); err != nil {
		return nil, err
	}

	return eventData, nil
}
