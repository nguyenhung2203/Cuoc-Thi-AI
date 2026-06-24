package handler

import "net/http"

type AuditHandler struct {
	// Handler fields
}

func NewAuditHandler() *AuditHandler {
	return &AuditHandler{}
}

func (h *AuditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("audit handler placeholder"))
}
