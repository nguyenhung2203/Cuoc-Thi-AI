package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"backend/internal/models"
	"backend/internal/pkg/pagination"
)

type InterviewTemplateRepository interface {
	List(ctx context.Context, companyID string, p pagination.Params) ([]models.InterviewTemplate, int, error)
	GetByID(ctx context.Context, companyID, templateID string) (*models.InterviewTemplate, error)
	Create(ctx context.Context, t *models.InterviewTemplate) error
	Update(ctx context.Context, companyID, templateID string, patch map[string]any) error
	Delete(ctx context.Context, companyID, templateID string) error
}

type interviewTemplateRepository struct {
	db *sql.DB
}

func NewInterviewTemplateRepository(db *sql.DB) InterviewTemplateRepository {
	return &interviewTemplateRepository{db: db}
}

func (r *interviewTemplateRepository) List(ctx context.Context, companyID string, p pagination.Params) ([]models.InterviewTemplate, int, error) {
	const whereClause = "WHERE company_id = $1::uuid"
	args := []any{companyID}

	var total int
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM interview_templates %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("interview template list count: %w", err)
	}

	listQ := fmt.Sprintf(
		"SELECT id, company_id, name, type, duration_minutes, description, config_json, created_by, created_at, updated_at FROM interview_templates %s ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		whereClause,
	)
	args = append(args, p.PageSize, p.Offset())

	rows, err := r.db.QueryContext(ctx, listQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("interview template list query: %w", err)
	}
	defer rows.Close()

	var templates []models.InterviewTemplate
	for rows.Next() {
		var t models.InterviewTemplate
		if err := rows.Scan(
			&t.ID, &t.CompanyID, &t.Name, &t.Type, &t.DurationMinutes,
			&t.Description, &t.ConfigJSON, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("interview template scan: %w", err)
		}
		templates = append(templates, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("interview template rows: %w", err)
	}
	if templates == nil {
		templates = []models.InterviewTemplate{}
	}
	return templates, total, nil
}

func (r *interviewTemplateRepository) GetByID(ctx context.Context, companyID, templateID string) (*models.InterviewTemplate, error) {
	const q = `SELECT id, company_id, name, type, duration_minutes, description, config_json, created_by, created_at, updated_at FROM interview_templates WHERE id = $1::uuid AND company_id = $2::uuid`
	var t models.InterviewTemplate
	if err := r.db.QueryRowContext(ctx, q, templateID, companyID).Scan(
		&t.ID, &t.CompanyID, &t.Name, &t.Type, &t.DurationMinutes,
		&t.Description, &t.ConfigJSON, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("interview template get by id: %w", err)
	}
	return &t, nil
}

func (r *interviewTemplateRepository) Create(ctx context.Context, t *models.InterviewTemplate) error {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	const query = `
		INSERT INTO interview_templates (
			id, company_id, name, type, duration_minutes, description, config_json, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.CompanyID, t.Name, t.Type, t.DurationMinutes,
		t.Description, t.ConfigJSON, t.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("interview template create: %w", err)
	}
	return nil
}

func (r *interviewTemplateRepository) Update(ctx context.Context, companyID, templateID string, patch map[string]any) error {
	if len(patch) == 0 {
		return nil
	}
	allowedCols := map[string]bool{
		"name": true, "type": true, "duration_minutes": true,
		"description": true, "config_json": true,
	}
	setClauses := make([]string, 0, len(patch)+1)
	args := make([]any, 0, len(patch)+2)
	idx := 1
	for col, val := range patch {
		if !allowedCols[col] {
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, idx))
		args = append(args, val)
		idx++
	}
	if len(setClauses) == 0 {
		return nil
	}
	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, templateID, companyID)
	q := fmt.Sprintf(
		`UPDATE interview_templates SET %s WHERE id = $%d::uuid AND company_id = $%d::uuid`,
		strings.Join(setClauses, ", "), idx, idx+1,
	)
	_, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("interview template update: %w", err)
	}
	return nil
}

func (r *interviewTemplateRepository) Delete(ctx context.Context, companyID, templateID string) error {
	const q = `DELETE FROM interview_templates WHERE id = $1::uuid AND company_id = $2::uuid`
	_, err := r.db.ExecContext(ctx, q, templateID, companyID)
	if err != nil {
		return fmt.Errorf("interview template delete: %w", err)
	}
	return nil
}
