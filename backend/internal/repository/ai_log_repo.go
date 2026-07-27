package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
)

type AILogRepository struct {
	db *sqlx.DB
}

func NewAILogRepository(db *sqlx.DB) *AILogRepository {
	return &AILogRepository{db: db}
}

// Create inserts a new AI request log. It does not return the inserted row because
// logs are usually written asynchronously and we only care about success/failure.
func (r *AILogRepository) Create(ctx context.Context, log *models.AIRequestLog) error {
	const q = `
		INSERT INTO ai_request_logs (
			company_id, template_id, template_version, 
			interview_id, job_id, candidate_id, 
			input_json, output_json, latency_ms, 
			tokens_in, tokens_out, cost, status, error
		) VALUES (
			$1::uuid, $2, $3, 
			$4, $5, $6, 
			$7, $8, $9, 
			$10, $11, $12, $13, $14
		)`

	_, err := r.db.ExecContext(ctx, q,
		log.CompanyID, log.TemplateID, log.TemplateVersion,
		log.InterviewID, log.JobID, log.CandidateID,
		log.InputJSON, log.OutputJSON, log.LatencyMs,
		log.TokensIn, log.TokensOut, log.Cost, log.Status, log.Error,
	)

	if err != nil {
		return fmt.Errorf("insert ai request log: %w", err)
	}

	return nil
}

// AILogFilter narrows an AI request-log query. Empty fields are ignored.
type AILogFilter struct {
	InterviewID string
	JobID       string
	CandidateID string
	Limit       int
	Offset      int
}

// List returns AI request logs matching the filter, newest first. Used by the
// admin AI-logs endpoint for debugging and cost auditing.
func (r *AILogRepository) List(ctx context.Context, f AILogFilter) ([]models.AIRequestLog, error) {
	q := `SELECT * FROM ai_request_logs WHERE 1=1`
	args := []interface{}{}
	i := 1

	if f.InterviewID != "" {
		q += fmt.Sprintf(" AND interview_id = $%d", i)
		args = append(args, f.InterviewID)
		i++
	}
	if f.JobID != "" {
		q += fmt.Sprintf(" AND job_id = $%d", i)
		args = append(args, f.JobID)
		i++
	}
	if f.CandidateID != "" {
		q += fmt.Sprintf(" AND candidate_id = $%d", i)
		args = append(args, f.CandidateID)
		i++
	}

	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	q += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, f.Limit, f.Offset)

	logs := []models.AIRequestLog{}
	if err := r.db.SelectContext(ctx, &logs, q, args...); err != nil {
		return nil, fmt.Errorf("list ai request logs: %w", err)
	}
	return logs, nil
}
