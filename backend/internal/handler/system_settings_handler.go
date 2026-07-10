package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/middleware"
	"backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/service"
)

type SystemSettingsHandler struct {
	settingsSvc *service.SystemSettingsService
	auditSvc    *service.AuditService
}

func NewSystemSettingsHandler(settingsSvc *service.SystemSettingsService, auditSvc *service.AuditService) *SystemSettingsHandler {
	return &SystemSettingsHandler{settingsSvc: settingsSvc, auditSvc: auditSvc}
}

func (h *SystemSettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.settingsSvc.GetSettings(r.Context())
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to get system settings"), "")
		}
		return
	}
	pkgresponse.JSON(w, http.StatusOK, settings, nil, "")
}

func (h *SystemSettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, errors.NewBadRequest("invalid JSON body"), "")
		return
	}

	if err := h.settingsSvc.UpdateSettings(r.Context(), req); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to update system settings"), "")
		}
		return
	}

	actorID, _ := r.Context().Value(middleware.CtxUserID).(string)
	actorRole, _ := r.Context().Value(middleware.CtxUserRole).(string)
	if h.auditSvc != nil {
		h.auditSvc.LogAction(r.Context(), service.AuditLogInput{
			ActorUserID:  actorID,
			ActorRole:    actorRole,
			Action:       "UPDATE_SYSTEM_SETTINGS",
			ResourceType: "system",
			ResourceID:   "system_settings",
			AfterData:    req,
			IPAddress:    r.RemoteAddr,
			UserAgent:    r.UserAgent(),
		})
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{"status": "success", "settings": req}, nil, "")
}
