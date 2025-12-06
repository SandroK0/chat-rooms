package api

import (
	"chat-rooms-backend/entities"
	"chat-rooms-backend/state"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func SendEvent(eventType ServerEventType, data any, ws *websocket.Conn) error {
	event := NewServerEvent(eventType, data)
	msg, err := entities.NewClientMessage(ws, event)
	if err != nil {
		return fmt.Errorf("creating client message: %w", err)
	}
	state.Messages <- msg
	return nil
}

func SendEventToRoom(eventType ServerEventType, data any, roomName string) error {
	event := NewServerEvent(eventType, data)
	msg, err := entities.NewRoomMessage(roomName, event)
	if err != nil {
		return fmt.Errorf("creating room message: %w", err)
	}
	state.Messages <- msg
	return nil
}

func ValidateEventData(fields map[string]string) error {
	for fieldName, value := range fields {
		if value == "" {
			return fmt.Errorf("missing %s", fieldName)
		}
	}
	return nil
}

func AssertEventData[T any](clientEventData any, eventType ClientEventType) (T, error) {
	var zero T
	data, ok := clientEventData.(T)
	if !ok {
		return zero, fmt.Errorf("invalid event data type for %s", eventType)
	}
	return data, nil
}

func GetRoom(roomName string) (*entities.Room, error) {
	room, exists := state.Rooms[roomName]
	if !exists {
		return nil, fmt.Errorf("room not found: %s", roomName)
	}
	return room, nil
}

func SendError(Code, Message string, ws *websocket.Conn) {
	errorData := ErrorEventData{Code: Code, Message: Message}
	if err := SendEvent(Error, errorData, ws); err != nil {
		HandleEventError(err, "sending error event")
	}
}

func WriteJsonHTTP(data any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func HandleEventError(err error, context string) {
	log.Printf("Error in %s: %v", context, err)
}
