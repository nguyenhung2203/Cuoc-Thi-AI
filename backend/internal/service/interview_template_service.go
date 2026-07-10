package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"backend/internal/dto/request"
	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/pagination"
	"backend/internal/repository"
)

type InterviewTemplateService struct {
	repo repository.InterviewTemplateRepository
}

func NewInterviewTemplateService(repo repository.InterviewTemplateRepository) *InterviewTemplateService {
	return &InterviewTemplateService{repo: repo}
}

func (s *InterviewTemplateService) List(ctx context.Context, companyID string, p pagination.Params) ([]models.InterviewTemplate, int, error) {
	return s.repo.List(ctx, companyID, p)
}

func (s *InterviewTemplateService) GetByID(ctx context.Context, companyID, templateID string) (*models.InterviewTemplate, error) {
	t, err := s.repo.GetByID(ctx, companyID, templateID)
	if err != nil {
		return nil, errors.NewInternal("failed to get template")
	}
	if t == nil {
		return nil, errors.NewNotFound("template not found")
	}
	return t, nil
}

func (s *InterviewTemplateService) Create(ctx context.Context, companyID, userID string, req *request.CreateInterviewTemplateRequest) (*models.InterviewTemplate, error) {
	if req.Name == "" {
		return nil, errors.NewValidation("name", []string{"name is required"})
	}

	configJSON := models.JSONB("{}")
	if len(req.Config) > 0 {
		configJSON = models.JSONB(req.Config)
	}

	// type and duration_minutes are NOT NULL in the schema — apply defaults.
	tType := req.Type
	if tType == "" {
		tType = "custom"
	}
	duration := req.DurationMinutes
	if duration <= 0 {
		duration = 60
	}

	t := &models.InterviewTemplate{
		ID:              uuid.NewString(),
		CompanyID:       sql.NullString{String: companyID, Valid: true},
		Name:            req.Name,
		Type:            tType,
		DurationMinutes: duration,
		Description:     sql.NullString{String: req.Description, Valid: req.Description != ""},
		ConfigJSON:      configJSON,
		CreatedBy:       sql.NullString{String: userID, Valid: true},
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, errors.NewInternal("failed to create template")
	}
	return t, nil
}

func (s *InterviewTemplateService) Update(ctx context.Context, companyID, templateID string, req *request.UpdateInterviewTemplateRequest) error {
	if _, err := s.GetByID(ctx, companyID, templateID); err != nil {
		return err
	}

	patch := make(map[string]any)
	if req.Name != nil {
		patch["name"] = *req.Name
	}
	if req.Type != nil {
		patch["type"] = *req.Type
	}
	if req.DurationMinutes != nil {
		patch["duration_minutes"] = *req.DurationMinutes
	}
	if req.Description != nil {
		patch["description"] = sql.NullString{String: *req.Description, Valid: *req.Description != ""}
	}
	if len(req.Config) > 0 {
		patch["config_json"] = models.JSONB(req.Config)
	}

	return s.repo.Update(ctx, companyID, templateID, patch)
}

func (s *InterviewTemplateService) Delete(ctx context.Context, companyID, templateID string) error {
	if _, err := s.GetByID(ctx, companyID, templateID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, companyID, templateID)
}
