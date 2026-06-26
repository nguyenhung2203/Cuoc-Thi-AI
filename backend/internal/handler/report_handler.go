package handler

import "net/http"

type ReportHandler struct {
	// Handler fields
}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{}
}

func (h *ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("report handler placeholder"))
}
