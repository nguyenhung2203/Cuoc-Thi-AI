package service

import (
	"context"
	"fmt"

	"backend/internal/models"
	"backend/internal/repository"
	apierrors "backend/internal/pkg/errors"
)

type RubricService struct {
	repo repository.RubricRepository
}

func NewRubricService(repo repository.RubricRepository) *RubricService {
	return &RubricService{repo: repo}
}

func (s *RubricService) GetRubricByID(ctx context.Context, companyID, rubricID string) (*models.Rubric, []models.RubricCriteria, error) {
	rubric, criteria, err := s.repo.GetByID(ctx, companyID, rubricID)
	if err != nil {
		return nil, nil, err
	}
	if rubric == nil {
		return nil, nil, apierrors.NewNotFound("rubric")
	}
	return rubric, criteria, nil
}

func (s *RubricService) ListCompanyRubrics(ctx context.Context, companyID, jobID string) ([]models.Rubric, error) {
	return s.repo.ListByCompany(ctx, companyID, jobID)
}

func (s *RubricService) DeleteRubric(ctx context.Context, companyID, rubricID string) error {
	// Let repo handle it.
	return s.repo.Delete(ctx, companyID, rubricID)
}

func (s *RubricService) CreateRubric(ctx context.Context, rubric *models.Rubric, criteria []models.RubricCriteria) error {
	// Simple validation
	if rubric.Name == "" {
		return apierrors.NewValidation("rubric name", []string{"name is required"})
	}
	if len(criteria) == 0 {
		return apierrors.NewValidation("criteria", []string{"at least one criterion is required"})
	}

	var totalWeight float64
	for _, c := range criteria {
		totalWeight += c.Weight
		if c.MinScore >= c.MaxScore {
			return apierrors.NewValidation("score bounds", []string{fmt.Sprintf("min_score must be less than max_score for criterion %s", c.Name)})
		}
	}
	rubric.TotalWeight = totalWeight

	return s.repo.CreateRubricWithCriteria(ctx, rubric, criteria)
}
