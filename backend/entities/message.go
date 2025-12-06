package entities

import (
	"errors"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	ClientSpecific  MessageType = "client"
	RoomBroadcast   MessageType = "room"
	GlobalBroadcast MessageType = "global"
)

type Message struct {
	Type     MessageType
	Client   *websocket.Conn // For client-specific messages
	RoomName string          // For room-specific messages
	Content  any
}

func NewClientMessage(client *websocket.Conn, content any) (*Message, error) {
	if client == nil {
		return nil, errors.New("client connection cannot be nil")
	}
	return &Message{
		Type:    ClientSpecific,
		Client:  client,
		Content: content,
	}, nil
}

func NewRoomMessage(roomName string, content any) (*Message, error) {
	if roomName == "" {
		return nil, errors.New("roomName cannot be empty")
	}
	return &Message{
		Type:     RoomBroadcast,
		RoomName: roomName,
		Content:  content,
	}, nil
}

func NewBroadcastMessage(content any) (*Message, error) {
	return &Message{
		Type:    GlobalBroadcast,
		Content: content,
	}, nil
}
