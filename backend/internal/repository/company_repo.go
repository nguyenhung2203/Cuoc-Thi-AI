package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"backend/internal/models"
)

type CompanyRepository struct {
	db *sqlx.DB
}

func NewCompanyRepository(db *sqlx.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(ctx context.Context, company *models.Company, ownerID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO companies (id, name, slug, website, industry, size, created_by)
		VALUES (:id, :name, :slug, :website, :industry, :size, :created_by)
	`
	if _, err := tx.NamedExecContext(ctx, query, company); err != nil {
		tx.Rollback()
		return err
	}

	queryMember := `
		INSERT INTO company_members (company_id, user_id, role, status)
		VALUES ($1, $2, 'owner', 'active')
	`
	if _, err := tx.ExecContext(ctx, queryMember, company.ID, ownerID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *CompanyRepository) FindByID(ctx context.Context, id string) (*models.Company, error) {
	var c models.Company
	err := r.db.GetContext(ctx, &c, "SELECT * FROM companies WHERE id = $1 AND deleted_at IS NULL", id)
	return &c, err
}

func (r *CompanyRepository) ListByUserID(ctx context.Context, userID string) ([]models.Company, error) {
	query := `
		SELECT c.*
		FROM companies c
		JOIN company_members cm ON c.id = cm.company_id
		WHERE cm.user_id = $1 AND cm.status = 'active' AND c.deleted_at IS NULL
	`
	var companies []models.Company
	err := r.db.SelectContext(ctx, &companies, query, userID)
	return companies, err
}

func (r *CompanyRepository) ListAll(ctx context.Context) ([]models.Company, error) {
	var companies []models.Company
	err := r.db.SelectContext(ctx, &companies, "SELECT * FROM companies WHERE deleted_at IS NULL ORDER BY created_at DESC")
	return companies, err
}

func (r *CompanyRepository) Update(ctx context.Context, company *models.Company) error {
	query := `
		UPDATE companies
		SET name = :name, website = :website, industry = :industry, size = :size, updated_at = NOW()
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, company)
	return err
}
