package service

import (
	"context"

	"github.com/google/uuid"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/pkg/errors"
)

type CompanyTemplateService struct {
	repo *repository.CompanyTemplateRepository
}

func NewCompanyTemplateService(repo *repository.CompanyTemplateRepository) *CompanyTemplateService {
	return &CompanyTemplateService{repo: repo}
}

type CreateCompanyTemplateReq struct {
	Title       string   `json:"title" validate:"required"`
	Type        string   `json:"type" validate:"required"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func (s *CompanyTemplateService) Create(ctx context.Context, companyID string, req *CreateCompanyTemplateReq) (*models.CompanyTemplate, error) {
	t := &models.CompanyTemplate{
		ID:          uuid.NewString(),
		CompanyID:   companyID,
		Title:       req.Title,
		Type:        req.Type,
		Description: req.Description,
	}
	t.SetTags(req.Tags)

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, errors.NewInternal("failed to create template")
	}

	return t, nil
}

type UpdateCompanyTemplateReq struct {
	Title       string   `json:"title" validate:"required"`
	Type        string   `json:"type" validate:"required"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func (s *CompanyTemplateService) Update(ctx context.Context, companyID, id string, req *UpdateCompanyTemplateReq) (*models.CompanyTemplate, error) {
	t, err := s.repo.FindByID(ctx, companyID, id)
	if err != nil {
		return nil, errors.NewNotFound("template not found")
	}

	t.Title = req.Title
	t.Type = req.Type
	t.Description = req.Description
	t.SetTags(req.Tags)

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, errors.NewInternal("failed to update template")
	}

	return t, nil
}

func (s *CompanyTemplateService) Delete(ctx context.Context, companyID, id string) error {
	return s.repo.Delete(ctx, companyID, id)
}

func (s *CompanyTemplateService) List(ctx context.Context, companyID string) ([]models.CompanyTemplate, error) {
	return s.repo.FindByCompanyID(ctx, companyID)
}

func (s *CompanyTemplateService) GetByID(ctx context.Context, companyID, id string) (*models.CompanyTemplate, error) {
	t, err := s.repo.FindByID(ctx, companyID, id)
	if err != nil {
		return nil, errors.NewNotFound("template not found")
	}
	return t, nil
}
