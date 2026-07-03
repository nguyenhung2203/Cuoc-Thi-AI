package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
)

type ReportRepository struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetByInterviewID(ctx context.Context, interviewID string) (*models.InterviewReport, error) {
	q := `SELECT * FROM interview_reports WHERE interview_id = $1 LIMIT 1`
	var report models.InterviewReport
	err := r.db.GetContext(ctx, &report, q, interviewID)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) Create(ctx context.Context, report *models.InterviewReport) error {
	q := `
		INSERT INTO interview_reports (
			interview_id, summary, final_score, recommendation, strengths, weaknesses, risks, evidence_json, ai_reasoning_summary, report_json, generated_by
		) VALUES (
			:interview_id, :summary, :final_score, :recommendation, :strengths, :weaknesses, :risks, :evidence_json, :ai_reasoning_summary, :report_json, :generated_by
		) RETURNING id, generated_at, created_at, updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, q, report)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.StructScan(report)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReportRepository) UpdateDecision(ctx context.Context, interviewID string, decision string, comment string) error {
	q := `
		UPDATE interview_reports
		SET recruiter_decision = $1, recruiter_comment = $2, updated_at = NOW()
		WHERE interview_id = $3
	`
	_, err := r.db.ExecContext(ctx, q, decision, comment, interviewID)
	return err
}
