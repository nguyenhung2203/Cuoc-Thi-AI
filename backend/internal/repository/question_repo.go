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

type QuestionRepository interface {
	CreateQuestions(ctx context.Context, questions []models.QuestionBank) error
	List(ctx context.Context, companyID string, jobID, questionType, level, keyword string, p pagination.Params) ([]models.QuestionBank, int, error)
	GetByID(ctx context.Context, companyID, questionID string) (*models.QuestionBank, error)
	Create(ctx context.Context, q *models.QuestionBank) error
	Update(ctx context.Context, companyID, questionID string, patch map[string]any) error
	Delete(ctx context.Context, companyID, questionID string) error
}

type questionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) List(ctx context.Context, companyID string, jobID, questionType, level, keyword string, p pagination.Params) ([]models.QuestionBank, int, error) {
	args := []any{companyID}
	argIdx := 2
	where := []string{"company_id = $1::uuid"}

	if jobID != "" {
		where = append(where, fmt.Sprintf("(job_id = $%d::uuid OR job_id IS NULL)", argIdx))
		args = append(args, jobID)
		argIdx++
	}
	if questionType != "" {
		where = append(where, fmt.Sprintf("question_type = $%d", argIdx))
		args = append(args, questionType)
		argIdx++
	}
	if level != "" {
		where = append(where, fmt.Sprintf("level = $%d", argIdx))
		args = append(args, level)
		argIdx++
	}
	if keyword != "" {
		where = append(where, fmt.Sprintf("question_text ILIKE $%d", argIdx))
		args = append(args, "%"+keyword+"%")
		argIdx++
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	var total int
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM question_bank %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("question list count: %w", err)
	}

	listQ := fmt.Sprintf(
		"SELECT * FROM question_bank %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		whereClause, argIdx, argIdx+1,
	)
	args = append(args, p.PageSize, p.Offset())

	rows, err := r.db.QueryContext(ctx, listQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("question list query: %w", err)
	}
	defer rows.Close()

	var questions []models.QuestionBank
	for rows.Next() {
		var q models.QuestionBank
		if err := rows.Scan(
			&q.ID, &q.CompanyID, &q.JobID, &q.CreatedBy,
			&q.QuestionText, &q.QuestionType, &q.SkillTags,
			&q.Level, &q.ExpectedSignals, &q.IsAIGenerated,
			&q.CreatedAt, &q.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("question scan: %w", err)
		}
		questions = append(questions, q)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("question rows: %w", err)
	}
	if questions == nil {
		questions = []models.QuestionBank{}
	}
	return questions, total, nil
}

func (r *questionRepository) GetByID(ctx context.Context, companyID, questionID string) (*models.QuestionBank, error) {
	const q = `SELECT * FROM question_bank WHERE id = $1::uuid AND company_id = $2::uuid`
	var qb models.QuestionBank
	if err := r.db.QueryRowContext(ctx, q, questionID, companyID).Scan(
		&qb.ID, &qb.CompanyID, &qb.JobID, &qb.CreatedBy,
		&qb.QuestionText, &qb.QuestionType, &qb.SkillTags,
		&qb.Level, &qb.ExpectedSignals, &qb.IsAIGenerated,
		&qb.CreatedAt, &qb.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("question get by id: %w", err)
	}
	return &qb, nil
}

func (r *questionRepository) Create(ctx context.Context, q *models.QuestionBank) error {
	if q.ID == "" {
		q.ID = uuid.NewString()
	}
	const query = `
		INSERT INTO question_bank (
			id, company_id, job_id, created_by, question_text, question_type,
			skill_tags, level, expected_signals, is_ai_generated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		q.ID, q.CompanyID, q.JobID, q.CreatedBy, q.QuestionText, q.QuestionType,
		q.SkillTags, q.Level, q.ExpectedSignals, q.IsAIGenerated,
	)
	if err != nil {
		return fmt.Errorf("question create: %w", err)
	}
	return nil
}

func (r *questionRepository) Update(ctx context.Context, companyID, questionID string, patch map[string]any) error {
	if len(patch) == 0 {
		return nil
	}
	allowedCols := map[string]bool{
		"question_text": true, "question_type": true, "skill_tags": true,
		"level": true, "expected_signals": true,
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
	args = append(args, questionID, companyID)
	q := fmt.Sprintf(
		`UPDATE question_bank SET %s WHERE id = $%d::uuid AND company_id = $%d::uuid`,
		strings.Join(setClauses, ", "), idx, idx+1,
	)
	_, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("question update: %w", err)
	}
	return nil
}

func (r *questionRepository) Delete(ctx context.Context, companyID, questionID string) error {
	const q = `DELETE FROM question_bank WHERE id = $1::uuid AND company_id = $2::uuid`
	_, err := r.db.ExecContext(ctx, q, questionID, companyID)
	if err != nil {
		return fmt.Errorf("question delete: %w", err)
	}
	return nil
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
