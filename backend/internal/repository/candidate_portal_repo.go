package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"

	"backend/internal/dto/response"
	"backend/internal/models"
)

type CandidatePortalRepository struct {
	db *sqlx.DB
}

func NewCandidatePortalRepository(db *sqlx.DB) *CandidatePortalRepository {
	return &CandidatePortalRepository{db: db}
}

func (r *CandidatePortalRepository) GetLatestInterviewScore(ctx context.Context, userID string) (*response.CandidateInterviewScore, error) {
	q := `SELECT COALESCE(ir.final_score, 0), COALESCE(ir.summary, ''), COALESCE(ir.strengths, '[]'::jsonb), COALESCE(ir.weaknesses, '[]'::jsonb), COALESCE(ir.report_json->'improvement_advice', ir.report_json->'suggested_next_steps', '[]'::jsonb), COALESCE(ir.report_json->>'communication_score', '0'), COALESCE(ir.report_json->>'tone_score', '0'), COALESCE(ir.report_json->>'personality_score', '0')
		FROM interview_reports ir JOIN interviews i ON i.id = ir.interview_id JOIN candidates c ON c.id = i.candidate_id
		WHERE c.user_id = $1 AND i.status = 'completed' ORDER BY ir.updated_at DESC LIMIT 1`
	var score response.CandidateInterviewScore
	var strengths, weaknesses, advice []byte
	if err := r.db.QueryRowxContext(ctx, q, userID).Scan(&score.FinalScore, &score.Summary, &strengths, &weaknesses, &advice, &score.CommunicationScore, &score.ToneScore, &score.PersonalityScore); err != nil {
		return nil, err
	}
	_ = json.Unmarshal(strengths, &score.Strengths)
	_ = json.Unmarshal(weaknesses, &score.Weaknesses)
	_ = json.Unmarshal(advice, &score.Advice)
	return &score, nil
}

func (r *CandidatePortalRepository) GetInterviewsByUserID(ctx context.Context, userID string) ([]response.CandidatePortalInterview, error) {
	q := `
		SELECT 
			i.id, i.title, i.mode, i.status, i.scheduled_at, coalesce(i.invite_token_hash, '') as invite_token_hash,
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
	q := `SELECT count(*) FROM mock_interviews WHERE user_id = $1 AND status = 'completed'`
	var count int
	err := r.db.GetContext(ctx, &count, q, userID)
	return count, err
}

func (r *CandidatePortalRepository) GetAverageMockScore(ctx context.Context, userID string) (float64, error) {
	q := `SELECT coalesce(avg(final_score), 0) FROM mock_interviews WHERE user_id = $1 AND status = 'completed'`
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
	if err == nil && fileID != "" {
		return fileID, origName, storageKey, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return "", "", "", err
	}

	// Fallback: CV uploaded on portal before the user has any candidate/application row.
	qFile := `
		SELECT id, coalesce(original_name, ''), coalesce(storage_key, '')
		FROM files
		WHERE owner_user_id = $1::uuid AND file_type = 'cv'
		ORDER BY created_at DESC
		LIMIT 1
	`
	err = r.db.QueryRowContext(ctx, qFile, userID).Scan(&fileID, &origName, &storageKey)
	if err == sql.ErrNoRows {
		return "", "", "", nil
	}
	return fileID, origName, storageKey, err
}

// GetUserCVParseStatus returns the latest parse status for the user's portal CV.
// Prefers the owned files row (works without a candidates application), then
// falls back to candidates.cv_ai_status.
func (r *CandidatePortalRepository) GetUserCVParseStatus(ctx context.Context, userID string) (string, error) {
	qFile := `
		SELECT coalesce(parse_status, '')
		FROM files
		WHERE owner_user_id = $1::uuid AND file_type = 'cv'
		ORDER BY created_at DESC
		LIMIT 1
	`
	var status string
	err := r.db.QueryRowContext(ctx, qFile, userID).Scan(&status)
	if err == nil && status != "" {
		return status, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	q := `
		SELECT coalesce(cv_ai_status, '')
		FROM candidates
		WHERE user_id = $1 AND cv_file_id IS NOT NULL
		ORDER BY updated_at DESC
		LIMIT 1
	`
	err = r.db.QueryRowContext(ctx, q, userID).Scan(&status)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return status, err
}

// SaveParsedCV stores AI output on the CV file row (always) and mirrors it to
// any candidate rows that still point at the same file.
func (r *CandidatePortalRepository) SaveParsedCV(ctx context.Context, userID, cvFileID, parsedJSON, summary string) (bool, error) {
	qFile := `
		UPDATE files
		SET parsed_json = $1::jsonb,
		    ai_summary = $2,
		    parse_status = 'ready',
		    parse_error = NULL,
		    parsed_at = NOW()
		WHERE id = $3::uuid
		  AND owner_user_id = $4::uuid
		  AND file_type = 'cv'`
	fileRes, err := r.db.ExecContext(ctx, qFile, parsedJSON, summary, cvFileID, userID)
	if err != nil {
		return false, err
	}
	fileRows, _ := fileRes.RowsAffected()

	q := `UPDATE candidates SET parsed_cv_json = $1, ai_cv_summary = $2, cv_ai_status = 'ready',
	              cv_ai_error = NULL, cv_ai_updated_at = NOW(), updated_at = NOW()
	      WHERE user_id = $3 AND cv_file_id = $4::uuid`
	candRes, err := r.db.ExecContext(ctx, q, parsedJSON, summary, userID, cvFileID)
	if err != nil {
		return fileRows > 0, err
	}
	candRows, _ := candRes.RowsAffected()
	return fileRows > 0 || candRows > 0, nil
}

// UpdateCVAIStatus records a retryable CV parsing transition for the user's
// latest owned CV file and any candidate rows owned by the portal user.
func (r *CandidatePortalRepository) UpdateCVAIStatus(ctx context.Context, userID, status string, parseErr error) error {
	var errorText sql.NullString
	if parseErr != nil {
		errorText = sql.NullString{String: parseErr.Error(), Valid: true}
	}
	qFile := `
		UPDATE files
		SET parse_status = $1::text,
		    parse_error = $2,
		    parsed_at = CASE WHEN $1::text IN ('ready', 'failed') THEN NOW() ELSE parsed_at END
		WHERE id = (
			SELECT id FROM files
			WHERE owner_user_id = $3::uuid AND file_type = 'cv'
			ORDER BY created_at DESC
			LIMIT 1
		)`
	if _, err := r.db.ExecContext(ctx, qFile, status, errorText, userID); err != nil {
		return err
	}

	q := `UPDATE candidates
	      SET cv_ai_status = $1::text,
	          cv_ai_error = $2,
	          cv_ai_attempts = cv_ai_attempts + CASE WHEN $1::text = 'processing' THEN 1 ELSE 0 END,
	          cv_ai_updated_at = NOW(), updated_at = NOW()
	      WHERE user_id = $3`
	_, err := r.db.ExecContext(ctx, q, status, errorText, userID)
	return err
}

// GetParsedCV returns the latest AI-extracted CV JSON for this user (empty string if none).
// Prefers owned CV files so portal uploads without applications still work.
func (r *CandidatePortalRepository) GetParsedCV(ctx context.Context, userID string) (string, error) {
	qFile := `
		SELECT coalesce(parsed_json::text, '')
		FROM files
		WHERE owner_user_id = $1::uuid
		  AND file_type = 'cv'
		  AND parsed_json IS NOT NULL
		ORDER BY created_at DESC
		LIMIT 1`
	var parsed string
	err := r.db.QueryRowContext(ctx, qFile, userID).Scan(&parsed)
	if err == nil && parsed != "" && parsed != "null" {
		return parsed, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	q := `SELECT coalesce(parsed_cv_json::text, '') FROM candidates
	      WHERE user_id = $1 AND parsed_cv_json IS NOT NULL
	      ORDER BY updated_at DESC LIMIT 1`
	err = r.db.QueryRowContext(ctx, q, userID).Scan(&parsed)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return parsed, err
}

// UpdateUserCV attaches the already-created upload metadata to all candidate
// applications owned by this user. The file row itself is created by
// FileService.ProcessUpload, so no duplicate/dummy metadata is inserted here.
func (r *CandidatePortalRepository) UpdateUserCV(ctx context.Context, userID string, cvFileID string, cvOriginalName string) error {
	// Mark the new file as waiting for AI parse.
	qMark := `
		UPDATE files
		SET parse_status = 'pending',
		    parse_error = NULL,
		    parsed_json = NULL,
		    ai_summary = NULL,
		    parsed_at = NULL
		WHERE id = $1::uuid AND owner_user_id = $2::uuid AND file_type = 'cv'`
	if _, err := r.db.ExecContext(ctx, qMark, cvFileID, userID); err != nil {
		return err
	}

	// The file metadata is created by FileService.ProcessUpload. Keep the
	// candidate pointing at that record instead of creating a second dummy row
	// whose storage_key would incorrectly contain the file ID.
	q := `
		UPDATE candidates
		SET cv_file_id = $1::uuid, parsed_cv_json = NULL, ai_cv_summary = NULL,
		    cv_ai_status = 'pending', cv_ai_error = NULL, cv_ai_updated_at = NOW(), updated_at = NOW()
		WHERE user_id = $2
	`
	_, err := r.db.ExecContext(ctx, q, cvFileID, userID)
	return err
}

// ClearAllUserCVs detaches every CV from the user's candidate/mock rows.
func (r *CandidatePortalRepository) ClearAllUserCVs(ctx context.Context, userID string) error {
	qCand := `
		UPDATE candidates
		SET cv_file_id = NULL,
		    parsed_cv_json = NULL,
		    ai_cv_summary = NULL,
		    cv_ai_status = 'pending',
		    cv_ai_error = NULL,
		    cv_ai_updated_at = NOW(),
		    updated_at = NOW()
		WHERE user_id = $1 AND cv_file_id IS NOT NULL`
	if _, err := r.db.ExecContext(ctx, qCand, userID); err != nil {
		return err
	}
	qMock := `
		UPDATE mock_interviews
		SET cv_file_id = NULL, updated_at = NOW()
		WHERE user_id = $1::uuid AND cv_file_id IS NOT NULL`
	_, err := r.db.ExecContext(ctx, qMock, userID)
	return err
}

type duplicateApplicationError struct{}

func (e *duplicateApplicationError) Error() string { return "duplicate application" }

// DuplicateApplication is returned by ApplyForJob when the (job, candidate) pair already exists.
var DuplicateApplication error = &duplicateApplicationError{}

func (r *CandidatePortalRepository) ApplyForJob(ctx context.Context, cID, companyID, userID, fullName, email, phone, cvFileID, cvOriginalName, jobID string) error {
	qFind := `SELECT id FROM candidates WHERE user_id = $1 AND company_id = $2 AND deleted_at IS NULL LIMIT 1`
	var existingCID string
	err := r.db.GetContext(ctx, &existingCID, qFind, userID, companyID)

	candidateID := cID
	if err == nil && existingCID != "" {
		candidateID = existingCID
		if cvFileID != "" {
			qUpd := `UPDATE candidates SET phone = COALESCE(NULLIF($1, ''), phone),
				cv_file_id = $2::uuid, updated_at = NOW() WHERE id = $3`
			_, err = r.db.ExecContext(ctx, qUpd, phone, cvFileID, candidateID)
		} else {
			qUpd := `UPDATE candidates SET phone = COALESCE(NULLIF($1, ''), phone), updated_at = NOW() WHERE id = $2`
			_, err = r.db.ExecContext(ctx, qUpd, phone, candidateID)
		}
		if err != nil {
			return err
		}
	} else {
		var cvParam interface{}
		if cvFileID != "" {
			cvParam = cvFileID
		}
		qIns := `
			INSERT INTO candidates (id, company_id, user_id, full_name, email, phone, source, status, cv_file_id)
			VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), 'manual', 'new', $7)
		`
		_, err = r.db.ExecContext(ctx, qIns, candidateID, companyID, userID, fullName, email, phone, cvParam)
		if err != nil {
			return err
		}
	}

	qJc := `
		INSERT INTO job_candidates (company_id, candidate_id, job_id, pipeline_status, applied_at)
		VALUES ($1, $2, $3, 'new', NOW())
		ON CONFLICT (job_id, candidate_id) DO NOTHING
	`
	result, err := r.db.ExecContext(ctx, qJc, companyID, candidateID, jobID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return DuplicateApplication
	}
	return nil
}

// GetCandidateIDByUserCompany resolves the candidate row id for a user in a company.
func (r *CandidatePortalRepository) GetCandidateIDByUserCompany(ctx context.Context, userID, companyID string) (string, error) {
	q := `SELECT id FROM candidates WHERE user_id = $1 AND company_id = $2 AND deleted_at IS NULL LIMIT 1`
	var id string
	err := r.db.GetContext(ctx, &id, q, userID, companyID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return id, err
}

// UpdateMatch persists the AI-computed fit score and match JSON onto a job application.
func (r *CandidatePortalRepository) UpdateMatch(ctx context.Context, jobID, candidateID string, fitScore sql.NullFloat64, matchJSON models.JSONB) error {
	q := `UPDATE job_candidates SET fit_score = $1, ai_match_json = $2 WHERE job_id = $3 AND candidate_id = $4`
	_, err := r.db.ExecContext(ctx, q, fitScore, matchJSON, jobID, candidateID)
	return err
}

// UserApplication is a (job, candidate) pair for a user's job application, used
// to recompute fit scores when the user's CV changes.
type UserApplication struct {
	JobID       string `db:"job_id"`
	CandidateID string `db:"candidate_id"`
}

// ListUserApplications returns all (job_id, candidate_id) pairs the user has
// applied to, so their fit scores can be recomputed after a CV change.
func (r *CandidatePortalRepository) ListUserApplications(ctx context.Context, userID string) ([]UserApplication, error) {
	q := `
		SELECT jc.job_id, jc.candidate_id
		FROM job_candidates jc
		JOIN candidates c ON jc.candidate_id = c.id
		WHERE c.user_id = $1 AND c.deleted_at IS NULL
	`
	var items []UserApplication
	err := r.db.SelectContext(ctx, &items, q, userID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CandidatePortalRepository) GetApplicationsByUserID(ctx context.Context, userID string) ([]response.CandidatePortalApplication, error) {
	q := `
		SELECT 
			jc.id,
			jc.job_id,
			coalesce(j.title, '') as job_title,
			jc.company_id,
			coalesce(comp.name, '') as company_name,
			jc.pipeline_status as status,
			jc.applied_at,
			coalesce(f.original_name, '') as cv_name,
			coalesce(f.storage_key, '') as cv_storage_key
		FROM job_candidates jc
		JOIN candidates c ON jc.candidate_id = c.id
		JOIN jobs j ON jc.job_id = j.id
		LEFT JOIN companies comp ON jc.company_id = comp.id
		LEFT JOIN files f ON c.cv_file_id = f.id
		WHERE c.user_id = $1
		ORDER BY jc.applied_at DESC NULLS LAST
	`
	var items []response.CandidatePortalApplication
	err := r.db.SelectContext(ctx, &items, q, userID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CandidatePortalRepository) CancelApplication(ctx context.Context, userID string, applicationID string) error {
	q := `
		DELETE FROM job_candidates jc
		USING candidates c
		WHERE jc.candidate_id = c.id AND jc.id = $1 AND c.user_id = $2
	`
	_, err := r.db.ExecContext(ctx, q, applicationID, userID)
	return err
}
