package service

import (
	"context"
	"fmt"
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

// Render replaces variables safely by wrapping user input in delimiter tags.
func (s *PromptService) Render(content string, variables map[string]string) string {
	rendered := content
	for k, v := range variables {
		// Neutralise delimiter-breaking sequences from untrusted input.
		safe := strings.ReplaceAll(v, "```", "ʼʼʼ")
		safe = strings.ReplaceAll(safe, "<<<END_USER_DATA>>>", "") // Extra safety to prevent breakout
		wrapped := fmt.Sprintf("\n<<<USER_DATA:%s>>>\n%s\n<<<END_USER_DATA>>>\n", k, safe)
		rendered = strings.ReplaceAll(rendered, "{{"+k+"}}", wrapped)
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
