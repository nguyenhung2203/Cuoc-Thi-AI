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

// UpdateNotes saves the recruiter's internal notes for an interview (company-scoped).
func (r *InterviewRepository) UpdateNotes(ctx context.Context, id, companyID, notes string) error {
	q := `UPDATE interviews SET recruiter_notes = $1, updated_at = NOW() WHERE id = $2 AND company_id = $3`
	_, err := r.db.ExecContext(ctx, q, notes, id, companyID)
	return err
}

// UpdateScheduledAt reschedules an interview (company-scoped).
func (r *InterviewRepository) UpdateScheduledAt(ctx context.Context, id, companyID string, scheduledAt time.Time) error {
	q := `UPDATE interviews SET scheduled_at = $1, updated_at = NOW() WHERE id = $2 AND company_id = $3`
	_, err := r.db.ExecContext(ctx, q, scheduledAt, id, companyID)
	return err
}

func (r *InterviewRepository) UpdateRoomStatus(ctx context.Context, id, status string) error {
	q := `UPDATE interview_rooms SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, status, id)
	return err
}

// GetCompanyIDByInterviewID returns the company that owns an interview.
// The realtime gateway needs it to call company-scoped AI services (IDOR guard)
// since the LiveKit room token carries only interview_id, not company_id.
func (r *InterviewRepository) GetCompanyIDByInterviewID(ctx context.Context, interviewID string) (string, error) {
	var companyID string
	err := r.db.GetContext(ctx, &companyID, `SELECT company_id FROM interviews WHERE id = $1`, interviewID)
	if err != nil {
		return "", err
	}
	return companyID, nil
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

type CandidateJoinRepoInfo struct {
	InterviewID     string          `db:"interview_id"`
	RoomID          string          `db:"room_id"`
	CandidateID     string          `db:"candidate_id"`
	UserID          *string         `db:"user_id"`
	CandidateName   string          `db:"candidate_name"`
	CompanyName     string          `db:"company_name"`
	JobTitle        string          `db:"job_title"`
	ScheduledAt     time.Time       `db:"scheduled_at"`
	InviteExpiresAt time.Time       `db:"invite_expires_at"`
}

func (r *InterviewRepository) GetCandidateJoinInfoByInviteTokenHash(ctx context.Context, hash string) (*CandidateJoinRepoInfo, error) {
	q := `
		SELECT 
			i.id as interview_id,
			coalesce(ir.room_code, 'room-' || i.id::text) as room_id,
			i.candidate_id as candidate_id,
			c.user_id as user_id,
			c.full_name as candidate_name,
			coalesce(comp.name, '') as company_name,
			coalesce(j.title, '') as job_title,
			coalesce(i.scheduled_at, NOW()) as scheduled_at,
			coalesce(i.invite_expires_at, NOW()) as invite_expires_at
		FROM interviews i
		JOIN candidates c ON i.candidate_id = c.id
		LEFT JOIN companies comp ON i.company_id = comp.id
		LEFT JOIN jobs j ON i.job_id = j.id
		LEFT JOIN interview_rooms ir ON ir.interview_id = i.id
		WHERE i.invite_token_hash = $1 AND i.status != 'cancelled' AND i.status != 'completed'
	`
	var info CandidateJoinRepoInfo
	err := r.db.GetContext(ctx, &info, q, hash)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
