package handler

import "net/http"

type FileHandler struct {
	// Handler fields
}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (h *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("file handler placeholder"))
}
