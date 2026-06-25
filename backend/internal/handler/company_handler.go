package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"backend/internal/middleware"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/response"
	"backend/internal/service"
)

type CompanyHandler struct {
	companyService *service.CompanyService
}

func NewCompanyHandler(companyService *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: companyService}
}

func (h *CompanyHandler) Routes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{company_id}", h.Get)
	r.Put("/{company_id}", h.Update)
}

type CreateCompanyReq struct {
	Name     string `json:"name"`
	Website  string `json:"website"`
	Industry string `json:"industry"`
	Size     string `json:"size"`
}

func (h *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok {
		response.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}

	companies, err := h.companyService.ListCompanies(r.Context(), userID)
	if err != nil {
		response.Error(w, errors.NewInternal("failed to list companies"), "")
		return
	}

	response.JSON(w, http.StatusOK, companies, nil, "")
}

func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok {
		response.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}

	var req CreateCompanyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}

	company, err := h.companyService.CreateCompany(r.Context(), userID, req.Name, req.Website, req.Industry, req.Size)
	if err != nil {
		response.Error(w, errors.NewInternal("failed to create company"), "")
		return
	}

	response.JSON(w, http.StatusCreated, company, nil, "")
}

func (h *CompanyHandler) Get(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	company, err := h.companyService.GetCompany(r.Context(), companyID)
	if err != nil {
		response.Error(w, errors.NewNotFound("company not found"), "")
		return
	}
	response.JSON(w, http.StatusOK, company, nil, "")
}

func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	
	var req CreateCompanyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}

	company, err := h.companyService.UpdateCompany(r.Context(), companyID, req.Name, req.Website, req.Industry, req.Size)
	if err != nil {
		response.Error(w, errors.NewInternal("failed to update company"), "")
		return
	}

	response.JSON(w, http.StatusOK, company, nil, "")
}
