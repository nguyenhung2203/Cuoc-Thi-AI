package handler

import "net/http"

type AuthHandler struct {
	// Handler fields
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("auth handler placeholder"))
}
