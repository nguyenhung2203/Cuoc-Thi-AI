package repository

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
)

type TranscriptRepository struct {
	db *sqlx.DB
}

func NewTranscriptRepository(db *sqlx.DB) *TranscriptRepository {
	return &TranscriptRepository{db: db}
}

func (r *TranscriptRepository) Create(ctx context.Context, transcript *models.InterviewTranscript) error {
	q := `
		INSERT INTO interview_transcripts (
			id, interview_id, participant_id, speaker_type, speaker_name, content,
			language, start_time_ms, end_time_ms, confidence, source, is_final, created_at
		) VALUES (
			:id, :interview_id, :participant_id, :speaker_type, :speaker_name, :content,
			:language, :start_time_ms, :end_time_ms, :confidence, :source, :is_final, :created_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, q, transcript)
	return err
}

func (r *TranscriptRepository) ListByInterview(ctx context.Context, interviewID string) ([]models.InterviewTranscript, error) {
	q := `
		SELECT id, interview_id, participant_id, speaker_type, speaker_name, content,
			language, start_time_ms, end_time_ms, confidence, source, is_final,
			edited_content, edited_by, edited_at, created_at
		FROM interview_transcripts
		WHERE interview_id = $1
		ORDER BY start_time_ms ASC NULLS LAST, created_at ASC
	`
	var items []models.InterviewTranscript
	err := r.db.SelectContext(ctx, &items, q, interviewID)
	return items, err
}

func (r *TranscriptRepository) UpdateEditedContent(ctx context.Context, id, interviewID, editedContent, editedBy string) error {
	q := `
		UPDATE interview_transcripts
		SET edited_content = $1, edited_by = $2, edited_at = NOW()
		WHERE id = $3 AND interview_id = $4
	`
	res, err := r.db.ExecContext(ctx, q, editedContent, editedBy, id, interviewID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return errors.New("transcript not found or access denied")
	}
	return nil
}
