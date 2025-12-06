package api

import (
	"log"
	"net/http"

	"chat-rooms-backend/entities"
	"chat-rooms-backend/state"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()
	state.Clients[ws] = true
	log.Println("Client connected:", len(state.Clients))

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			delete(state.Clients, ws)
			break
		}

		event, err := NewClientEvent(msg)
		if err != nil {
			log.Println(err)
			break
		}

		HandleEvents(*event, ws)
	}

}

func HandleMessages() {
	for {
		msgInterface := <-state.Messages
		msg, ok := msgInterface.(*entities.Message)
		if !ok {
			log.Printf("Invalid message type received: %T", msgInterface)
			continue
		}

		switch msg.Type {
		case entities.GlobalBroadcast:
			for client := range state.Clients {
				err := client.WriteJSON(msg.Content)
				if err != nil {
					log.Println("Write error:", err)
					client.Close()
					delete(state.Clients, client)
				}
			}
		case entities.ClientSpecific:
			err := msg.Client.WriteJSON(msg.Content)
			if err != nil {
				log.Println("Write error:", err)
				msg.Client.Close()
				delete(state.Clients, msg.Client)
			}
		case entities.RoomBroadcast:
			room, exists := state.Rooms[msg.RoomName]
			if !exists {
				log.Println("Room not found:", msg.RoomName)
				continue
			}

			for _, user := range room.Users {
				err := user.Conn.WriteJSON(msg.Content)
				if err != nil {
					log.Println("Write error:", err)
					user.Conn.Close()
					delete(state.Clients, user.Conn)
				}
			}
		}

	}
}
