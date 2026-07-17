package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
)

type MockRepository struct {
	db *sqlx.DB
}

func NewMockRepository(db *sqlx.DB) *MockRepository {
	return &MockRepository{db: db}
}

// --- Mock Interviews ---

func (r *MockRepository) Create(ctx context.Context, m *models.MockInterview) error {
	q := `
		INSERT INTO mock_interviews (id, user_id, target_role, target_level, cv_file_id, status, created_at, updated_at)
		VALUES (:id, :user_id, :target_role, :target_level, :cv_file_id, :status, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	stmt, err := r.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.QueryRowContext(ctx, m).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *MockRepository) GetByID(ctx context.Context, id string) (*models.MockInterview, error) {
	q := `SELECT * FROM mock_interviews WHERE id = $1`
	var m models.MockInterview
	err := r.db.GetContext(ctx, &m, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *MockRepository) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]models.MockInterview, error) {
	q := `SELECT * FROM mock_interviews WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	var items []models.MockInterview
	err := r.db.SelectContext(ctx, &items, q, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MockRepository) UpdateStatus(ctx context.Context, id, status string) error {
	q := `UPDATE mock_interviews SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, status, id)
	return err
}

// UpdateStatusIf atomically updates status only when current status equals expected.
// Returns true if a row was updated.
func (r *MockRepository) UpdateStatusIf(ctx context.Context, id, status, expectedCurrentStatus string) (bool, error) {
	q := `UPDATE mock_interviews SET status = $1, updated_at = NOW() WHERE id = $2 AND status = $3`
	res, err := r.db.ExecContext(ctx, q, status, id, expectedCurrentStatus)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

func (r *MockRepository) UpdateStartedAt(ctx context.Context, id string, t time.Time) error {
	q := `UPDATE mock_interviews SET started_at = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, t, id)
	return err
}

func (r *MockRepository) UpdateEndedAt(ctx context.Context, id string, t time.Time) error {
	q := `UPDATE mock_interviews SET ended_at = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, t, id)
	return err
}

func (r *MockRepository) SaveFeedback(ctx context.Context, id string, feedbackJSON models.JSONB, finalScore sql.NullFloat64) error {
	q := `UPDATE mock_interviews SET feedback_json = $1, final_score = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, q, feedbackJSON, finalScore, id)
	return err
}

// --- Mock Messages ---

func (r *MockRepository) CreateMessage(ctx context.Context, msg *models.MockInterviewMessage) error {
	q := `
		INSERT INTO mock_interview_messages (id, mock_interview_id, sender_type, content, question_type, score_json, created_at)
		VALUES (:id, :mock_interview_id, :sender_type, :content, :question_type, :score_json, NOW())
		RETURNING id, created_at
	`
	stmt, err := r.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.QueryRowContext(ctx, msg).Scan(&msg.ID, &msg.CreatedAt)
}

// UpdateMessageScore sets score_json on a single message (used by AI scoring at End).
func (r *MockRepository) UpdateMessageScore(ctx context.Context, messageID string, scoreJSON models.JSONB) error {
	q := `UPDATE mock_interview_messages SET score_json = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, scoreJSON, messageID)
	return err
}

func (r *MockRepository) ListMessages(ctx context.Context, mockID string) ([]models.MockInterviewMessage, error) {
	q := `SELECT * FROM mock_interview_messages WHERE mock_interview_id = $1 ORDER BY created_at ASC`
	var items []models.MockInterviewMessage
	err := r.db.SelectContext(ctx, &items, q, mockID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MockRepository) CountMessagesBySender(ctx context.Context, mockID, senderType string) (int, error) {
	q := `SELECT COUNT(*) FROM mock_interview_messages WHERE mock_interview_id = $1 AND sender_type = $2`
	var count int
	err := r.db.GetContext(ctx, &count, q, mockID, senderType)
	return count, err
}
