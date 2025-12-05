package main

import (
	"log"
	"net/http"

	"github.com/SandroK0/chat-rooms/backend/api"
	"github.com/SandroK0/chat-rooms/backend/pkg/logger"
	"github.com/rs/cors"
)

func main() {
	// Setup file logging
	if err := logger.SetupFileLogging(); err != nil {
		log.Fatalf("Failed to setup file logging: %v", err)
	}

	log.Println("Starting chat rooms backend server...")

	mux := http.NewServeMux()

	mux.HandleFunc("/rooms", api.GetRoomsHandler)

	mux.HandleFunc("/ws", api.HandleConnections)
	go api.HandleMessages()

	handler := cors.Default().Handler(mux)
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
