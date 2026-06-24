package handler

import "net/http"

type RoomHandler struct {
	// Handler fields
}

func NewRoomHandler() *RoomHandler {
	return &RoomHandler{}
}

func (h *RoomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("room handler placeholder"))
}
