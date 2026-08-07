package service

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"

	"backend/internal/dto/request"
	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/pagination"
	"backend/internal/repository"
)

type QuestionBankService struct {
	repo repository.QuestionRepository
}

func NewQuestionBankService(repo repository.QuestionRepository) *QuestionBankService {
	return &QuestionBankService{repo: repo}
}

func (s *QuestionBankService) List(ctx context.Context, companyID string, jobID, questionType, level, keyword string, p pagination.Params) ([]models.QuestionBank, int, error) {
	return s.repo.List(ctx, companyID, jobID, questionType, level, keyword, p)
}

func (s *QuestionBankService) GetByID(ctx context.Context, companyID, questionID string) (*models.QuestionBank, error) {
	qb, err := s.repo.GetByID(ctx, companyID, questionID)
	if err != nil {
		return nil, errors.NewInternal("failed to get question")
	}
	if qb == nil {
		return nil, errors.NewNotFound("question not found")
	}
	return qb, nil
}

func (s *QuestionBankService) Create(ctx context.Context, companyID, userID string, req *request.CreateQuestionBankRequest) (*models.QuestionBank, error) {
	skillTags := models.JSONB("[]")
	if len(req.SkillTags) > 0 {
		b, _ := json.Marshal(req.SkillTags)
		skillTags = models.JSONB(b)
	}
	signals := models.JSONB("[]")
	if len(req.ExpectedSignals) > 0 {
		b, _ := json.Marshal(req.ExpectedSignals)
		signals = models.JSONB(b)
	}

	q := &models.QuestionBank{
		ID:              uuid.NewString(),
		CompanyID:       sql.NullString{String: companyID, Valid: true},
		CreatedBy:       sql.NullString{String: userID, Valid: true},
		QuestionText:    req.QuestionText,
		QuestionType:    req.QuestionType,
		SkillTags:       skillTags,
		Level:           sql.NullString{String: req.Level, Valid: req.Level != ""},
		ExpectedSignals: signals,
		IsAIGenerated:   false,
	}
	if req.JobID != nil && *req.JobID != "" {
		q.JobID = sql.NullString{String: *req.JobID, Valid: true}
	}

	if err := s.repo.Create(ctx, q); err != nil {
		return nil, errors.NewInternal("failed to create question")
	}
	return q, nil
}

func (s *QuestionBankService) Update(ctx context.Context, companyID, questionID string, req *request.UpdateQuestionBankRequest) error {
	if _, err := s.GetByID(ctx, companyID, questionID); err != nil {
		return err
	}

	patch := make(map[string]any)
	if req.QuestionText != nil {
		patch["question_text"] = *req.QuestionText
	}
	if req.QuestionType != nil {
		patch["question_type"] = *req.QuestionType
	}
	if req.Level != nil {
		patch["level"] = sql.NullString{String: *req.Level, Valid: *req.Level != ""}
	}
	if req.SkillTags != nil {
		b, _ := json.Marshal(*req.SkillTags)
		patch["skill_tags"] = models.JSONB(b)
	}
	if req.ExpectedSignals != nil {
		b, _ := json.Marshal(*req.ExpectedSignals)
		patch["expected_signals"] = models.JSONB(b)
	}

	return s.repo.Update(ctx, companyID, questionID, patch)
}

func (s *QuestionBankService) Delete(ctx context.Context, companyID, questionID string) error {
	if _, err := s.GetByID(ctx, companyID, questionID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, companyID, questionID)
}
