package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
)

type AILogRepository struct{ db *sqlx.DB }

func NewAILogRepository(db *sqlx.DB) *AILogRepository { return &AILogRepository{db: db} }

func (r *AILogRepository) Create(ctx context.Context, log *models.AIRequestLog) error {
	const q = `
		INSERT INTO ai_request_logs (
			company_id, template_id, template_version, interview_id, job_id, candidate_id,
			input_json, output_json, latency_ms, tokens_in, tokens_out, total_tokens,
			provider, model, operation, retry_count, cost, status, error
		) VALUES (NULLIF($1, '')::uuid, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`
	_, err := r.db.ExecContext(ctx, q,
		log.CompanyID, log.TemplateID, log.TemplateVersion, log.InterviewID, log.JobID, log.CandidateID,
		log.InputJSON, log.OutputJSON, log.LatencyMs, log.TokensIn, log.TokensOut, log.TotalTokens,
		log.Provider, log.Model, log.Operation, log.RetryCount, log.Cost, log.Status, log.Error,
	)
	if err != nil {
		return fmt.Errorf("insert ai request log: %w", err)
	}
	return nil
}

type AILogFilter struct {
	InterviewID, JobID, CandidateID string
	Limit, Offset                   int
}

func (r *AILogRepository) List(ctx context.Context, f AILogFilter) ([]models.AIRequestLog, error) {
	q := `SELECT * FROM ai_request_logs WHERE 1=1`
	args := []interface{}{}
	i := 1
	for _, filter := range []struct{ value, column string }{{f.InterviewID, "interview_id"}, {f.JobID, "job_id"}, {f.CandidateID, "candidate_id"}} {
		if filter.value != "" {
			q += fmt.Sprintf(" AND %s = $%d", filter.column, i)
			args = append(args, filter.value)
			i++
		}
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
