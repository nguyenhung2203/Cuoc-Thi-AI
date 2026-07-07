package repository

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/models"
)

type ScoreRepository interface {
	SaveScore(ctx context.Context, score *models.InterviewScore) error
	GetScoresByInterviewID(ctx context.Context, interviewID string) ([]models.InterviewScore, error)
}

type scoreRepository struct {
	db *sql.DB
}

func NewScoreRepository(db *sql.DB) ScoreRepository {
	return &scoreRepository{db: db}
}

func (r *scoreRepository) SaveScore(ctx context.Context, score *models.InterviewScore) error {
	query := `
		INSERT INTO interview_scores (
			interview_id, rubric_criterion_id, criterion_name, score, max_score, 
			weight, weighted_score, evidence, ai_comment, confidence, status, scored_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (interview_id, criterion_name) DO UPDATE SET
			score = EXCLUDED.score,
			max_score = EXCLUDED.max_score,
			weight = EXCLUDED.weight,
			weighted_score = EXCLUDED.weighted_score,
			evidence = EXCLUDED.evidence,
			ai_comment = EXCLUDED.ai_comment,
			confidence = EXCLUDED.confidence,
			status = EXCLUDED.status,
			scored_by = EXCLUDED.scored_by,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		score.InterviewID,
		score.RubricCriterionID,
		score.CriterionName,
		score.Score,
		score.MaxScore,
		score.Weight,
		score.WeightedScore,
		score.Evidence,
		score.AIComment,
		score.Confidence,
		score.Status,
		score.ScoredBy,
	).Scan(&score.ID, &score.CreatedAt, &score.UpdatedAt)

	if err != nil {
		return fmt.Errorf("upsert interview_score: %w", err)
	}
	return nil
}

func (r *scoreRepository) GetScoresByInterviewID(ctx context.Context, interviewID string) ([]models.InterviewScore, error) {
	query := `
		SELECT id, interview_id, rubric_criterion_id, criterion_name, score, max_score, 
			weight, weighted_score, evidence, ai_comment, confidence, status, scored_by, 
			created_at, updated_at
		FROM interview_scores
		WHERE interview_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, interviewID)
	if err != nil {
		return nil, fmt.Errorf("query interview_scores: %w", err)
	}
	defer rows.Close()

	var scores []models.InterviewScore
	for rows.Next() {
		var s models.InterviewScore
		if err := rows.Scan(
			&s.ID, &s.InterviewID, &s.RubricCriterionID, &s.CriterionName, &s.Score, &s.MaxScore,
			&s.Weight, &s.WeightedScore, &s.Evidence, &s.AIComment, &s.Confidence, &s.Status, &s.ScoredBy,
			&s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan interview_score: %w", err)
		}
		scores = append(scores, s)
	}
	return scores, nil
}
