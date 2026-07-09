package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/repository"
)

type CompanyService struct {
	companyRepo *repository.CompanyRepository
}

func NewCompanyService(companyRepo *repository.CompanyRepository) *CompanyService {
	return &CompanyService{companyRepo: companyRepo}
}

func (s *CompanyService) CreateCompany(ctx context.Context, userID, name, website, industry, size string) (*models.Company, error) {
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-")) + "-" + uuid.NewString()[:8]
	
	company := &models.Company{
		ID:        uuid.NewString(),
		Name:      name,
		Slug:      slug,
		Website:   website,
		Industry:  industry,
		Size:      size,
		CreatedBy: userID,
	}

	if err := s.companyRepo.Create(ctx, company, userID); err != nil {
		return nil, errors.NewInternal("failed to create company")
	}

	return company, nil
}

func (s *CompanyService) GetCompany(ctx context.Context, id string) (*models.Company, error) {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFound("company not found")
	}
	return company, nil
}

func (s *CompanyService) ListCompanies(ctx context.Context, userID string) ([]models.Company, error) {
	companies, err := s.companyRepo.ListByUserID(ctx, userID)
	if err != nil {
		println("ListByUserID error:", err.Error())
		return nil, errors.NewInternal("failed to list companies")
	}
	return companies, nil
}

func (s *CompanyService) UpdateCompany(ctx context.Context, id, name, website, industry, size string) (*models.Company, error) {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFound("company not found")
	}
	
	company.Name = name
	company.Website = website
	company.Industry = industry
	company.Size = size

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return nil, errors.NewInternal("failed to update company")
	}

	return company, nil
}
