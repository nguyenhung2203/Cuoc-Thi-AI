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
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "decision saved"}, nil, requestID)
}
