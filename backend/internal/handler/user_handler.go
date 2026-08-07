package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/pagination"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/repository"
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
	r.With(middleware.AuthMiddleware(h.jwtSecret), middleware.RoleMiddleware("admin")).Get("/admin/reports", h.GetReports)
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

func (h *UserHandler) GetReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.userSvc.GetReports(r.Context())
	if err != nil {
		pkgresponse.Error(w, errors.NewInternal("failed to get reports"), "")
		return
	}
	pkgresponse.JSON(w, http.StatusOK, reports, nil, "")
}

func (h *UserHandler) listUsers(w http.ResponseWriter, r *http.Request, pending bool) {
	p := pagination.FromRequest(r)
	filter := repository.UserListFilter{Search: strings.TrimSpace(r.URL.Query().Get("search")), Role: strings.TrimSpace(r.URL.Query().Get("role"))}
	result, err := h.userSvc.ListUsersPage(r.Context(), p, filter, pending)
	if err != nil {
		pkgresponse.Error(w, errors.NewInternal("failed to list users"), "")
		return
	}
	all, pendingTotal, err := h.userSvc.UserCounts(r.Context())
	if err != nil {
		pkgresponse.Error(w, errors.NewInternal("failed to count users"), "")
		return
	}
	meta := &pkgresponse.Meta{Page: p.Page, PageSize: p.PageSize, Total: result.Total, TotalPages: pagination.CalcTotalPages(result.Total, p.PageSize)}
	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{"users": result.Users, "counts": map[string]int{"all": all, "pending": pendingTotal}}, meta, "")
}

func (h *UserHandler) ListPendingUsers(w http.ResponseWriter, r *http.Request) {
	h.listUsers(w, r, true)
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

func (h *UserHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) { h.listUsers(w, r, false) }

func (h *UserHandler) ListAllCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.companySvc.ListAllCompanies(r.Context())
	if err != nil {
		pkgresponse.Error(w, errors.NewInternal("failed to list all companies"), "")
		return
	}
	pkgresponse.JSON(w, http.StatusOK, companies, nil, "")
}
