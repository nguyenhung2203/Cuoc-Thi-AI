package repository

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/models"
)

type QuestionRepository interface {
	CreateQuestions(ctx context.Context, questions []models.QuestionBank) error
}

type questionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) CreateQuestions(ctx context.Context, questions []models.QuestionBank) error {
	if len(questions) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO question_bank (
			company_id, job_id, created_by, question_text, question_type,
			skill_tags, level, expected_signals, is_ai_generated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, q := range questions {
		_, err := stmt.ExecContext(ctx,
			q.CompanyID, q.JobID, q.CreatedBy, q.QuestionText, q.QuestionType,
			q.SkillTags, q.Level, q.ExpectedSignals, q.IsAIGenerated,
		)
		if err != nil {
			return fmt.Errorf("execute statement: %w", err)
		}
	}

	return tx.Commit()
}
