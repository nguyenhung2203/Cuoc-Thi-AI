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
