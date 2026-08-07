package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"backend/internal/dto/request"
	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/pagination"
	"backend/internal/repository"
)

type QuestionService struct {
	repo repository.QuestionRepository
}

func NewQuestionService(repo repository.QuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) List(ctx context.Context, companyID, jobID, qType, level, keyword string, p *pagination.Params) ([]models.QuestionBank, int64, error) {
	if companyID == "" {
		return nil, 0, errors.NewValidation("company_id", []string{"is required"})
	}
	items, total, err := s.repo.List(ctx, companyID, jobID, qType, level, keyword, *p)
	if err != nil {
		return nil, 0, errors.NewInternal("failed to list questions")
	}
	return items, total, nil
}

func (s *QuestionService) Create(ctx context.Context, companyID, userID string, req *request.CreateQuestionRequest) (*models.QuestionBank, error) {
	q := &models.QuestionBank{
		ID:            uuid.NewString(),
		CompanyID:     sql.NullString{String: companyID, Valid: true},
		CreatedBy:     sql.NullString{String: userID, Valid: true},
		QuestionText:  req.QuestionText,
		QuestionType:  req.QuestionType,
		SkillTags:     req.SkillTags,
		ExpectedSignals: req.ExpectedSignals,
		IsAIGenerated: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if req.JobID != "" {
		q.JobID = sql.NullString{String: req.JobID, Valid: true}
	}
	if req.Level != "" {
		q.Level = sql.NullString{String: req.Level, Valid: true}
	}

	if err := s.repo.Create(ctx, q); err != nil {
		return nil, errors.NewInternal("failed to create question")
	}
	return q, nil
}

func (s *QuestionService) Update(ctx context.Context, companyID, id string, req *request.UpdateQuestionRequest) error {
	updates := make(map[string]interface{})
	if req.QuestionText != "" {
		updates["question_text"] = req.QuestionText
	}
	if req.QuestionType != "" {
		updates["question_type"] = req.QuestionType
	}
	if req.Level != "" {
		updates["level"] = req.Level
	}
	if len(req.SkillTags) > 0 {
		updates["skill_tags"] = req.SkillTags
	}
	if len(req.ExpectedSignals) > 0 {
		updates["expected_signals"] = req.ExpectedSignals
	}

	err := s.repo.Update(ctx, companyID, id, updates)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewNotFound("question not found")
		}
		return errors.NewInternal("failed to update question")
	}
	return nil
}

func (s *QuestionService) Delete(ctx context.Context, companyID, id string) error {
	err := s.repo.Delete(ctx, companyID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewNotFound("question not found")
		}
		return errors.NewInternal("failed to delete question")
	}
	return nil
}
