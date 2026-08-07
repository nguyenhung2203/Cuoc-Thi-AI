package repository

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/models"
)

type RubricRepository interface {
	CreateRubricWithCriteria(ctx context.Context, rubric *models.Rubric, criteria []models.RubricCriteria) error
	UpdateRubricWithCriteria(ctx context.Context, rubric *models.Rubric, criteria []models.RubricCriteria) error
	GetByID(ctx context.Context, companyID, rubricID string) (*models.Rubric, []models.RubricCriteria, error)
	ListByCompany(ctx context.Context, companyID, jobID string) ([]models.Rubric, error)
	Delete(ctx context.Context, companyID, rubricID string) error
	GetCriteriaByIDs(ctx context.Context, criteriaIDs []string) ([]models.RubricCriteria, error)
}

type rubricRepository struct {
	db *sql.DB
}

func NewRubricRepository(db *sql.DB) RubricRepository {
	return &rubricRepository{db: db}
}

func (r *rubricRepository) CreateRubricWithCriteria(ctx context.Context, rubric *models.Rubric, criteria []models.RubricCriteria) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert rubric
	rubricQuery := `
		INSERT INTO rubrics (company_id, job_id, name, description, total_weight, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err = tx.QueryRowContext(ctx, rubricQuery,
		rubric.CompanyID, rubric.JobID, rubric.Name, rubric.Description, rubric.TotalWeight, rubric.CreatedBy,
	).Scan(&rubric.ID)
	if err != nil {
		return fmt.Errorf("insert rubric: %w", err)
	}

	if len(criteria) == 0 {
		return tx.Commit()
	}

	// Insert criteria
	criteriaQuery := `
		INSERT INTO rubric_criteria (
			rubric_id, name, description, weight, min_score, max_score, scoring_guide, order_index
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	stmt, err := tx.PrepareContext(ctx, criteriaQuery)
	if err != nil {
		return fmt.Errorf("prepare criteria statement: %w", err)
	}
	defer stmt.Close()

	for _, c := range criteria {
		_, err := stmt.ExecContext(ctx,
			rubric.ID, c.Name, c.Description, c.Weight, c.MinScore, c.MaxScore, c.ScoringGuide, c.OrderIndex,
		)
		if err != nil {
			return fmt.Errorf("insert criteria: %w", err)
		}
	}

	return tx.Commit()
}

// UpdateRubricWithCriteria updates a rubric's fields and replaces its criteria
// atomically. Scoped by company_id to prevent cross-tenant edits.
func (r *rubricRepository) UpdateRubricWithCriteria(ctx context.Context, rubric *models.Rubric, criteria []models.RubricCriteria) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update rubric row (only if it belongs to the company).
	res, err := tx.ExecContext(ctx, `
		UPDATE rubrics SET name = $1, description = $2, job_id = $3, total_weight = $4, updated_at = NOW()
		WHERE id = $5 AND company_id = $6`,
		rubric.Name, rubric.Description, rubric.JobID, rubric.TotalWeight, rubric.ID, rubric.CompanyID,
	)
	if err != nil {
		return fmt.Errorf("update rubric: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}

	// Replace criteria: delete existing, insert the new set.
	if _, err := tx.ExecContext(ctx, `DELETE FROM rubric_criteria WHERE rubric_id = $1`, rubric.ID); err != nil {
		return fmt.Errorf("clear criteria: %w", err)
	}

	if len(criteria) > 0 {
		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO rubric_criteria (rubric_id, name, description, weight, min_score, max_score, scoring_guide, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
		if err != nil {
			return fmt.Errorf("prepare criteria: %w", err)
		}
		defer stmt.Close()
		for _, c := range criteria {
			if _, err := stmt.ExecContext(ctx,
				rubric.ID, c.Name, c.Description, c.Weight, c.MinScore, c.MaxScore, c.ScoringGuide, c.OrderIndex,
			); err != nil {
				return fmt.Errorf("insert criteria: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (r *rubricRepository) GetByID(ctx context.Context, companyID, rubricID string) (*models.Rubric, []models.RubricCriteria, error) {
	// 1. Get rubric
	rubricQ := `SELECT id, company_id, job_id, name, description, total_weight, created_by, created_at, updated_at
	            FROM rubrics WHERE id = $1 AND company_id = $2`
	var rubric models.Rubric
	err := r.db.QueryRowContext(ctx, rubricQ, rubricID, companyID).Scan(
		&rubric.ID, &rubric.CompanyID, &rubric.JobID, &rubric.Name, &rubric.Description,
		&rubric.TotalWeight, &rubric.CreatedBy, &rubric.CreatedAt, &rubric.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil // not found
		}
		return nil, nil, err
	}

	// 2. Get criteria
	critQ := `SELECT id, rubric_id, name, description, weight, min_score, max_score, scoring_guide, order_index, created_at, updated_at
	          FROM rubric_criteria WHERE rubric_id = $1 ORDER BY order_index ASC`
	rows, err := r.db.QueryContext(ctx, critQ, rubric.ID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var criteria []models.RubricCriteria
	for rows.Next() {
		var c models.RubricCriteria
		err := rows.Scan(
			&c.ID, &c.RubricID, &c.Name, &c.Description, &c.Weight, &c.MinScore, &c.MaxScore,
			&c.ScoringGuide, &c.OrderIndex, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		criteria = append(criteria, c)
	}

	return &rubric, criteria, nil
}

func (r *rubricRepository) ListByCompany(ctx context.Context, companyID, jobID string) ([]models.Rubric, error) {
	query := `SELECT id, company_id, job_id, name, description, total_weight, created_by, created_at, updated_at
	          FROM rubrics WHERE company_id = $1`
	var args []interface{}
	args = append(args, companyID)

	if jobID != "" {
		query += ` AND (job_id = $2 OR job_id IS NULL)`
		args = append(args, jobID)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Rubric
	for rows.Next() {
		var rubric models.Rubric
		if err := rows.Scan(
			&rubric.ID, &rubric.CompanyID, &rubric.JobID, &rubric.Name, &rubric.Description,
			&rubric.TotalWeight, &rubric.CreatedBy, &rubric.CreatedAt, &rubric.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, rubric)
	}
	return list, nil
}

func (r *rubricRepository) Delete(ctx context.Context, companyID, rubricID string) error {
	// Note: using hard delete or soft delete depending on requirement.
	// Since rubric might be referenced by interview_scores, usually we shouldn't hard delete if used.
	// We'll use a hard delete here for simplicity, assuming cascade or no references if deleted early.
	// Wait, we should probably soft delete if DATABASE_DESIGN has deleted_at.
	// Let's check DATABASE_DESIGN: "Rubrics" has no deleted_at. So hard delete.

	// Wait, we need to delete criteria first? No, CASCADE should handle it if set, otherwise manual.
	// Let's delete criteria manually just in case.
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Ensure the rubric belongs to the company
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT true FROM rubrics WHERE id = $1 AND company_id = $2", rubricID, companyID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // nothing to delete
		}
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM rubric_criteria WHERE rubric_id = $1", rubricID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM rubrics WHERE id = $1", rubricID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *rubricRepository) GetCriteriaByIDs(ctx context.Context, criteriaIDs []string) ([]models.RubricCriteria, error) {
	if len(criteriaIDs) == 0 {
		return nil, nil
	}

	query := `SELECT id, rubric_id, name, description, weight, min_score, max_score, scoring_guide, order_index, created_at, updated_at
	          FROM rubric_criteria WHERE id IN (`

	args := make([]interface{}, len(criteriaIDs))
	for i, id := range criteriaIDs {
		args[i] = id
		query += fmt.Sprintf("$%d", i+1)
		if i < len(criteriaIDs)-1 {
			query += ", "
		}
	}
	query += `)`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var criteria []models.RubricCriteria
	for rows.Next() {
		var c models.RubricCriteria
		if err := rows.Scan(
			&c.ID, &c.RubricID, &c.Name, &c.Description, &c.Weight, &c.MinScore, &c.MaxScore,
			&c.ScoringGuide, &c.OrderIndex, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		criteria = append(criteria, c)
	}

	return criteria, nil
}
