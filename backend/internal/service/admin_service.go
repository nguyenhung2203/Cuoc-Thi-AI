package service

import (
	"context"

	"backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type AdminService struct {
	userRepo    *repository.UserRepository
	companyRepo *repository.CompanyRepository
	auditRepo   *repository.AuditRepository
	db          *sqlx.DB
}

func NewAdminService(userRepo *repository.UserRepository, companyRepo *repository.CompanyRepository, auditRepo *repository.AuditRepository, db *sqlx.DB) *AdminService {
	return &AdminService{
		userRepo:    userRepo,
		companyRepo: companyRepo,
		auditRepo:   auditRepo,
		db:          db,
	}
}

func (s *AdminService) GetStats(ctx context.Context) (map[string]int, error) {
	var totalUsers, totalCompanies, totalJobs, totalInterviews int

	if s.db != nil {
		_ = s.db.GetContext(ctx, &totalUsers, "SELECT count(*) FROM users")
		_ = s.db.GetContext(ctx, &totalCompanies, "SELECT count(*) FROM companies")
		_ = s.db.GetContext(ctx, &totalJobs, "SELECT count(*) FROM jobs")
		_ = s.db.GetContext(ctx, &totalInterviews, "SELECT count(*) FROM interviews")
	}

	return map[string]int{
		"total_users":      totalUsers,
		"total_companies":  totalCompanies,
		"total_jobs":       totalJobs,
		"total_interviews": totalInterviews,
	}, nil
}
