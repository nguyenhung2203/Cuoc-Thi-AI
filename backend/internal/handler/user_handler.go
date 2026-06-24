package handler

import "net/http"

type UserHandler struct {
	// Handler fields
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user handler placeholder"))
}
