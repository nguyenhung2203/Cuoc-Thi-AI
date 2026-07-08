package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
)

type ReportRepository interface {
	UpsertReport(ctx context.Context, report *models.InterviewReport) error
	GetByInterviewID(ctx context.Context, interviewID string) (*models.InterviewReport, error)
	UpdateRecruiterDecision(ctx context.Context, interviewID string, decision string, comment string) error
}

type reportRepository struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) UpsertReport(ctx context.Context, report *models.InterviewReport) error {
	q := `
		INSERT INTO interview_reports (
			interview_id, summary, final_score, recommendation, strengths, weaknesses, risks,
			evidence_json, ai_reasoning_summary, report_json, generated_by
		) VALUES (
			:interview_id, :summary, :final_score, :recommendation, :strengths, :weaknesses, :risks,
			:evidence_json, :ai_reasoning_summary, :report_json, :generated_by
		)
		ON CONFLICT (interview_id) DO UPDATE SET
			summary = EXCLUDED.summary,
			final_score = EXCLUDED.final_score,
			recommendation = EXCLUDED.recommendation,
			strengths = EXCLUDED.strengths,
			weaknesses = EXCLUDED.weaknesses,
			risks = EXCLUDED.risks,
			evidence_json = EXCLUDED.evidence_json,
			ai_reasoning_summary = EXCLUDED.ai_reasoning_summary,
			report_json = EXCLUDED.report_json,
			generated_by = EXCLUDED.generated_by,
			generated_at = NOW(),
			updated_at = NOW()
		RETURNING id, generated_at, created_at, updated_at
	`
	
	stmt, err := r.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	err = stmt.QueryRowContext(ctx, report).Scan(&report.ID, &report.GeneratedAt, &report.CreatedAt, &report.UpdatedAt)
	return err
}

func (r *reportRepository) GetByInterviewID(ctx context.Context, interviewID string) (*models.InterviewReport, error) {
	q := `SELECT * FROM interview_reports WHERE interview_id = $1`
	var report models.InterviewReport
	err := r.db.GetContext(ctx, &report, q, interviewID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) UpdateRecruiterDecision(ctx context.Context, interviewID string, decision string, comment string) error {
	q := `
		UPDATE interview_reports
		SET recruiter_decision = $1, recruiter_comment = $2, updated_at = NOW()
		WHERE interview_id = $3
	`

	res, err := r.db.ExecContext(ctx, q, decision, comment, interviewID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return errors.New("report not found")
	}
	return nil
}
