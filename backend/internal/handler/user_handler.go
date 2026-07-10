package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	"backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/service"
)

type UserHandler struct {
	userSvc    *service.UserService
	companySvc *service.CompanyService
	auditSvc   *service.AuditService
	jwtSecret  string
}

func NewUserHandler(userSvc *service.UserService, companySvc *service.CompanyService, auditSvc *service.AuditService, jwtSecret string) *UserHandler {
	return &UserHandler{userSvc: userSvc, companySvc: companySvc, auditSvc: auditSvc, jwtSecret: jwtSecret}
}

func (h *UserHandler) Routes(r chi.Router) {
	// Admin-only routes
	r.With(middleware.AuthMiddleware(h.jwtSecret), middleware.RoleMiddleware("admin")).Get("/admin/users", h.ListAllUsers)
	r.With(middleware.AuthMiddleware(h.jwtSecret), middleware.RoleMiddleware("admin")).Get("/admin/users/pending", h.ListPendingUsers)
	r.With(middleware.AuthMiddleware(h.jwtSecret), middleware.RoleMiddleware("admin")).Put("/admin/users/{user_id}/approve", h.ApproveUser)
	r.With(middleware.AuthMiddleware(h.jwtSecret), middleware.RoleMiddleware("admin")).Put("/admin/users/{user_id}/status", h.UpdateUserStatus)
	r.With(middleware.AuthMiddleware(h.jwtSecret), middleware.RoleMiddleware("admin")).Get("/admin/dashboard-stats", h.GetDashboardStats)
	r.With(middleware.AuthMiddleware(h.jwtSecret), middleware.RoleMiddleware("admin")).Get("/admin/companies", h.ListAllCompanies)
}

func (h *UserHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.userSvc.GetDashboardStats(r.Context())
	if err != nil {
		pkgresponse.Error(w, errors.NewInternal("failed to get dashboard stats"), "")
		return
	}
	pkgresponse.JSON(w, http.StatusOK, stats, nil, "")
}

func (h *UserHandler) ListPendingUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userSvc.ListPendingUsers(r.Context())
	if err != nil {
		pkgresponse.Error(w, errors.NewInternal("failed to list users"), "")
		return
	}
	pkgresponse.JSON(w, http.StatusOK, users, nil, "")
}

func (h *UserHandler) ApproveUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		pkgresponse.Error(w, errors.NewBadRequest("user_id is required"), "")
		return
	}
	if err := h.userSvc.ApproveUser(r.Context(), userID); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to approve user"), "")
		}
		return
	}

	actorID, _ := r.Context().Value(middleware.CtxUserID).(string)
	actorRole, _ := r.Context().Value(middleware.CtxUserRole).(string)
	if h.auditSvc != nil {
		h.auditSvc.LogAction(r.Context(), service.AuditLogInput{
			ActorUserID:  actorID,
			ActorRole:    actorRole,
			Action:       "APPROVE",
			ResourceType: "user",
			ResourceID:   userID,
			IPAddress:    r.RemoteAddr,
			UserAgent:    r.UserAgent(),
		})
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"status": "approved"}, nil, "")
}

func (h *UserHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		pkgresponse.Error(w, errors.NewBadRequest("user_id is required"), "")
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		pkgresponse.Error(w, errors.NewBadRequest("invalid or missing status"), "")
		return
	}
	if err := h.userSvc.UpdateUserStatus(r.Context(), userID, req.Status); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to update user status"), "")
		}
		return
	}

	actorID, _ := r.Context().Value(middleware.CtxUserID).(string)
	actorRole, _ := r.Context().Value(middleware.CtxUserRole).(string)
	actionName := "UPDATE_STATUS_" + req.Status
	if req.Status == "blocked" {
		actionName = "BLOCK"
	} else if req.Status == "active" {
		actionName = "ACTIVATE"
	}
	if h.auditSvc != nil {
		h.auditSvc.LogAction(r.Context(), service.AuditLogInput{
			ActorUserID:  actorID,
			ActorRole:    actorRole,
			Action:       actionName,
			ResourceType: "user",
			ResourceID:   userID,
			IPAddress:    r.RemoteAddr,
			UserAgent:    r.UserAgent(),
		})
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"status": req.Status}, nil, "")
}

func (h *UserHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userSvc.ListAllUsers(r.Context())
	if err != nil {
		pkgresponse.Error(w, errors.NewInternal("failed to list all users"), "")
		return
	}
	pkgresponse.JSON(w, http.StatusOK, users, nil, "")
}

func (h *UserHandler) ListAllCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.companySvc.ListAllCompanies(r.Context())
	if err != nil {
		pkgresponse.Error(w, errors.NewInternal("failed to list all companies"), "")
		return
	}
	pkgresponse.JSON(w, http.StatusOK, companies, nil, "")
}
