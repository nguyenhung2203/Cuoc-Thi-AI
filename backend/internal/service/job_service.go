package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"backend/internal/ai"
	"backend/internal/dto/request"
	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/pagination"
	"backend/internal/pkg/utils"
	"backend/internal/repository"
	"strings"
)

// JobService implements business logic for Job CRUD and analysis.
type JobService struct {
	jobRepo      *repository.JobRepository
	jdAnalyzer   *ai.JDAnalyzer
	qGenerator   *ai.QuestionGenerator
	rubricRepo   repository.RubricRepository
	questionRepo repository.QuestionRepository
}

// NewJobService constructs a JobService with its required repositories and AI clients.
func NewJobService(
	jobRepo *repository.JobRepository,
	jdAnalyzer *ai.JDAnalyzer,
	qGenerator *ai.QuestionGenerator,
	rubricRepo repository.RubricRepository,
	questionRepo repository.QuestionRepository,
) *JobService {
	return &JobService{
		jobRepo:      jobRepo,
		jdAnalyzer:   jdAnalyzer,
		qGenerator:   qGenerator,
		rubricRepo:   rubricRepo,
		questionRepo: questionRepo,
	}
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

// ListAllOpen returns all open jobs across all companies.
func (s *JobService) ListAllOpen(ctx context.Context, keyword string, p pagination.Params) ([]models.Job, int, error) {
	result, err := s.jobRepo.ListAllOpen(ctx, keyword, p)
	if err != nil {
		return nil, 0, errors.NewInternal("failed to list all open jobs")
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
		AIAnalysisJSON: models.JSONB("null"),
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

	// Check if there are active candidates
	count, err := s.jobRepo.CountCandidates(ctx, jobID)
	if err != nil {
		return errors.NewInternal("failed to count candidates for job")
	}
	if count > 0 {
		return errors.NewConflict("Không thể xóa công việc này vì đang có ứng viên ứng tuyển.")
	}

	if err := s.jobRepo.SoftDelete(ctx, companyID, jobID); err != nil {
		return errors.NewInternal("failed to delete job")
	}
	return nil
}

// GetStats returns the candidate and interview counts for a job.
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

// Analyze triggers AI analysis for a job.
func (s *JobService) Analyze(ctx context.Context, companyID, jobID string) error {
	job, err := s.GetByID(ctx, companyID, jobID)
	if err != nil {
		return err
	}

	// 1. Prevent duplicate analysis
	if job.AIAnalysisJSON != nil && string(job.AIAnalysisJSON) != "null" {
		return errors.NewConflict("job has already been analyzed")
	}

	// 2. Truncate long descriptions to prevent Resource Exhaustion (DoS)
	safeDescription := utils.TruncateText(job.Description, 15000)

	result, err := s.jdAnalyzer.AnalyzeJD(ctx, safeDescription, job.Title, job.Level.String, job.Department.String, companyID)
	if err != nil {
		// Preserve upstream AppErrors (502 circuit breaker, 422 insufficient
		// data); everything else is an AI outage → 502, not a raw 500.
		if appErr, ok := errors.IsAppError(err); ok {
			return appErr
		}
		return errors.NewAIServiceError("AI phân tích JD thất bại: " + err.Error())
	}
	if result == nil {
		return errors.NewAIServiceError("AI phân tích JD không trả về kết quả")
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return errors.NewInternal(fmt.Sprintf("failed to marshal AI analysis result: %v", err))
	}

	// Prepare Rubric & Criteria
	var rubric *models.Rubric
	var criteria []models.RubricCriteria

	if len(result.SuggestedRubric) > 0 {
		rubric = &models.Rubric{
			ID:          uuid.NewString(),
			CompanyID:   companyID,
			JobID:       sql.NullString{String: jobID, Valid: true},
			Name:        fmt.Sprintf("Rubric for %s", job.Title),
			TotalWeight: 100,
			CreatedBy:   job.CreatedBy,
		}

		for i, c := range result.SuggestedRubric {
			criteria = append(criteria, models.RubricCriteria{
				ID:          uuid.NewString(),
				RubricID:    rubric.ID,
				Name:        c.Name,
				Description: sql.NullString{String: c.Description, Valid: true},
				Weight:      float64(c.Weight),
				MinScore:    1,
				MaxScore:    5,
				OrderIndex:  i,
			})
		}
	}

	// Prepare Questions
	var questions []models.QuestionBank
	if len(result.SuggestedQuestions) > 0 {
		for _, sq := range result.SuggestedQuestions {
			tagsBytes, err := json.Marshal([]string{sq.TargetSkill})
			if err != nil {
				return errors.NewInternal(fmt.Sprintf("failed to marshal target skill tags: %v", err))
			}
			signalsBytes, err := json.Marshal(sq.ExpectedSignals)
			if err != nil {
				return errors.NewInternal(fmt.Sprintf("failed to marshal expected signals: %v", err))
			}
			questions = append(questions, models.QuestionBank{
				CompanyID:       sql.NullString{String: companyID, Valid: true},
				JobID:           sql.NullString{String: jobID, Valid: true},
				CreatedBy:       sql.NullString{String: job.CreatedBy, Valid: true},
				QuestionText:    sq.QuestionText,
				QuestionType:    sq.QuestionType,
				SkillTags:       models.JSONB(tagsBytes),
				Level:           sql.NullString{String: sq.Difficulty, Valid: true},
				ExpectedSignals: models.JSONB(signalsBytes),
				IsAIGenerated:   true,
			})
		}
	}

	// 3. Save all results transactionally
	if err := s.jobRepo.SaveAIAnalysisWithTx(ctx, companyID, jobID, result.Summary, resultBytes, rubric, criteria, questions); err != nil {
		if strings.Contains(err.Error(), "CONFLICT") {
			return errors.NewConflict("job has already been analyzed")
		}
		return errors.NewInternal(fmt.Sprintf("failed to save AI analysis results: %v", err))
	}

	return nil
}

// GenerateQuestions uses AI to generate interview questions for a job.
func (s *JobService) GenerateQuestions(ctx context.Context, companyID, jobID string, req *request.GenerateQuestionsRequest) ([]models.QuestionBank, error) {
	job, err := s.GetByID(ctx, companyID, jobID)
	if err != nil {
		return nil, err
	}

	count := req.Count
	if count <= 0 {
		count = 10
	}
	questionTypes := "behavioral,technical"
	if len(req.QuestionTypes) > 0 {
		questionTypes = strings.Join(req.QuestionTypes, ",")
	}

	safeDescription := utils.TruncateText(job.Description, 15000)

	result, err := s.qGenerator.GenerateQuestions(ctx, safeDescription, "", "", req.Level, fmt.Sprintf("%d", count), questionTypes, companyID)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			return nil, appErr
		}
		return nil, errors.NewAIServiceError("AI sinh câu hỏi thất bại: " + err.Error())
	}
	if result == nil {
		return nil, errors.NewAIServiceError("AI sinh câu hỏi không trả về kết quả")
	}

	var questions []models.QuestionBank
	for _, sq := range result.Questions {
		tagsBytes, _ := json.Marshal([]string{sq.TargetSkill})
		signalsBytes, _ := json.Marshal(sq.ExpectedSignals)
		questions = append(questions, models.QuestionBank{
			ID:              uuid.NewString(),
			CompanyID:       sql.NullString{String: companyID, Valid: true},
			JobID:           sql.NullString{String: jobID, Valid: true},
			CreatedBy:       sql.NullString{String: job.CreatedBy, Valid: true},
			QuestionText:    sq.QuestionText,
			QuestionType:    sq.QuestionType,
			SkillTags:       models.JSONB(tagsBytes),
			Level:           sql.NullString{String: sq.Difficulty, Valid: true},
			ExpectedSignals: models.JSONB(signalsBytes),
			IsAIGenerated:   true,
		})
	}

	if err := s.questionRepo.CreateQuestions(ctx, questions); err != nil {
		return nil, errors.NewInternal("failed to save generated questions")
	}

	return questions, nil
}
