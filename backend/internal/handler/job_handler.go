package handler

import "net/http"

type JobHandler struct {
	// Handler fields
}

func NewJobHandler() *JobHandler {
	return &JobHandler{}
}

func (h *JobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("job handler placeholder"))
}
