package handler

import "net/http"

type MockHandler struct {
	// Handler fields
}

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (h *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mock handler placeholder"))
}
