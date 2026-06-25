package handler

import "net/http"

type CompanyHandler struct {
	// Handler fields
}

func NewCompanyHandler() *CompanyHandler {
	return &CompanyHandler{}
}

func (h *CompanyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("company handler placeholder"))
}
