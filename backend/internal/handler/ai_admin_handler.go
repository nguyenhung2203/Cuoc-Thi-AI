package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/models"
	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/repository"
	"backend/internal/service"
)

type AIAdminHandler struct {
	promptSvc *service.PromptService
	logSvc    *service.AILogService
}

func NewAIAdminHandler(promptSvc *service.PromptService, logSvc *service.AILogService) *AIAdminHandler {
	return &AIAdminHandler{promptSvc: promptSvc, logSvc: logSvc}
}

// ListAILogs handles GET /api/v1/admin/ai-logs?interview_id=&job_id=&candidate_id=&page=&page_size=
// Returns AI request logs for debugging and cost auditing (admin only).
func (h *AIAdminHandler) ListAILogs(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 50
	}

	logs, err := h.logSvc.List(r.Context(), repository.AILogFilter{
		InterviewID: q.Get("interview_id"),
		JobID:       q.Get("job_id"),
		CandidateID: q.Get("candidate_id"),
		Limit:       pageSize,
		Offset:      (page - 1) * pageSize,
	})
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, logs, nil, requestID)
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
