package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/dto/request"
	"backend/internal/middleware"
	"backend/internal/models"
	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/pkg/validator"
	"backend/internal/service"
)

type RubricHandler struct {
	rubricSvc *service.RubricService
}

func NewRubricHandler(rubricSvc *service.RubricService) *RubricHandler {
	return &RubricHandler{
		rubricSvc: rubricSvc,
	}
}

func (h *RubricHandler) CreateRubric(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	userID, _ := r.Context().Value(middleware.CtxUserID).(string) // assuming middleware sets this
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req request.CreateRubricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("payload", []string{"invalid json payload"}), requestID)
		return
	}

	if errs := validator.Validate(req); len(errs) > 0 {
		pkgresponse.Error(w, apierrors.NewValidation("payload", errs), requestID)
		return
	}

	rubric := &models.Rubric{
		CompanyID: companyID,
		Name:      req.Name,
		CreatedBy: userID,
	}
	if req.JobID != nil {
		rubric.JobID = sql.NullString{String: *req.JobID, Valid: true}
	}
	if req.Description != nil {
		rubric.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	var criteria []models.RubricCriteria
	for _, cReq := range req.Criteria {
		c := models.RubricCriteria{
			Name:         cReq.Name,
			Weight:       cReq.Weight,
			MinScore:     cReq.MinScore,
			MaxScore:     cReq.MaxScore,
			ScoringGuide: models.JSONB([]byte(`"` + cReq.ScoringGuide + `"`)), // simple wrap for string
			OrderIndex:   cReq.OrderIndex,
		}
		if cReq.Description != nil {
			c.Description = sql.NullString{String: *cReq.Description, Valid: true}
		}
		criteria = append(criteria, c)
	}

	if err := h.rubricSvc.CreateRubric(r.Context(), rubric, criteria); err != nil {
		if appErr, ok := apierrors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, requestID)
		} else {
			pkgresponse.Error(w, apierrors.NewInternal(err.Error()), requestID)
		}
		return
	}

	// We can return the created rubric here, but for now just success
	pkgresponse.JSON(w, http.StatusCreated, map[string]interface{}{
		"id": rubric.ID,
	}, nil, requestID)
}

func (h *RubricHandler) UpdateRubric(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	rubricID := chi.URLParam(r, "rubric_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req request.CreateRubricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("payload", []string{"invalid json payload"}), requestID)
		return
	}
	if errs := validator.Validate(req); len(errs) > 0 {
		pkgresponse.Error(w, apierrors.NewValidation("payload", errs), requestID)
		return
	}

	rubric := &models.Rubric{
		ID:        rubricID,
		CompanyID: companyID,
		Name:      req.Name,
	}
	if req.JobID != nil {
		rubric.JobID = sql.NullString{String: *req.JobID, Valid: true}
	}
	if req.Description != nil {
		rubric.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	var criteria []models.RubricCriteria
	for _, cReq := range req.Criteria {
		c := models.RubricCriteria{
			Name:         cReq.Name,
			Weight:       cReq.Weight,
			MinScore:     cReq.MinScore,
			MaxScore:     cReq.MaxScore,
			ScoringGuide: models.JSONB([]byte(`"` + cReq.ScoringGuide + `"`)),
			OrderIndex:   cReq.OrderIndex,
		}
		if cReq.Description != nil {
			c.Description = sql.NullString{String: *cReq.Description, Valid: true}
		}
		criteria = append(criteria, c)
	}

	if err := h.rubricSvc.UpdateRubric(r.Context(), rubric, criteria); err != nil {
		if err == sql.ErrNoRows {
			pkgresponse.Error(w, apierrors.NewNotFound("rubric"), requestID)
			return
		}
		if appErr, ok := apierrors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, requestID)
		} else {
			pkgresponse.Error(w, apierrors.NewInternal(err.Error()), requestID)
		}
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{"id": rubricID}, nil, requestID)
}

func (h *RubricHandler) ListCompanyRubrics(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	jobID := r.URL.Query().Get("job_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	rubrics, err := h.rubricSvc.ListCompanyRubrics(r.Context(), companyID, jobID)
	if err != nil {
		if appErr, ok := apierrors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, requestID)
		} else {
			pkgresponse.Error(w, apierrors.NewInternal(err.Error()), requestID)
		}
		return
	}

	pkgresponse.JSON(w, http.StatusOK, rubrics, nil, requestID)
}

func (h *RubricHandler) GetRubric(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	rubricID := chi.URLParam(r, "rubric_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	rubric, criteria, err := h.rubricSvc.GetRubricByID(r.Context(), companyID, rubricID)
	if err != nil {
		if appErr, ok := apierrors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, requestID)
		} else {
			pkgresponse.Error(w, apierrors.NewInternal(err.Error()), requestID)
		}
		return
	}

	res := map[string]interface{}{
		"rubric":   rubric,
		"criteria": criteria,
	}
	pkgresponse.JSON(w, http.StatusOK, res, nil, requestID)
}

func (h *RubricHandler) DeleteRubric(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	rubricID := chi.URLParam(r, "rubric_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	if err := h.rubricSvc.DeleteRubric(r.Context(), companyID, rubricID); err != nil {
		if appErr, ok := apierrors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, requestID)
		} else {
			pkgresponse.Error(w, apierrors.NewInternal(err.Error()), requestID)
		}
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "deleted successfully"}, nil, requestID)
}
