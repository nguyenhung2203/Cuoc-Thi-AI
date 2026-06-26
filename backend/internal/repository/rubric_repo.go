package repository

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/models"
)

type RubricRepository interface {
	CreateRubricWithCriteria(ctx context.Context, rubric *models.Rubric, criteria []models.RubricCriteria) error
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
