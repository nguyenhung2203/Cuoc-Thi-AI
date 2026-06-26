package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"backend/internal/dto/request"
	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/pagination"
	"backend/internal/repository"
)

// CandidateService implements business logic for Candidate CRUD and pipeline management.
type CandidateService struct {
	candidateRepo *repository.CandidateRepository
	jobRepo       *repository.JobRepository
}

// NewCandidateService constructs a CandidateService with its required repositories.
func NewCandidateService(
	candidateRepo *repository.CandidateRepository,
	jobRepo *repository.JobRepository,
) *CandidateService {
	return &CandidateService{
		candidateRepo: candidateRepo,
		jobRepo:       jobRepo,
	}
}

// List returns a paginated list of candidates scoped to companyID.
func (s *CandidateService) List(
	ctx context.Context,
	companyID, jobID, status, keyword string,
	p pagination.Params,
) ([]models.Candidate, int, error) {
	candidates, total, err := s.candidateRepo.List(ctx, companyID, jobID, status, keyword, p)
	if err != nil {
		return nil, 0, errors.NewInternal("failed to list candidates")
	}
	return candidates, total, nil
}

// GetByID returns a single candidate or a NOT_FOUND error.
func (s *CandidateService) GetByID(ctx context.Context, companyID, candidateID string) (*models.Candidate, error) {
	candidate, err := s.candidateRepo.GetByID(ctx, companyID, candidateID)
	if err != nil {
		return nil, errors.NewInternal("failed to get candidate")
	}
	if candidate == nil {
		return nil, errors.NewNotFound("candidate not found")
	}
	return candidate, nil
}

// Create validates uniqueness, persists a new Candidate, and optionally links to a job.
func (s *CandidateService) Create(
	ctx context.Context,
	companyID, createdByUserID string,
	req *request.CreateCandidateRequest,
) (*models.Candidate, error) {
	exists, err := s.candidateRepo.EmailExists(ctx, companyID, req.Email)
	if err != nil {
		return nil, errors.NewInternal("failed to check email uniqueness")
	}
	if exists {
		return nil, errors.NewConflict("a candidate with this email already exists")
	}

	// If a JobID is supplied, verify the job exists before creating anything.
	if req.JobID != "" {
		job, err := s.jobRepo.GetByID(ctx, companyID, req.JobID)
		if err != nil {
			return nil, errors.NewInternal("failed to verify job")
		}
		if job == nil {
			return nil, errors.NewNotFound("job not found")
		}
	}

	now := time.Now()
	candidate := &models.Candidate{
		ID:           uuid.NewString(),
		CompanyID:    companyID,
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Source:       sql.NullString{String: req.Source, Valid: req.Source != ""},
		Status:       models.CandidateStatus("new"),
		ParsedCVJSON: models.JSONB("null"),
		Tags:         models.JSONB("[]"),
		CreatedBy:    sql.NullString{String: createdByUserID, Valid: createdByUserID != ""},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := s.candidateRepo.Create(ctx, candidate)
	if err != nil {
		return nil, errors.NewInternal("failed to create candidate")
	}

	if req.JobID != "" {
		jc := &models.JobCandidate{
			ID:             uuid.NewString(),
			CompanyID:      companyID,
			JobID:          req.JobID,
			CandidateID:    created.ID,
			PipelineStatus: "new",
			AIMatchJSON:    models.JSONB("null"),
			CreatedBy:      sql.NullString{String: createdByUserID, Valid: createdByUserID != ""},
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if _, err := s.candidateRepo.CreateJobCandidate(ctx, jc); err != nil {
			return nil, errors.NewInternal("failed to link candidate to job")
		}
	}

	return created, nil
}

// Update applies a partial patch to an existing candidate.
func (s *CandidateService) Update(
	ctx context.Context,
	companyID, candidateID string,
	req *request.UpdateCandidateRequest,
) (*models.Candidate, error) {
	if _, err := s.GetByID(ctx, companyID, candidateID); err != nil {
		return nil, err
	}

	patch := make(map[string]any)

	if req.FullName != nil {
		patch["full_name"] = *req.FullName
	}
	if req.Email != nil {
		patch["email"] = *req.Email
	}
	if req.Phone != nil {
		patch["phone"] = sql.NullString{String: *req.Phone, Valid: *req.Phone != ""}
	}
	if req.Source != nil {
		patch["source"] = sql.NullString{String: *req.Source, Valid: *req.Source != ""}
	}
	if req.Status != nil {
		patch["status"] = models.CandidateStatus(*req.Status)
	}
	if req.Tags != nil {
		tagsJSON, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, errors.NewInternal("failed to marshal tags")
		}
		patch["tags"] = models.JSONB(tagsJSON)
	}

	if len(patch) == 0 {
		return s.GetByID(ctx, companyID, candidateID)
	}

	updated, err := s.candidateRepo.Update(ctx, companyID, candidateID, patch)
	if err != nil {
		return nil, errors.NewInternal("failed to update candidate")
	}
	return updated, nil
}

// Delete soft-deletes a candidate after confirming it exists.
func (s *CandidateService) Delete(ctx context.Context, companyID, candidateID string) error {
	if _, err := s.GetByID(ctx, companyID, candidateID); err != nil {
		return err
	}
	if err := s.candidateRepo.SoftDelete(ctx, companyID, candidateID); err != nil {
		return errors.NewInternal("failed to delete candidate")
	}
	return nil
}

// AssignToJob links a candidate to a job, enforcing existence and uniqueness.
func (s *CandidateService) AssignToJob(
	ctx context.Context,
	companyID, jobID, candidateID string,
	req *request.AssignCandidateRequest,
) (*models.JobCandidate, error) {
	if _, err := s.GetByID(ctx, companyID, candidateID); err != nil {
		return nil, err
	}

	job, err := s.jobRepo.GetByID(ctx, companyID, jobID)
	if err != nil {
		return nil, errors.NewInternal("failed to verify job")
	}
	if job == nil {
		return nil, errors.NewNotFound("job not found")
	}

	existing, err := s.candidateRepo.GetJobCandidate(ctx, jobID, candidateID)
	if err != nil {
		return nil, errors.NewInternal("failed to check existing assignment")
	}
	if existing != nil {
		return nil, errors.NewConflict("candidate is already assigned to this job")
	}

	now := time.Now()
	jc := &models.JobCandidate{
		ID:             uuid.NewString(),
		CompanyID:      companyID,
		JobID:          jobID,
		CandidateID:    candidateID,
		PipelineStatus: req.PipelineStatus,
		AIMatchJSON:    models.JSONB("null"),
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	created, err := s.candidateRepo.CreateJobCandidate(ctx, jc)
	if err != nil {
		return nil, errors.NewInternal("failed to assign candidate to job")
	}
	return created, nil
}

// UpdatePipelineStatus changes a candidate's pipeline stage for a given job.
func (s *CandidateService) UpdatePipelineStatus(
	ctx context.Context,
	companyID, jobID, candidateID, status string,
) error {
	if err := s.candidateRepo.UpdateJobCandidateStatus(ctx, companyID, jobID, candidateID, status); err != nil {
		return errors.NewInternal("failed to update pipeline status")
	}
	return nil
}

// ListByJob returns a paginated list of job-candidate records for a specific job.
func (s *CandidateService) ListByJob(
	ctx context.Context,
	companyID, jobID string,
	p pagination.Params,
) ([]models.JobCandidate, int, error) {
	records, total, err := s.candidateRepo.ListByJobID(ctx, companyID, jobID, p)
	if err != nil {
		return nil, 0, errors.NewInternal("failed to list candidates by job")
	}
	return records, total, nil
}

func (s *CandidateService) UpdateCVParseResult(ctx context.Context, companyID, candidateID, parsedJSON, summary string) error {
	patch := map[string]any{
		"parsed_cv_json": models.JSONB(parsedJSON),
		"ai_cv_summary":  sql.NullString{String: summary, Valid: true},
	}
	_, err := s.candidateRepo.Update(ctx, companyID, candidateID, patch)
	return err
}
