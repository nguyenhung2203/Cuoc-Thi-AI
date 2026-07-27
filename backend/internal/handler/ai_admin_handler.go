package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/models"
	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/service"
)

type AIAdminHandler struct {
	promptSvc *service.PromptService
}

func NewAIAdminHandler(promptSvc *service.PromptService) *AIAdminHandler {
	return &AIAdminHandler{promptSvc: promptSvc}
}

// CreatePromptTemplate handles POST /api/v1/admin/ai-prompts
func (h *AIAdminHandler) CreatePromptTemplate(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	
	// Assuming an admin middleware already ran and authenticated the user.
	var req models.AIPromptTemplate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid JSON body", []string{err.Error()}), requestID)
		return
	}

	// For simplicity, skip detailed validation here. In production, use a dedicated request struct.
	if req.Name == "" || req.Content == "" || req.Model == "" {
		pkgresponse.Error(w, apierrors.NewValidation("missing required fields", []string{"name, content, model are required"}), requestID)
		return
	}

	// Always active on creation for simplicity
	req.IsActive = true

	created, err := h.promptSvc.CreateNewVersion(r.Context(), &req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusCreated, created, nil, requestID)
}

// ListTemplates handles GET /api/v1/admin/ai-prompts
func (h *AIAdminHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	templates, err := h.promptSvc.ListAllTemplates(r.Context())
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, templates, nil, requestID)
}
