package service

import (
	"context"
	"strings"

	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/repository"
)

type PromptService struct {
	repo *repository.AIPromptRepository
}

func NewPromptService(repo *repository.AIPromptRepository) *PromptService {
	return &PromptService{repo: repo}
}

// LoadTemplate retrieves the prompt template by name. If version is 0, it gets the latest active.
func (s *PromptService) LoadTemplate(ctx context.Context, name, companyID string, version int) (*models.AIPromptTemplate, error) {
	var tmpl *models.AIPromptTemplate
	var err error

	if version == 0 {
		tmpl, err = s.repo.GetLatestActive(ctx, name, companyID)
	} else {
		tmpl, err = s.repo.GetByVersion(ctx, name, companyID, version)
	}

	if err != nil {
		return nil, errors.NewInternal("failed to load prompt template")
	}
	if tmpl == nil {
		return nil, errors.NewNotFound("prompt template not found")
	}

	return tmpl, nil
}

// Render replaces variables in the format {{key}} with their corresponding values.
// This is a basic implementation. For complex rendering, text/template could be used.
func (s *PromptService) Render(content string, variables map[string]string) string {
	rendered := content
	for k, v := range variables {
		// Replace {{key}} with the value
		rendered = strings.ReplaceAll(rendered, "{{"+k+"}}", v)
	}
	return rendered
}

// CreateNewVersion saves a new version of the template.
func (s *PromptService) CreateNewVersion(ctx context.Context, t *models.AIPromptTemplate) (*models.AIPromptTemplate, error) {
	created, err := s.repo.CreateNewVersion(ctx, t)
	if err != nil {
		return nil, errors.NewInternal("failed to create prompt template version")
	}
	return created, nil
}
