package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"backend/internal/models"
	"backend/internal/pkg/pagination"
)

type QuestionRepository interface {
	List(ctx context.Context, companyID string, jobID, qType, level, keyword string, p pagination.Params) ([]models.QuestionBank, int64, error)
	GetByID(ctx context.Context, companyID, id string) (*models.QuestionBank, error)
	Create(ctx context.Context, q *models.QuestionBank) error
	CreateQuestions(ctx context.Context, questions []models.QuestionBank) error
	Update(ctx context.Context, companyID, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, companyID, id string) error
}

type questionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) List(ctx context.Context, companyID string, jobID, qType, level, keyword string, p pagination.Params) ([]models.QuestionBank, int64, error) {
	where := []string{"company_id = $1"}
	args := []interface{}{companyID}
	argID := 2

	if jobID != "" {
		where = append(where, fmt.Sprintf("job_id = $%d", argID))
		args = append(args, jobID)
		argID++
	}
	if qType != "" {
		where = append(where, fmt.Sprintf("question_type = $%d", argID))
		args = append(args, qType)
		argID++
	}
	if level != "" {
		where = append(where, fmt.Sprintf("level = $%d", argID))
		args = append(args, level)
		argID++
	}
	if keyword != "" {
		where = append(where, fmt.Sprintf("question_text ILIKE $%d", argID))
		args = append(args, "%"+keyword+"%")
		argID++
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	countQuery := "SELECT count(*) FROM question_bank " + whereClause
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count query: %w", err)
	}

	if total == 0 {
		return []models.QuestionBank{}, 0, nil
	}

	query := fmt.Sprintf(`
		SELECT id, company_id, job_id, created_by, question_text, question_type, 
			   skill_tags, level, expected_signals, is_ai_generated, created_at, updated_at
		FROM question_bank
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argID, argID+1)
	args = append(args, p.PageSize, p.Offset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("select query: %w", err)
	}
	defer rows.Close()

	var items []models.QuestionBank
	for rows.Next() {
		var q models.QuestionBank
		if err := rows.Scan(
			&q.ID, &q.CompanyID, &q.JobID, &q.CreatedBy, &q.QuestionText, &q.QuestionType,
			&q.SkillTags, &q.Level, &q.ExpectedSignals, &q.IsAIGenerated, &q.CreatedAt, &q.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan: %w", err)
		}
		items = append(items, q)
	}

	return items, total, nil
}

func (r *questionRepository) GetByID(ctx context.Context, companyID, id string) (*models.QuestionBank, error) {
	query := `
		SELECT id, company_id, job_id, created_by, question_text, question_type, 
			   skill_tags, level, expected_signals, is_ai_generated, created_at, updated_at
		FROM question_bank
		WHERE id = $1 AND company_id = $2
	`
	var q models.QuestionBank
	err := r.db.QueryRowContext(ctx, query, id, companyID).Scan(
		&q.ID, &q.CompanyID, &q.JobID, &q.CreatedBy, &q.QuestionText, &q.QuestionType,
		&q.SkillTags, &q.Level, &q.ExpectedSignals, &q.IsAIGenerated, &q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("get by id query: %w", err)
	}
	return &q, nil
}

func (r *questionRepository) Create(ctx context.Context, q *models.QuestionBank) error {
	query := `
		INSERT INTO question_bank (
			id, company_id, job_id, created_by, question_text, question_type,
			skill_tags, level, expected_signals, is_ai_generated, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		q.ID, q.CompanyID, q.JobID, q.CreatedBy, q.QuestionText, q.QuestionType,
		q.SkillTags, q.Level, q.ExpectedSignals, q.IsAIGenerated, q.CreatedAt, q.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create query: %w", err)
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
			id, company_id, job_id, created_by, question_text, question_type,
			skill_tags, level, expected_signals, is_ai_generated, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, q := range questions {
		_, err := stmt.ExecContext(ctx,
			q.ID, q.CompanyID, q.JobID, q.CreatedBy, q.QuestionText, q.QuestionType,
			q.SkillTags, q.Level, q.ExpectedSignals, q.IsAIGenerated, q.CreatedAt, q.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("execute statement: %w", err)
		}
	}

	return tx.Commit()
}

func (r *questionRepository) Update(ctx context.Context, companyID, id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClauses := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+2)
	argID := 1

	for k, v := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", k, argID))
		args = append(args, v)
		argID++
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = NOW()"))

	query := fmt.Sprintf("UPDATE question_bank SET %s WHERE id = $%d AND company_id = $%d",
		strings.Join(setClauses, ", "), argID, argID+1)
	
	args = append(args, id, companyID)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update query: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *questionRepository) Delete(ctx context.Context, companyID, id string) error {
	query := "DELETE FROM question_bank WHERE id = $1 AND company_id = $2"
	res, err := r.db.ExecContext(ctx, query, id, companyID)
	if err != nil {
		return fmt.Errorf("delete query: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
