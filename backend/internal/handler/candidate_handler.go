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

// CandidateHandler wires the CandidateService to HTTP endpoints.
type CandidateHandler struct {
	svc     *service.CandidateService
	aiSvc   *service.AIService
	fileSvc *service.FileService
}

// NewCandidateHandler constructs a CandidateHandler.
func NewCandidateHandler(svc *service.CandidateService, aiSvc *service.AIService, fileSvc *service.FileService) *CandidateHandler {
	return &CandidateHandler{svc: svc, aiSvc: aiSvc, fileSvc: fileSvc}
}

// Routes registers all candidate endpoints on r.
// Expected mount point: /companies/{company_id}
//
//	GET    /candidates                              → List
//	POST   /candidates                              → Create
//	GET    /candidates/{candidate_id}               → GetByID
//	PUT    /candidates/{candidate_id}               → Update
//	DELETE /candidates/{candidate_id}               → Delete
//	POST   /candidates/{candidate_id}/parse-cv      → ParseCV
//	GET    /jobs/{job_id}/candidates                → ListByJob
//	POST   /jobs/{job_id}/candidates/{candidate_id}/assign   → AssignToJob
//	PUT    /jobs/{job_id}/candidates/{candidate_id}/pipeline → UpdatePipeline
//	DELETE /jobs/{job_id}/candidates/{candidate_id}/unassign → UnassignFromJob
func (h *CandidateHandler) Routes(r chi.Router) {
	r.Route("/candidates", func(r chi.Router) {
		r.With(middleware.RequirePermission("candidate:read")).Get("/", h.List)
		r.With(middleware.RequirePermission("candidate:create")).Post("/", h.Create)
		r.Route("/{candidate_id}", func(r chi.Router) {
			r.With(middleware.RequirePermission("candidate:read")).Get("/", h.GetByID)
			r.With(middleware.RequirePermission("candidate:update")).Post("/cv", h.UploadCV)
			r.With(middleware.RequirePermission("candidate:update")).Post("/parse-cv", h.ParseCV)
			r.With(middleware.RequirePermission("candidate:update")).Put("/", h.Update)
			r.With(middleware.RequirePermission("candidate:delete")).Delete("/", h.Delete)
		})
	})

	r.Route("/jobs/{job_id}/candidates", func(r chi.Router) {
		r.With(middleware.RequirePermission("candidate:read")).Get("/", h.ListByJob)
		r.With(middleware.RequirePermission("candidate:update")).Post("/{candidate_id}/assign", h.AssignToJob)
		r.With(middleware.RequirePermission("candidate:update")).Put("/{candidate_id}/pipeline", h.UpdatePipeline)
		// Unassign removes the job link, not the candidate — update, not delete.
		r.With(middleware.RequirePermission("candidate:update")).Delete("/{candidate_id}/unassign", h.UnassignFromJob)
	})
}

// List handles GET /companies/{company_id}/candidates
func (h *CandidateHandler) List(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)

	q := r.URL.Query()
	jobID := q.Get("job_id")
	status := q.Get("status")
	keyword := q.Get("keyword")
	p := pagination.FromRequest(r)

	candidates, total, err := h.svc.List(r.Context(), companyID, jobID, status, keyword, p)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	items := make([]response.CandidateListItem, 0, len(candidates))
	for _, c := range candidates {
		items = append(items, toCandidateListItem(c))
	}

	meta := &pkgresponse.Meta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: pagination.CalcTotalPages(total, p.PageSize),
	}
	pkgresponse.JSON(w, http.StatusOK, items, meta, requestID)
}

// Create handles POST /companies/{company_id}/candidates
func (h *CandidateHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)

	var req request.CreateCandidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid JSON body", []string{err.Error()}), requestID)
		return
	}
	if msgs := validator.Validate(&req); msgs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("validation failed", msgs), requestID)
		return
	}

	candidate, err := h.svc.Create(r.Context(), companyID, userID, &req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Audit log candidate creation
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("candidate:create", "candidate", candidate.ID, companyID, nil, map[string]interface{}{
			"full_name": candidate.FullName,
			"email":     candidate.Email,
		})
	}

	pkgresponse.JSON(w, http.StatusCreated, toCandidateDetail(*candidate), nil, requestID)
}

// GetByID handles GET /companies/{company_id}/candidates/{candidate_id}
func (h *CandidateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	candidateID := chi.URLParam(r, "candidate_id")

	candidate, err := h.svc.GetByID(r.Context(), companyID, candidateID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, toCandidateDetail(*candidate), nil, requestID)
}

// Update handles PUT /companies/{company_id}/candidates/{candidate_id}
func (h *CandidateHandler) Update(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	candidateID := chi.URLParam(r, "candidate_id")

	var req request.UpdateCandidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid JSON body", []string{err.Error()}), requestID)
		return
	}
	if msgs := validator.Validate(&req); msgs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("validation failed", msgs), requestID)
		return
	}

	// Fetch old candidate before updating (for audit diff)
	oldCandidate, _ := h.svc.GetByID(r.Context(), companyID, candidateID)

	candidate, err := h.svc.Update(r.Context(), companyID, candidateID, &req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Audit log candidate update
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("candidate:update", "candidate", candidateID, companyID, oldCandidate, candidate)
	}

	pkgresponse.JSON(w, http.StatusOK, toCandidateDetail(*candidate), nil, requestID)
}

// Delete handles DELETE /companies/{company_id}/candidates/{candidate_id}
func (h *CandidateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	candidateID := chi.URLParam(r, "candidate_id")

	if err := h.svc.Delete(r.Context(), companyID, candidateID); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Audit log candidate deletion
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("candidate:delete", "candidate", candidateID, companyID, map[string]string{"candidate_id": candidateID}, nil)
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "candidate deleted"}, nil, requestID)
}

// ListByJob handles GET /companies/{company_id}/jobs/{job_id}/candidates
func (h *CandidateHandler) ListByJob(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	jobID := chi.URLParam(r, "job_id")
	p := pagination.FromRequest(r)

	records, total, err := h.svc.ListByJob(r.Context(), companyID, jobID, p)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	items := make([]response.JobCandidateItem, 0, len(records))
	for _, jc := range records {
		items = append(items, toJobCandidateItem(jc))
	}

	meta := &pkgresponse.Meta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		Total:      total,
		TotalPages: pagination.CalcTotalPages(total, p.PageSize),
	}
	pkgresponse.JSON(w, http.StatusOK, items, meta, requestID)
}

// AssignToJob handles POST /companies/{company_id}/jobs/{job_id}/candidates/{candidate_id}/assign
func (h *CandidateHandler) AssignToJob(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	jobID := chi.URLParam(r, "job_id")
	candidateID := chi.URLParam(r, "candidate_id")

	var req request.AssignCandidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid JSON body", []string{err.Error()}), requestID)
		return
	}
	if msgs := validator.Validate(&req); msgs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("validation failed", msgs), requestID)
		return
	}

	jc, err := h.svc.AssignToJob(r.Context(), companyID, jobID, candidateID, &req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Audit log candidate assign to job
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("candidate:assign", "job_candidate", jc.ID, companyID, nil, map[string]string{
			"job_id":       jobID,
			"candidate_id": candidateID,
		})
	}

	pkgresponse.JSON(w, http.StatusCreated, toJobCandidateItem(*jc), nil, requestID)
}

// UpdatePipeline handles PUT /companies/{company_id}/jobs/{job_id}/candidates/{candidate_id}/pipeline
func (h *CandidateHandler) UpdatePipeline(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	jobID := chi.URLParam(r, "job_id")
	candidateID := chi.URLParam(r, "candidate_id")

	var body struct {
		PipelineStatus string `json:"pipeline_status" validate:"required,oneof=new screening invited interviewing completed passed rejected talent_pool"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid JSON body", []string{err.Error()}), requestID)
		return
	}
	if msgs := validator.Validate(&body); msgs != nil {
		pkgresponse.Error(w, apierrors.NewValidation("validation failed", msgs), requestID)
		return
	}

	if err := h.svc.UpdatePipelineStatus(r.Context(), companyID, jobID, candidateID, body.PipelineStatus); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Audit log pipeline status change
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("candidate:pipeline_update", "job_candidate", candidateID, companyID, nil, map[string]string{
			"job_id":          jobID,
			"pipeline_status": body.PipelineStatus,
		})
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "pipeline status updated"}, nil, requestID)
}

// UnassignFromJob handles DELETE /companies/{company_id}/jobs/{job_id}/candidates/{candidate_id}/unassign
func (h *CandidateHandler) UnassignFromJob(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	jobID := chi.URLParam(r, "job_id")
	candidateID := chi.URLParam(r, "candidate_id")

	if err := h.svc.UnassignFromJob(r.Context(), companyID, jobID, candidateID); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Audit log unassign
	if ah := middleware.GetAuditHelper(r); ah != nil {
		ah.Log("candidate:unassign", "job_candidate", candidateID, companyID, nil, map[string]string{
			"job_id": jobID,
		})
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "candidate unassigned from job"}, nil, requestID)
}

// ---------------------------------------------------------------------------
// Mapping helpers
// ---------------------------------------------------------------------------

func toCandidateListItem(c models.Candidate) response.CandidateListItem {
	item := response.CandidateListItem{
		ID:       c.ID,
		FullName: c.FullName,
		Email:    c.Email,
		Phone:    c.Phone.String,
		Status:   string(c.Status),
		Source:   c.Source.String,
	}
	if c.LatestJobID.Valid && c.LatestJobTitle.Valid {
		item.LatestJob = &response.JobBasicInfo{
			ID:    c.LatestJobID.String,
			Title: c.LatestJobTitle.String,
		}
	}
	return item
}

func toCandidateDetail(c models.Candidate) response.CandidateDetail {
	var tags []string
	if len(c.Tags) > 0 {
		_ = json.Unmarshal(c.Tags, &tags)
	}
	
	var skills []string
	var experience string
	var education string
	if len(c.ParsedCVJSON) > 0 && string(c.ParsedCVJSON) != "null" {
		var parsed struct {
			Skills     []string `json:"skills"`
			Experience string   `json:"experience"`
			Education  string   `json:"education"`
		}
		if err := json.Unmarshal(c.ParsedCVJSON, &parsed); err == nil {
			skills = parsed.Skills
			experience = parsed.Experience
			education = parsed.Education
		}
	}

	detail := response.CandidateDetail{
		ID:          c.ID,
		CompanyID:   c.CompanyID,
		FullName:    c.FullName,
		Email:       c.Email,
		Phone:       c.Phone.String,
		AvatarURL:   c.AvatarURL.String,
		Status:      string(c.Status),
		Source:      c.Source.String,
		Tags:        tags,
		AICVSummary: c.AICVSummary.String,
		Skills:      skills,
		Experience:  experience,
		Education:   education,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
	if c.CVFileID.Valid {
		detail.CVFile = &response.CVFile{
			ID:           c.CVFileID.String,
			OriginalName: c.CVOriginalName.String,
			DownloadURL:  "/api/v1/files/" + c.CVFileID.String + "/signed-url",
		}
	}
	if c.LatestJobID.Valid && c.LatestJobTitle.Valid {
		detail.LatestJob = &response.JobBasicInfo{
			ID:    c.LatestJobID.String,
			Title: c.LatestJobTitle.String,
		}
	}
	return detail
}

func toJobCandidateItem(jc models.JobCandidate) response.JobCandidateItem {
	item := response.JobCandidateItem{
		ID:             jc.ID,
		CandidateID:    jc.CandidateID,
		PipelineStatus: jc.PipelineStatus,
		FitScore:       jc.FitScore.Float64,
	}
	if jc.AppliedAt.Valid {
		t := jc.AppliedAt.Time
		item.AppliedAt = &t
	}
	return item
}

func (h *CandidateHandler) ParseCV(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	candidateID := chi.URLParam(r, "candidate_id")

	err := h.aiSvc.ParseCV(r.Context(), companyID, candidateID)
	if err != nil {
		writeServiceError(w, apierrors.NewInternal(err.Error()), requestID)
		return
	}
	
pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "CV parsed successfully"}, nil, requestID)
}

// UploadCV handles POST /companies/{company_id}/candidates/{candidate_id}/cv
func (h *CandidateHandler) UploadCV(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	companyID, _ := r.Context().Value(middleware.CtxCompanyID).(string)
	candidateID := chi.URLParam(r, "candidate_id")

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		pkgresponse.Error(w, apierrors.NewBadRequest("failed to parse form data"), requestID)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		pkgresponse.Error(w, apierrors.NewBadRequest("file is required"), requestID)
		return
	}
	defer file.Close()

	// 1. Process upload using file service (using candidate_id as userID since it's recruiter uploading for candidate)
	// We use "cv" as folder
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	fileRecord, err := h.fileSvc.ProcessUpload(r.Context(), file, header, userID, companyID, "cv")
	if err != nil {
		pkgresponse.Error(w, apierrors.NewInternal("failed to process file upload"), requestID)
		return
	}

	// 2. Update candidate with file ID if candidate_id is not "new"
	if candidateID != "new" {
		updateReq := &request.UpdateCandidateRequest{
			CVFileID: &fileRecord.ID,
		}
		_, err = h.svc.Update(r.Context(), companyID, candidateID, updateReq)
		if err != nil {
			writeServiceError(w, err, requestID)
			return
		}
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{
		"message": "CV uploaded successfully",
		"id":      fileRecord.ID,
		"url":     "/api/v1/files/" + fileRecord.ID + "/signed-url",
	}, nil, requestID)
}
