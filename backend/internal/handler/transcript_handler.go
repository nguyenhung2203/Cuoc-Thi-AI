package handler

import "net/http"

type TranscriptHandler struct {
	// Handler fields
}

func NewTranscriptHandler() *TranscriptHandler {
	return &TranscriptHandler{}
}

func (h *TranscriptHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("transcript handler placeholder"))
}
