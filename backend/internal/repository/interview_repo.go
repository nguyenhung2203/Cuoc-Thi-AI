package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
	"backend/internal/pkg/errors"
)

type InterviewRepository struct {
	db *sqlx.DB
}

func NewInterviewRepository(db *sqlx.DB) *InterviewRepository {
	return &InterviewRepository{db: db}
}

// Create generates both the interview and the interview room inside a single transaction.
func (r *InterviewRepository) Create(ctx context.Context, interview *models.Interview, room *models.InterviewRoom) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.NewInternal("failed to begin transaction")
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	q1 := `
		INSERT INTO interviews (
			id, company_id, job_id, candidate_id, recruiter_id, template_id, rubric_id,
			mode, title, scheduled_at, status, room_id, invite_token_hash, invite_expires_at,
			consent_recording, consent_ai, created_by, created_at, updated_at
		) VALUES (
			:id, :company_id, :job_id, :candidate_id, :recruiter_id, :template_id, :rubric_id,
			:mode, :title, :scheduled_at, :status, :room_id, :invite_token_hash, :invite_expires_at,
			:consent_recording, :consent_ai, :created_by, :created_at, :updated_at
		)
	`
	if _, err = tx.NamedExecContext(ctx, q1, interview); err != nil {
		return errors.NewInternal("failed to insert interview: " + err.Error())
	}

	q2 := `
		INSERT INTO interview_rooms (
			id, interview_id, room_code, status, provider, connection_config,
			created_at, updated_at
		) VALUES (
			:id, :interview_id, :room_code, :status, :provider, :connection_config,
			:created_at, :updated_at
		)
	`
	if _, err = tx.NamedExecContext(ctx, q2, room); err != nil {
		return errors.NewInternal("failed to insert interview room: " + err.Error())
	}

	return nil
}

func (r *InterviewRepository) GetByIDAndCompany(ctx context.Context, id, companyID string) (*models.Interview, error) {
	q := `SELECT * FROM interviews WHERE id = $1 AND company_id = $2`
	var i models.Interview
	err := r.db.GetContext(ctx, &i, q, id, companyID)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *InterviewRepository) ListByCompany(ctx context.Context, companyID string, limit, offset int) ([]models.Interview, error) {
	q := `SELECT * FROM interviews WHERE company_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	var items []models.Interview
	err := r.db.SelectContext(ctx, &items, q, companyID, limit, offset)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *InterviewRepository) GetRoomByInterviewID(ctx context.Context, interviewID string) (*models.InterviewRoom, error) {
	q := `SELECT * FROM interview_rooms WHERE interview_id = $1`
	var room models.InterviewRoom
	err := r.db.GetContext(ctx, &room, q, interviewID)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *InterviewRepository) UpdateStatus(ctx context.Context, id, status string) error {
	q := `UPDATE interviews SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, status, id)
	return err
}

func (r *InterviewRepository) UpdateRoomStatus(ctx context.Context, id, status string) error {
	q := `UPDATE interview_rooms SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, status, id)
	return err
}

func (r *InterviewRepository) GetByInviteTokenHash(ctx context.Context, hash string) (*models.Interview, error) {
	q := `SELECT * FROM interviews WHERE invite_token_hash = $1 AND status != 'cancelled' AND status != 'completed'`
	var i models.Interview
	err := r.db.GetContext(ctx, &i, q, hash)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

type CandidateJoinInfo struct {
	InterviewID              string    `db:"interview_id" json:"interview_id"`
	RoomID                   string    `db:"room_id" json:"room_id"`
	CandidateName            string    `db:"candidate_name" json:"candidate_name"`
	CandidateID              string    `db:"candidate_id" json:"candidate_id"`
	JobTitle                 string    `db:"job_title" json:"job_title"`
	ScheduledAt              time.Time `db:"scheduled_at" json:"scheduled_at"`
	RequiresConsentAI        bool      `db:"consent_ai" json:"requires_consent_ai"`
	RequiresConsentRecording bool      `db:"consent_recording" json:"requires_consent_recording"`
	RoomCode                 string    `db:"room_code" json:"room_code"`
	Status                   string    `db:"status" json:"status"`

	RoomAccessToken          string    `json:"room_access_token"`
	RoomAccessTokenExpiresAt string    `json:"room_access_token_expires_at"`
}

func (r *InterviewRepository) GetCandidateJoinInfo(ctx context.Context, hash string) (*CandidateJoinInfo, error) {
	q := `
		SELECT 
			i.id as interview_id,
			ir.id as room_id,
			coalesce(c.full_name, '') as candidate_name,
			c.id as candidate_id,
			coalesce(j.title, '') as job_title,
			i.scheduled_at,
			i.consent_ai,
			i.consent_recording,
			ir.room_code,
			ir.status
		FROM interviews i
		LEFT JOIN interview_rooms ir ON ir.interview_id = i.id
		LEFT JOIN candidates c ON c.id = i.candidate_id
		LEFT JOIN jobs j ON j.id = i.job_id
		WHERE i.invite_token_hash = $1 
		  AND i.status != 'cancelled' 
		  AND i.status != 'completed'
	`
	var info CandidateJoinInfo
	err := r.db.GetContext(ctx, &info, q, hash)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *InterviewRepository) UpdateReportStatus(ctx context.Context, id, status string) error {
	q := `UPDATE interviews SET report_status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, status, id)
	return err
}

// UpdateReportStatusIf atomically updates report_status only when current status equals expected.
// Returns true if a row was updated, false if no row matched the condition.
func (r *InterviewRepository) UpdateReportStatusIf(ctx context.Context, id, status, expectedCurrentStatus string) (bool, error) {
	q := `UPDATE interviews SET report_status = $1, updated_at = NOW() WHERE id = $2 AND report_status = $3`
	res, err := r.db.ExecContext(ctx, q, status, id, expectedCurrentStatus)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}
