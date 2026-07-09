package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/dto/request"
	"backend/internal/middleware"
	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/pkg/validator"
	"backend/internal/service"
)

type AiHandler struct {
	scoreSvc      *service.ScoreService
	reportSvc     *service.ReportService
	suggestionSvc *service.SuggestionService
}

func NewAiHandler(scoreSvc *service.ScoreService, reportSvc *service.ReportService, suggestionSvc *service.SuggestionService) *AiHandler {
	return &AiHandler{
		scoreSvc:      scoreSvc,
		reportSvc:     reportSvc,
		suggestionSvc: suggestionSvc,
	}
}

func (h *AiHandler) ScoreAnswer(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")

	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req request.ScoreAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid json payload", []string{err.Error()}), requestID)
		return
	}

	if errs := validator.Validate(&req); errs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("payload", errs), requestID)
		return
	}

	scores, err := h.scoreSvc.ScoreAnswer(r.Context(), companyID, interviewID, req.TranscriptIDs, req.CriterionIDs)
	if err != nil {
		if appErr, ok := apierrors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, requestID)
		} else {
			pkgresponse.Error(w, apierrors.NewInternal(err.Error()), requestID)
		}
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{
		"scores": scores,
	}, nil, requestID)
}

func (h *AiHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)
	userID, _ := r.Context().Value(middleware.CtxUserID).(string) // Assume generated_by is current user

	var req request.GenerateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid json payload", []string{err.Error()}), requestID)
		return
	}

	if errs := validator.Validate(&req); errs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("payload", errs), requestID)
		return
	}

	// This is a synchronous call. For a production app, we might want to run this in a goroutine
	// and return a 202 Accepted, then update the interview status when done.
	report, err := h.reportSvc.GenerateReport(r.Context(), companyID, interviewID, req.JobID, userID)
	if err != nil {
		if appErr, ok := apierrors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, requestID)
		} else {
			pkgresponse.Error(w, apierrors.NewInternal(err.Error()), requestID)
		}
		return
	}

	pkgresponse.JSON(w, http.StatusAccepted, map[string]interface{}{
		"report": report,
		"status": "generating",
	}, nil, requestID)
}

func (h *AiHandler) SuggestFollowUp(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req request.SuggestFollowUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid json payload", []string{err.Error()}), requestID)
		return
	}

	// Audit log
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("ai:suggest_follow_up", "ai_suggestion", interviewID, companyID, nil, map[string]string{
			"focus": req.Focus,
		})
	}

	result, err := h.suggestionSvc.SuggestFollowUp(r.Context(), companyID, interviewID, req.LastTranscriptID, req.Focus)
	if err != nil {
		if appErr, ok := apierrors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, requestID)
		} else {
			pkgresponse.Error(w, apierrors.NewInternal(err.Error()), requestID)
		}
		return
	}

	pkgresponse.JSON(w, http.StatusOK, result, nil, requestID)
}
