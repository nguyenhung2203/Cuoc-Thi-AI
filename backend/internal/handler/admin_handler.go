package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/pagination"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/repository"
	"backend/internal/service"
)

type AdminHandler struct {
	adminSvc  *service.AdminService
	userRepo  *repository.UserRepository
	compRepo  *repository.CompanyRepository
	auditRepo *repository.AuditRepository
}

func NewAdminHandler(adminSvc *service.AdminService, userRepo *repository.UserRepository, compRepo *repository.CompanyRepository, auditRepo *repository.AuditRepository) *AdminHandler {
	return &AdminHandler{
		adminSvc:  adminSvc,
		userRepo:  userRepo,
		compRepo:  compRepo,
		auditRepo: auditRepo,
	}
}

func (h *AdminHandler) Routes(r chi.Router) {
	r.Use(middleware.RequireRole("admin"))

	r.Get("/stats", h.GetStats)

	// Users
	r.Get("/users", h.ListUsers)
	r.Put("/users/{user_id}/role", h.UpdateUserRole)
	r.Put("/users/{user_id}/status", h.UpdateUserStatus)

	// Companies
	r.Get("/companies", h.ListCompanies)
	r.Put("/companies/{company_id}/status", h.UpdateCompanyStatus)

	// Audit Logs
	r.Get("/audit-logs", h.ListAuditLogs)
}

func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(middleware.CtxRequestID).(string)
	stats, err := h.adminSvc.GetStats(r.Context())
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, stats, nil, requestID)
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(middleware.CtxRequestID).(string)
	users, _, err := h.userRepo.ListAllUsersPage(r.Context(), pagination.Params{Page: 1, PageSize: 1000}, repository.UserListFilter{})
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	// Redact password hash
	for i := range users {
		users[i].PasswordHash = ""
	}
	pkgresponse.JSON(w, http.StatusOK, users, nil, requestID)
}

func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(middleware.CtxRequestID).(string)
	userID := chi.URLParam(r, "user_id")

	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewBadRequest(err.Error()), requestID)
		return
	}

	err := h.userRepo.UpdateRole(r.Context(), userID, req.Role)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "Role updated"}, nil, requestID)
}

func (h *AdminHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(middleware.CtxRequestID).(string)
	userID := chi.URLParam(r, "user_id")

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewBadRequest(err.Error()), requestID)
		return
	}

	statusStr := "active"
	if !req.IsActive {
		statusStr = "inactive"
	}
	err := h.userRepo.UpdateStatus(r.Context(), userID, statusStr)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "Status updated"}, nil, requestID)
}

func (h *AdminHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(middleware.CtxRequestID).(string)
	companies, _, err := h.compRepo.List(r.Context(), 1000, 0)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, companies, nil, requestID)
}

func (h *AdminHandler) UpdateCompanyStatus(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(middleware.CtxRequestID).(string)
	companyID := chi.URLParam(r, "company_id")
	
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewBadRequest(err.Error()), requestID)
		return
	}

	err := h.compRepo.UpdateStatus(r.Context(), companyID, req.Status)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "Company status updated"}, nil, requestID)
}

func (h *AdminHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(middleware.CtxRequestID).(string)
	
	logs, err := h.auditRepo.ListAll(r.Context(), 1000, 0)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, logs, nil, requestID)
}
