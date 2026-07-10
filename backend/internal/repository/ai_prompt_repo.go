package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
)

type AIPromptRepository struct {
	db *sqlx.DB
}

func NewAIPromptRepository(db *sqlx.DB) *AIPromptRepository {
	return &AIPromptRepository{db: db}
}

// GetLatestActive returns the latest active version of a prompt template by name.
// If companyID is provided, it tries to fetch the company-specific template first.
// If not found or companyID is empty, it falls back to the system-wide template (company_id IS NULL).
func (r *AIPromptRepository) GetLatestActive(ctx context.Context, name, companyID string) (*models.AIPromptTemplate, error) {
	// First try to get company specific template if companyID is not empty
	if companyID != "" {
		const qCompany = `
			SELECT * FROM ai_prompt_templates
			WHERE name = $1 AND company_id = $2::uuid AND is_active = true AND deleted_at IS NULL
			ORDER BY version DESC LIMIT 1`

		var tmpl models.AIPromptTemplate
		err := r.db.GetContext(ctx, &tmpl, qCompany, name, companyID)
		if err == nil {
			return &tmpl, nil
		}
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("get company prompt: %w", err)
		}
	}

	// Fallback to system-wide template
	const qSystem = `
		SELECT * FROM ai_prompt_templates
		WHERE name = $1 AND company_id IS NULL AND is_active = true AND deleted_at IS NULL
		ORDER BY version DESC LIMIT 1`

	var sysTmpl models.AIPromptTemplate
	err := r.db.GetContext(ctx, &sysTmpl, qSystem, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("get system prompt: %w", err)
	}

	return &sysTmpl, nil
}

// GetByVersion returns a specific version of a prompt template.
func (r *AIPromptRepository) GetByVersion(ctx context.Context, name, companyID string, version int) (*models.AIPromptTemplate, error) {
	if companyID != "" {
		const qCompany = `
			SELECT * FROM ai_prompt_templates
			WHERE name = $1 AND company_id = $2::uuid AND version = $3 AND deleted_at IS NULL`

		var tmpl models.AIPromptTemplate
		err := r.db.GetContext(ctx, &tmpl, qCompany, name, companyID, version)
		if err == nil {
			return &tmpl, nil
		}
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("get company prompt by version: %w", err)
		}
	}

	const qSystem = `
		SELECT * FROM ai_prompt_templates
		WHERE name = $1 AND company_id IS NULL AND version = $2 AND deleted_at IS NULL`

	var sysTmpl models.AIPromptTemplate
	err := r.db.GetContext(ctx, &sysTmpl, qSystem, name, version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("get system prompt by version: %w", err)
	}

	return &sysTmpl, nil
}

// CreateNewVersion inserts a new version of a prompt template. 
// It automatically calculates the next version number.
func (r *AIPromptRepository) CreateNewVersion(ctx context.Context, t *models.AIPromptTemplate) (*models.AIPromptTemplate, error) {
	// Determine next version
	var nextVersion int
	var qMaxVersion string
	var args []any
	
	if t.CompanyID.Valid {
		qMaxVersion = `SELECT COALESCE(MAX(version), 0) FROM ai_prompt_templates WHERE name = $1 AND company_id = $2::uuid`
		args = []any{t.Name, t.CompanyID.String}
	} else {
		qMaxVersion = `SELECT COALESCE(MAX(version), 0) FROM ai_prompt_templates WHERE name = $1 AND company_id IS NULL`
		args = []any{t.Name}
	}

	if err := r.db.QueryRowContext(ctx, qMaxVersion, args...).Scan(&nextVersion); err != nil {
		return nil, fmt.Errorf("get max version: %w", err)
	}
	nextVersion++

	const qInsert = `
		INSERT INTO ai_prompt_templates (
			company_id, name, version, content, variables_schema, 
			model, params, is_active, created_by
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9
		) RETURNING *`

	var inserted models.AIPromptTemplate
	err := r.db.QueryRowxContext(ctx, qInsert,
		t.CompanyID, t.Name, nextVersion, t.Content, t.VariablesSchema,
		t.Model, t.Params, t.IsActive, t.CreatedBy,
	).StructScan(&inserted)
	
	if err != nil {
		return nil, fmt.Errorf("insert prompt template: %w", err)
	}

	return &inserted, nil
}

func (r *AIPromptRepository) ListAllTemplates(ctx context.Context) ([]models.AIPromptTemplate, error) {
	const qList = `
		SELECT * FROM ai_prompt_templates
		WHERE deleted_at IS NULL
		ORDER BY name ASC, version DESC`

	var tmpls []models.AIPromptTemplate
	err := r.db.SelectContext(ctx, &tmpls, qList)
	if err != nil {
		return nil, fmt.Errorf("list prompt templates: %w", err)
	}
	return tmpls, nil
}
