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

type ReportHandler struct {

	svc *service.ReportService
}

func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) Routes(r chi.Router) {
	r.Get("/", h.GetReport)
	r.Put("/decision", h.SaveDecision)
}

func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	interviewID := chi.URLParam(r, "interview_id")

	report, err := h.svc.GetOrGenerateReport(r.Context(), interviewID, companyID)

	reportSvc *service.ReportService
}

func NewReportHandler(reportSvc *service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportSvc: reportSvc,
	}
}

func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	report, err := h.reportSvc.GetReport(r.Context(), companyID, interviewID)

	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}


	pkgresponse.JSON(w, http.StatusOK, report, nil, requestID)
}

func (h *ReportHandler) SaveDecision(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	interviewID := chi.URLParam(r, "interview_id")

	var req request.DecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeServiceError(w, apierrors.NewBadRequest("invalid json payload"), requestID)
		return
	}
	if errs := validator.Validate(&req); len(errs) > 0 {
		writeServiceError(w, apierrors.NewBadRequest("validation failed: "+errs[0]), requestID)
		return
	}

	err := h.svc.UpdateDecision(r.Context(), interviewID, companyID, req.Decision, req.Comment)

	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{
		"report": report,
	}, nil, requestID)
}

func (h *ReportHandler) OverrideDecision(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req request.UpdateReportDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid json payload", []string{err.Error()}), requestID)
		return
	}

	if errs := validator.Validate(&req); errs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("payload", errs), requestID)
		return
	}

	if err := h.reportSvc.OverrideDecision(r.Context(), companyID, interviewID, req.Decision, req.Comment); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{
		"message": "decision updated",
	}, nil, requestID)
}

func (h *ReportHandler) RetryReport(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	interviewID := chi.URLParam(r, "interview_id")
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	report, err := h.reportSvc.RetryReport(r.Context(), companyID, interviewID, userID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}


	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "decision saved"}, nil, requestID)

	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{
		"report": report,
	}, nil, requestID)
}
