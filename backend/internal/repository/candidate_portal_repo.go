package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"backend/internal/dto/response"
)

type CandidatePortalRepository struct {
	db *sqlx.DB
}

func NewCandidatePortalRepository(db *sqlx.DB) *CandidatePortalRepository {
	return &CandidatePortalRepository{db: db}
}

func (r *CandidatePortalRepository) GetInterviewsByUserID(ctx context.Context, userID string) ([]response.CandidatePortalInterview, error) {
	q := `
		SELECT 
			i.id, i.title, i.mode, i.status, i.scheduled_at, i.invite_token_hash,
			i.company_id,
			coalesce(comp.name, '') as company_name,
			i.job_id,
			coalesce(j.title, '') as job_title
		FROM interviews i
		JOIN candidates c ON i.candidate_id = c.id
		LEFT JOIN companies comp ON i.company_id = comp.id
		LEFT JOIN jobs j ON i.job_id = j.id
		WHERE c.user_id = $1
		ORDER BY i.scheduled_at ASC NULLS LAST
	`
	var items []response.CandidatePortalInterview
	err := r.db.SelectContext(ctx, &items, q, userID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CandidatePortalRepository) GetCompletedMockCount(ctx context.Context, userID string) (int, error) {
	q := `SELECT count(*) FROM mock_results WHERE user_id = $1`
	var count int
	err := r.db.GetContext(ctx, &count, q, userID)
	return count, err
}

func (r *CandidatePortalRepository) GetAverageMockScore(ctx context.Context, userID string) (float64, error) {
	q := `SELECT coalesce(avg(score), 0) FROM mock_results WHERE user_id = $1`
	var avg float64
	err := r.db.GetContext(ctx, &avg, q, userID)
	return avg, err
}

func (r *CandidatePortalRepository) GetUserLatestCV(ctx context.Context, userID string) (string, string, string, error) {
	q := `
		SELECT c.cv_file_id, coalesce(f.original_name, '') as cv_original_name, coalesce(f.storage_key, '') as storage_key
		FROM candidates c
		LEFT JOIN files f ON c.cv_file_id = f.id
		WHERE c.user_id = $1 AND c.cv_file_id IS NOT NULL
		ORDER BY c.updated_at DESC
		LIMIT 1
	`
	var fileID, origName, storageKey string
	err := r.db.QueryRowContext(ctx, q, userID).Scan(&fileID, &origName, &storageKey)
	if err == sql.ErrNoRows {
		return "", "", "", nil
	}
	return fileID, origName, storageKey, err
}

func (r *CandidatePortalRepository) UpdateUserCV(ctx context.Context, userID string, cvFileID string, cvOriginalName string) error {
	// Insert a dummy file record to get a valid UUID for the foreign key
	var actualFileID string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO files (owner_user_id, original_name, storage_key, mime_type, size_bytes, file_type)
		VALUES ($1, $2, $3, 'application/pdf', 0, 'cv')
		RETURNING id
	`, userID, cvOriginalName, cvFileID).Scan(&actualFileID)
	if err != nil {
		return err
	}

	q := `
		UPDATE candidates 
		SET cv_file_id = $1, updated_at = NOW()
		WHERE user_id = $2
	`
	_, err = r.db.ExecContext(ctx, q, actualFileID, userID)
	return err
}

func (r *CandidatePortalRepository) ApplyForJob(ctx context.Context, cID, companyID, userID, fullName, email, phone, cvFileID, cvOriginalName, jobID string) error {
	qFind := `SELECT id FROM candidates WHERE user_id = $1 AND company_id = $2 AND deleted_at IS NULL LIMIT 1`
	var existingCID string
	err := r.db.GetContext(ctx, &existingCID, qFind, userID, companyID)

	candidateID := cID
	if err == nil && existingCID != "" {
		candidateID = existingCID
		qUpd := `UPDATE candidates SET phone = COALESCE(NULLIF($1, ''), phone), updated_at = NOW() WHERE id = $2`
		_, err = r.db.ExecContext(ctx, qUpd, phone, candidateID)
		if err != nil {
			return err
		}
	} else {
		qIns := `
			INSERT INTO candidates (id, company_id, user_id, full_name, email, phone, source, status)
			VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), 'manual', 'new')
		`
		_, err = r.db.ExecContext(ctx, qIns, candidateID, companyID, userID, fullName, email, phone)
		if err != nil {
			return err
		}
	}

	qJc := `
		INSERT INTO job_candidates (company_id, candidate_id, job_id, pipeline_status)
		VALUES ($1, $2, $3, 'new')
		ON CONFLICT (job_id, candidate_id) DO NOTHING
	`
	_, err = r.db.ExecContext(ctx, qJc, companyID, candidateID, jobID)
	return err
}
