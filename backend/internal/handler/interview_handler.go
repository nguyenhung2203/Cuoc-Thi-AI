package handler

import "net/http"

type InterviewHandler struct {
	// Handler fields
}

func NewInterviewHandler() *InterviewHandler {
	return &InterviewHandler{}
}

func (h *InterviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("interview handler placeholder"))
}
