package handler

import (
	"net/http"

	"backend/internal/pkg/response"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/pagination"
	"backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type PublicJobHandler struct {
	jobService *service.JobService
}

func NewPublicJobHandler(jobService *service.JobService) *PublicJobHandler {
	return &PublicJobHandler{jobService: jobService}
}

func (h *PublicJobHandler) Routes(r chi.Router) {
	r.Get("/companies/{companyID}/jobs", h.ListJobs)
	r.Get("/companies/{companyID}/jobs/{jobID}", h.GetJob)
	r.Get("/all-jobs", h.ListAllJobs)
}

func (h *PublicJobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "companyID")
	if companyID == "" {
		response.Error(w, errors.NewBadRequest("missing company ID"), "")
		return
	}

	status := "open" // Only show open jobs publicly
	keyword := r.URL.Query().Get("keyword")
	p := pagination.FromRequest(r)

	jobs, total, err := h.jobService.List(r.Context(), companyID, status, keyword, p)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(w, appErr, "")
		} else {
			response.Error(w, errors.NewInternal("failed to list public jobs"), "")
		}
		return
	}

	meta := &response.Meta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: pagination.CalcTotalPages(total, p.PageSize),
	}
	response.JSON(w, http.StatusOK, jobs, meta, "")
}

func (h *PublicJobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "companyID")
	jobID := chi.URLParam(r, "jobID")
	
	if companyID == "" || jobID == "" {
		response.Error(w, errors.NewBadRequest("missing company or job ID"), "")
		return
	}

	job, err := h.jobService.GetByID(r.Context(), companyID, jobID)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(w, appErr, "")
		} else {
			response.Error(w, errors.NewInternal("failed to get job details"), "")
		}
		return
	}

	if job.Status != "open" {
		response.Error(w, errors.NewNotFound("job not found or closed"), "")
		return
	}

	response.JSON(w, http.StatusOK, job, nil, "")
}

func (h *PublicJobHandler) ListAllJobs(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	p := pagination.FromRequest(r)

	jobs, total, err := h.jobService.ListAllOpen(r.Context(), keyword, p)
	if err != nil {
		response.Error(w, errors.NewInternal("failed to list jobs"), "")
		return
	}

	meta := &response.Meta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: pagination.CalcTotalPages(total, p.PageSize),
	}
	response.JSON(w, http.StatusOK, jobs, meta, "")
}
