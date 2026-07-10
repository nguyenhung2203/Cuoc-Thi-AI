package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/service"
)

type AuditHandler struct {
	auditSvc *service.AuditService
}

func NewAuditHandler(auditSvc *service.AuditService) *AuditHandler {
	return &AuditHandler{auditSvc: auditSvc}
}

func (h *AuditHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	resourceType := r.URL.Query().Get("resource_type")
	actorUserID := r.URL.Query().Get("actor_user_id")
	page := 1
	pageSize := 20

	logs, err := h.auditSvc.ListByCompany(r.Context(), companyID, resourceType, actorUserID, page, pageSize)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, logs, nil, requestID)
}

func (h *AuditHandler) LogActionDirect(w http.ResponseWriter, r *http.Request) {
	// Internal: not exposed via router, for use by other handlers through shared AuditService
	pkgresponse.Error(w, apierrors.NewValidation("direct", []string{"use audit service from within handlers"}), "")
}

func (h *AuditHandler) ListAllGlobalLogs(w http.ResponseWriter, r *http.Request) {
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	page := 1
	pageSize := 50

	logs, err := h.auditSvc.ListAll(r.Context(), page, pageSize)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, logs, nil, requestID)
}
