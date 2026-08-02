package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/middleware"
	"backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/service"
)

type AISettingsHandler struct {
	settingsSvc  *service.AISettingsService
	auditSvc     *service.AuditService
	serviceToken string
}

func NewAISettingsHandler(settingsSvc *service.AISettingsService, auditSvc *service.AuditService, serviceToken string) *AISettingsHandler {
	return &AISettingsHandler{settingsSvc: settingsSvc, auditSvc: auditSvc, serviceToken: serviceToken}
}

func (h *AISettingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	settings, err := h.settingsSvc.GetSettings(r.Context())
	if err != nil {
		pkgresponse.Error(w, errors.NewInternal("failed to get AI settings"), "")
		return
	}
	pkgresponse.JSON(w, http.StatusOK, settings, nil, "")
}

func (h *AISettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateAISettingsInput
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		pkgresponse.Error(w, errors.NewBadRequest("invalid JSON body"), "")
		return
	}
	settings, err := h.settingsSvc.UpdateSettings(r.Context(), req)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to update AI settings"), "")
		}
		return
	}
	if h.auditSvc != nil {
		actorID, _ := r.Context().Value(middleware.CtxUserID).(string)
		actorRole, _ := r.Context().Value(middleware.CtxUserRole).(string)
		_ = h.auditSvc.LogAction(r.Context(), service.AuditLogInput{
			ActorUserID: actorID, ActorRole: actorRole, Action: "UPDATE_AI_SETTINGS",
			ResourceType: "system", ResourceID: "ai_settings",
			AfterData: map[string]interface{}{"text_model": req.TextModel, "voice_model": req.VoiceModel, "api_key_changed": req.APIKey != "", "api_key_cleared": req.ClearAPIKey},
			IPAddress: r.RemoteAddr, UserAgent: r.UserAgent(),
		})
	}
	pkgresponse.JSON(w, http.StatusOK, settings, nil, "")
}

func (h *AISettingsHandler) GetRuntime(w http.ResponseWriter, r *http.Request) {
	if h.serviceToken == "" || r.Header.Get("X-Internal-Service-Token") != h.serviceToken {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	settings, err := h.settingsSvc.GetRuntimeConfig(r.Context())
	if err != nil {
		http.Error(w, "AI settings unavailable", http.StatusServiceUnavailable)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, settings, nil, "")
}
