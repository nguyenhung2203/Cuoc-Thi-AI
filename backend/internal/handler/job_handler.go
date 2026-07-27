package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/dto/request"
	"backend/internal/dto/response"
	"backend/internal/middleware"
	"backend/internal/models"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/pagination"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/pkg/validator"
	"backend/internal/service"
)

// JobHandler wires the JobService to HTTP endpoints.
type JobHandler struct {
	svc *service.JobService
}

// NewJobHandler constructs a JobHandler.
func NewJobHandler(svc *service.JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

// Routes registers all job endpoints on r.
// Expected mount point: /companies/{company_id}
//
//	GET    /jobs              → List
//	POST   /jobs              → Create
//	GET    /jobs/{job_id}     → GetByID
//	PUT    /jobs/{job_id}     → Update
//	DELETE /jobs/{job_id}     → Delete
func (h *JobHandler) Routes(r chi.Router) {
	r.Route("/jobs", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Route("/{job_id}", func(r chi.Router) {
			r.Get("/", h.GetByID)
			r.Put("/", h.Update)
			r.Delete("/", h.Delete)
			r.Post("/analyze", h.Analyze)
			r.Route("/ai", func(r chi.Router) {
				r.Post("/generate-questions", h.GenerateQuestions)
			})
		})
	})
}

// List handles GET /companies/{company_id}/jobs
func (h *JobHandler) List(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)

	q := r.URL.Query()
	status := q.Get("status")
	keyword := q.Get("keyword")
	p := pagination.FromRequest(r)

	jobs, total, err := h.svc.List(r.Context(), companyID, status, keyword, p)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	items := make([]response.JobListItem, 0, len(jobs))
	for _, job := range jobs {
		candidateCount, interviewCount, _ := h.svc.GetStats(r.Context(), job.ID)
		items = append(items, toJobListItem(job, candidateCount, interviewCount))
	}

	meta := &pkgresponse.Meta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: pagination.CalcTotalPages(total, p.PageSize),
	}
	pkgresponse.JSON(w, http.StatusOK, items, meta, requestID)
}

// Create handles POST /companies/{company_id}/jobs
func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)

	var req request.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid JSON body", []string{err.Error()}), requestID)
		return
	}
	if msgs := validator.Validate(&req); msgs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("validation failed", msgs), requestID)
		return
	}

	job, err := h.svc.Create(r.Context(), companyID, userID, &req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Audit log job creation
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("job:create", "job", job.ID, companyID, nil, job)
	}

	candidateCount, interviewCount, _ := h.svc.GetStats(r.Context(), job.ID)
	stats := response.JobStats{CandidateCount: candidateCount, InterviewCount: interviewCount}
	pkgresponse.JSON(w, http.StatusCreated, toJobDetail(*job, stats), nil, requestID)
}

// GetByID handles GET /companies/{company_id}/jobs/{job_id}
func (h *JobHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	jobID := chi.URLParam(r, "job_id")

	job, err := h.svc.GetByID(r.Context(), companyID, jobID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	candidateCount, interviewCount, _ := h.svc.GetStats(r.Context(), job.ID)
	stats := response.JobStats{CandidateCount: candidateCount, InterviewCount: interviewCount}
	pkgresponse.JSON(w, http.StatusOK, toJobDetail(*job, stats), nil, requestID)
}

// Update handles PUT /companies/{company_id}/jobs/{job_id}
func (h *JobHandler) Update(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	jobID := chi.URLParam(r, "job_id")

	var req request.UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid JSON body", []string{err.Error()}), requestID)
		return
	}
	if msgs := validator.Validate(&req); msgs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("validation failed", msgs), requestID)
		return
	}

	// Fetch old job before updating (for audit diff)
	oldJob, _ := h.svc.GetByID(r.Context(), companyID, jobID)

	job, err := h.svc.Update(r.Context(), companyID, jobID, &req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Audit log job update
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("job:update", "job", jobID, companyID, oldJob, job)
	}

	candidateCount, interviewCount, _ := h.svc.GetStats(r.Context(), job.ID)
	stats := response.JobStats{CandidateCount: candidateCount, InterviewCount: interviewCount}
	pkgresponse.JSON(w, http.StatusOK, toJobDetail(*job, stats), nil, requestID)
}

// Delete handles DELETE /companies/{company_id}/jobs/{job_id}
func (h *JobHandler) Delete(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	jobID := chi.URLParam(r, "job_id")

	if err := h.svc.Delete(r.Context(), companyID, jobID); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Audit log job deletion
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("job:delete", "job", jobID, companyID, map[string]string{"job_id": jobID}, nil)
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "job deleted"}, nil, requestID)
}

// Analyze handles POST /companies/{company_id}/jobs/{job_id}/analyze
func (h *JobHandler) Analyze(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	jobID := chi.URLParam(r, "job_id")

	if err := h.svc.Analyze(r.Context(), companyID, jobID); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "job analyzed successfully"}, nil, requestID)
}

// GenerateQuestions handles POST /companies/{company_id}/jobs/{job_id}/ai/generate-questions
func (h *JobHandler) GenerateQuestions(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "company_id")
	jobID := chi.URLParam(r, "job_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req request.GenerateQuestionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("payload", []string{"invalid json payload"}), requestID)
		return
	}

	questions, err := h.svc.GenerateQuestions(r.Context(), companyID, jobID, &req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{
		"questions": questions,
		"count":     len(questions),
	}, nil, requestID)
}

// ---------------------------------------------------------------------------
// Mapping helpers
// ---------------------------------------------------------------------------

// toJobListItem maps a models.Job to the list-view DTO.
func toJobListItem(job models.Job, candidateCount, interviewCount int) response.JobListItem {
	return response.JobListItem{
		ID:             job.ID,
		Title:          job.Title,
		Department:     job.Department.String,
		Level:          job.Level.String,
		Status:         string(job.Status),
		CandidateCount: candidateCount,
		InterviewCount: interviewCount,
		CreatedAt:      job.CreatedAt,
	}
}

// toJobDetail maps a models.Job plus pre-fetched stats to the detail DTO.
func toJobDetail(job models.Job, stats response.JobStats) response.JobDetail {
	return response.JobDetail{
		ID:             job.ID,
		CompanyID:      job.CompanyID,
		Title:          job.Title,
		Department:     job.Department.String,
		Level:          job.Level.String,
		Location:       job.Location.String,
		EmploymentType: job.EmploymentType.String,
		SalaryMin:      job.SalaryMin.Float64,
		SalaryMax:      job.SalaryMax.Float64,
		Currency:       job.Currency.String,
		Description:    job.Description,
		Requirements:   job.Requirements.String,
		Benefits:       job.Benefits.String,
		Status:         string(job.Status),
		AISummary:      job.AISummary.String,
		AIAnalysisJSON: job.AIAnalysisJSON,
		Stats:          stats,
		CreatedAt:      job.CreatedAt,
		UpdatedAt:      job.UpdatedAt,
	}
}

// ---------------------------------------------------------------------------
// Shared handler utilities
// ---------------------------------------------------------------------------

// getRequestID reads the request_id string injected by AuthMiddleware.
func getRequestID(r *http.Request) string {
	id, _ := r.Context().Value(middleware.CtxRequestID).(string)
	return id
}

// writeServiceError casts err to *AppError when possible and writes the
// appropriate response, falling back to 500 for unexpected errors.
func writeServiceError(w http.ResponseWriter, err error, requestID string) {
	if appErr, ok := apierrors.IsAppError(err); ok {
		pkgresponse.Error(w, appErr, requestID)
		return
	}
	pkgresponse.Error(w, apierrors.NewInternal("unexpected error"), requestID)
}
