package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/response"
	"backend/internal/service"
)

type TranscriptHandler struct {
	svc *service.TranscriptService
}

func NewTranscriptHandler(svc *service.TranscriptService) *TranscriptHandler {
	return &TranscriptHandler{svc: svc}
}

func (h *TranscriptHandler) Routes(r chi.Router) {
	// Need interview:update to push transcripts, or role:member
	r.With(middleware.RequirePermission("interview:update")).Post("/", h.PushTranscript)
	r.With(middleware.RequirePermission("interview:read")).Get("/", h.ListTranscripts)
	r.With(middleware.RequirePermission("interview:update")).Put("/{transcript_id}", h.EditTranscript)
}

func (h *TranscriptHandler) PushTranscript(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req service.PushTranscriptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierrors.NewBadRequest("invalid request body"), requestID)
		return
	}

	t, err := h.svc.PushTranscript(r.Context(), interviewID, companyID, req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusCreated, t, nil, requestID)
}

func (h *TranscriptHandler) ListTranscripts(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	items, err := h.svc.ListTranscripts(r.Context(), interviewID, companyID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusOK, items, nil, requestID)
}

type EditTranscriptRequest struct {
	EditedContent string `json:"edited_content"`
}

func (h *TranscriptHandler) EditTranscript(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	transcriptID := chi.URLParam(r, "transcript_id")
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req EditTranscriptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierrors.NewBadRequest("invalid request body"), requestID)
		return
	}

	err := h.svc.EditTranscript(r.Context(), transcriptID, interviewID, companyID, req.EditedContent, userID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "transcript updated"}, nil, requestID)
}


