package api

import (
	"net/http"

	"chat-rooms-backend/state"
)

func GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	WriteJsonHTTP(state.Rooms, w)
}
