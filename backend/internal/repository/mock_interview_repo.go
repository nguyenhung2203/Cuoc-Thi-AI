package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
)

type MockInterviewRepository struct {
	db *sqlx.DB
}

func NewMockInterviewRepository(db *sqlx.DB) *MockInterviewRepository {
	return &MockInterviewRepository{db: db}
}

func (r *MockInterviewRepository) Create(ctx context.Context, mi *models.MockInterview) error {
	q := `
		INSERT INTO mock_interviews (user_id, target_role, target_level, cv_file_id, status)
		VALUES (:user_id, :target_role, :target_level, :cv_file_id, :status)
		RETURNING id, created_at, updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, q, mi)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.StructScan(mi)
	}
	return nil
}

func (r *MockInterviewRepository) GetByID(ctx context.Context, id string) (*models.MockInterview, error) {
	var mi models.MockInterview
	err := r.db.GetContext(ctx, &mi, "SELECT * FROM mock_interviews WHERE id = $1", id)
	return &mi, err
}

func (r *MockInterviewRepository) AddMessage(ctx context.Context, msg *models.MockInterviewMessage) error {
	q := `
		INSERT INTO mock_interview_messages (mock_interview_id, sender_type, content, question_type, score_json)
		VALUES (:mock_interview_id, :sender_type, :content, :question_type, :score_json)
		RETURNING id, created_at
	`
	rows, err := r.db.NamedQueryContext(ctx, q, msg)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.StructScan(msg)
	}
	return nil
}

func (r *MockInterviewRepository) GetMessages(ctx context.Context, mockInterviewID string) ([]models.MockInterviewMessage, error) {
	var msgs []models.MockInterviewMessage
	err := r.db.SelectContext(ctx, &msgs, "SELECT * FROM mock_interview_messages WHERE mock_interview_id = $1 ORDER BY created_at ASC", mockInterviewID)
	return msgs, err
}

func (r *MockInterviewRepository) Finish(ctx context.Context, id string, finalScore float64, feedback string) error {
	q := `
		UPDATE mock_interviews
		SET status = 'completed', ended_at = NOW(), final_score = $1, feedback_json = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, q, finalScore, feedback, id)
	return err
}
