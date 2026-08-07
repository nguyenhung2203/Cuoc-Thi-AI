package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"backend/internal/models"
)

type CompanyTemplateRepository struct {
	db *sqlx.DB
}

func NewCompanyTemplateRepository(db *sqlx.DB) *CompanyTemplateRepository {
	return &CompanyTemplateRepository{db: db}
}

func (r *CompanyTemplateRepository) Create(ctx context.Context, template *models.CompanyTemplate) error {
	query := `
		INSERT INTO company_templates (id, company_id, title, type, description, tags)
		VALUES (:id, :company_id, :title, :type, :description, :tags)
	`
	_, err := r.db.NamedExecContext(ctx, query, template)
	return err
}

func (r *CompanyTemplateRepository) Update(ctx context.Context, template *models.CompanyTemplate) error {
	query := `
		UPDATE company_templates 
		SET title = :title, type = :type, description = :description, tags = :tags, updated_at = NOW()
		WHERE id = :id AND company_id = :company_id
	`
	_, err := r.db.NamedExecContext(ctx, query, template)
	return err
}

func (r *CompanyTemplateRepository) Delete(ctx context.Context, companyID, id string) error {
	query := `DELETE FROM company_templates WHERE id = $1 AND company_id = $2`
	_, err := r.db.ExecContext(ctx, query, id, companyID)
	return err
}

func (r *CompanyTemplateRepository) FindByCompanyID(ctx context.Context, companyID string) ([]models.CompanyTemplate, error) {
	var templates []models.CompanyTemplate
	query := `SELECT * FROM company_templates WHERE company_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &templates, query, companyID)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *CompanyTemplateRepository) FindByID(ctx context.Context, companyID, id string) (*models.CompanyTemplate, error) {
	var template models.CompanyTemplate
	query := `SELECT * FROM company_templates WHERE id = $1 AND company_id = $2`
	err := r.db.GetContext(ctx, &template, query, id, companyID)
	if err != nil {
		return nil, err
	}
	return &template, nil
}
