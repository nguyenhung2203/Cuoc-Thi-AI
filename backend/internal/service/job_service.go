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

// JobService implements business logic for Job CRUD and analysis.
type JobService struct {
	jobRepo *repository.JobRepository
}

// NewJobService constructs a JobService with its required repository.
func NewJobService(jobRepo *repository.JobRepository) *JobService {
	return &JobService{jobRepo: jobRepo}
}

// List returns a paginated list of jobs scoped to companyID.
func (s *JobService) List(
	ctx context.Context,
	companyID, status, keyword string,
	p pagination.Params,
) ([]models.Job, int, error) {
	result, err := s.jobRepo.List(ctx, companyID, status, keyword, p)
	if err != nil {
		return nil, 0, errors.NewInternal("failed to list jobs")
	}
	return result.Jobs, result.Total, nil
}

// GetByID returns a single job or a NOT_FOUND error.
func (s *JobService) GetByID(ctx context.Context, companyID, jobID string) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(ctx, companyID, jobID)
	if err != nil {
		return nil, errors.NewInternal("failed to get job")
	}
	if job == nil {
		return nil, errors.NewNotFound("job not found")
	}
	return job, nil
}

// Create validates and persists a new Job record.
func (s *JobService) Create(
	ctx context.Context,
	companyID, createdByUserID string,
	req *request.CreateJobRequest,
) (*models.Job, error) {
	now := time.Now()
	job := &models.Job{
		ID:             uuid.NewString(),
		CompanyID:      companyID,
		Title:          req.Title,
		Department:     sql.NullString{String: req.Department, Valid: req.Department != ""},
		Level:          sql.NullString{String: req.Level, Valid: req.Level != ""},
		Location:       sql.NullString{String: req.Location, Valid: req.Location != ""},
		EmploymentType: sql.NullString{String: req.EmploymentType, Valid: req.EmploymentType != ""},
		SalaryMin:      sql.NullFloat64{Float64: req.SalaryMin, Valid: req.SalaryMin != 0},
		SalaryMax:      sql.NullFloat64{Float64: req.SalaryMax, Valid: req.SalaryMax != 0},
		Currency:       sql.NullString{String: req.Currency, Valid: req.Currency != ""},
		Description:    req.Description,
		Requirements:   sql.NullString{String: req.Requirements, Valid: req.Requirements != ""},
		Benefits:       sql.NullString{String: req.Benefits, Valid: req.Benefits != ""},
		Status:         models.JobStatus(req.Status),
		AIAnalysisJSON: json.RawMessage("null"),
		CreatedBy:      createdByUserID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	created, err := s.jobRepo.Create(ctx, job)
	if err != nil {
		return nil, errors.NewInternal("failed to create job")
	}
	return created, nil
}

// Update applies a partial patch to an existing job.
func (s *JobService) Update(
	ctx context.Context,
	companyID, jobID string,
	req *request.UpdateJobRequest,
) (*models.Job, error) {
	if _, err := s.GetByID(ctx, companyID, jobID); err != nil {
		return nil, err
	}

	patch := make(map[string]any)

	if req.Title != nil {
		patch["title"] = *req.Title
	}
	if req.Department != nil {
		patch["department"] = sql.NullString{String: *req.Department, Valid: *req.Department != ""}
	}
	if req.Level != nil {
		patch["level"] = sql.NullString{String: *req.Level, Valid: *req.Level != ""}
	}
	if req.Location != nil {
		patch["location"] = sql.NullString{String: *req.Location, Valid: *req.Location != ""}
	}
	if req.EmploymentType != nil {
		patch["employment_type"] = sql.NullString{String: *req.EmploymentType, Valid: *req.EmploymentType != ""}
	}
	if req.SalaryMin != nil {
		patch["salary_min"] = sql.NullFloat64{Float64: *req.SalaryMin, Valid: *req.SalaryMin != 0}
	}
	if req.SalaryMax != nil {
		patch["salary_max"] = sql.NullFloat64{Float64: *req.SalaryMax, Valid: *req.SalaryMax != 0}
	}
	if req.Currency != nil {
		patch["currency"] = sql.NullString{String: *req.Currency, Valid: *req.Currency != ""}
	}
	if req.Description != nil {
		patch["description"] = *req.Description
	}
	if req.Requirements != nil {
		patch["requirements"] = sql.NullString{String: *req.Requirements, Valid: *req.Requirements != ""}
	}
	if req.Benefits != nil {
		patch["benefits"] = sql.NullString{String: *req.Benefits, Valid: *req.Benefits != ""}
	}
	if req.Status != nil {
		patch["status"] = models.JobStatus(*req.Status)
	}

	if len(patch) == 0 {
		return s.GetByID(ctx, companyID, jobID)
	}

	updated, err := s.jobRepo.Update(ctx, companyID, jobID, patch)
	if err != nil {
		return nil, errors.NewInternal("failed to update job")
	}
	return updated, nil
}

// Delete soft-deletes a job after confirming it exists.
func (s *JobService) Delete(ctx context.Context, companyID, jobID string) error {
	if _, err := s.GetByID(ctx, companyID, jobID); err != nil {
		return err
	}
	if err := s.jobRepo.SoftDelete(ctx, companyID, jobID); err != nil {
		return errors.NewInternal("failed to delete job")
	}
	return nil
}

// GetStats returns the candidate and interview counts for a job.
// Counts are fetched sequentially; the first error short-circuits.
func (s *JobService) GetStats(ctx context.Context, jobID string) (candidateCount, interviewCount int, err error) {
	candidateCount, err = s.jobRepo.CountCandidates(ctx, jobID)
	if err != nil {
		return 0, 0, errors.NewInternal("failed to count candidates")
	}
	interviewCount, err = s.jobRepo.CountInterviews(ctx, jobID)
	if err != nil {
		return 0, 0, errors.NewInternal("failed to count interviews")
	}
	return candidateCount, interviewCount, nil
}
