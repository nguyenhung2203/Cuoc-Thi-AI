package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
	"backend/internal/pkg/pagination"
)

type CandidateRepository struct {
	db *sqlx.DB
}

func NewCandidateRepository(db *sqlx.DB) *CandidateRepository {
	return &CandidateRepository{db: db}
}

// List returns candidates for a company with optional jobID/status/keyword filters.
// When jobID is provided, joins job_candidates to filter and include fit_score.
func (r *CandidateRepository) List(ctx context.Context, companyID, jobID, status, keyword string, p pagination.Params) ([]models.Candidate, int, error) {
	args := []any{companyID}
	argIdx := 2

	var from string
	var where []string

	if jobID != "" {
		from = `FROM candidates c
		         JOIN job_candidates jc ON jc.candidate_id = c.id AND jc.job_id = $%d::uuid`
		from = fmt.Sprintf(from, argIdx)
		args = append(args, jobID)
		argIdx++
		where = append(where, "c.company_id = $1::uuid", "c.deleted_at IS NULL")
	} else {
		from = "FROM candidates c"
		where = append(where, "c.company_id = $1::uuid", "c.deleted_at IS NULL")
	}

	if status != "" {
		where = append(where, fmt.Sprintf("c.status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}
	if keyword != "" {
		where = append(where, fmt.Sprintf("(c.full_name ILIKE $%d OR c.email ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+keyword+"%")
		argIdx++
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) %s %s", from, whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("candidate list count: %w", err)
	}

	sortBy := "c.created_at"
	allowedSort := map[string]bool{"created_at": true, "updated_at": true, "full_name": true, "email": true, "status": true}
	if allowedSort[p.SortBy] {
		sortBy = "c." + p.SortBy
	}
	sortDir := "ASC"
	if strings.ToUpper(p.SortDir) == "DESC" {
		sortDir = "DESC"
	}

	selectCols := `
		c.*,
		(
			SELECT jc.job_id 
			FROM job_candidates jc 
			WHERE jc.candidate_id = c.id 
			ORDER BY jc.created_at DESC LIMIT 1
		) as latest_job_id,
		(
			SELECT j.title 
			FROM job_candidates jc 
			JOIN jobs j ON j.id = jc.job_id 
			WHERE jc.candidate_id = c.id 
			ORDER BY jc.created_at DESC LIMIT 1
		) as latest_job_title
	`
	listQuery := fmt.Sprintf(
		"SELECT %s %s %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		selectCols, from, whereClause, sortBy, sortDir, argIdx, argIdx+1,
	)
	args = append(args, p.PageSize, p.Offset())

	var candidates []models.Candidate
	if err := r.db.SelectContext(ctx, &candidates, listQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("candidate list select: %w", err)
	}
	if candidates == nil {
		candidates = []models.Candidate{}
	}

	return candidates, total, nil
}

// GetByID returns a candidate scoped to companyID. Returns nil, nil when not found.
func (r *CandidateRepository) GetByID(ctx context.Context, companyID, candidateID string) (*models.Candidate, error) {
	q := `
		SELECT c.*,
		(
			SELECT jc.job_id 
			FROM job_candidates jc 
			WHERE jc.candidate_id = c.id 
			ORDER BY jc.created_at DESC LIMIT 1
		) as latest_job_id,
		(
			SELECT j.title 
			FROM job_candidates jc 
			JOIN jobs j ON j.id = jc.job_id 
			WHERE jc.candidate_id = c.id 
			ORDER BY jc.created_at DESC LIMIT 1
		) as latest_job_title,
		f.original_name as cv_original_name
		FROM candidates c
		LEFT JOIN files f ON c.cv_file_id = f.id
		WHERE c.id = $1::uuid
		  AND c.company_id = $2::uuid
		  AND c.deleted_at IS NULL`

	var c models.Candidate
	if err := r.db.GetContext(ctx, &c, q, candidateID, companyID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("candidate get by id: %w", err)
	}
	return &c, nil
}

// EmailExists checks if an email already exists within the company (excluding deleted rows).
func (r *CandidateRepository) EmailExists(ctx context.Context, companyID, email string) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1 FROM candidates
			WHERE company_id = $1::uuid
			  AND email = $2
			  AND deleted_at IS NULL
		)`

	var exists bool
	if err := r.db.QueryRowContext(ctx, q, companyID, email).Scan(&exists); err != nil {
		return false, fmt.Errorf("candidate email exists: %w", err)
	}
	return exists, nil
}

// Create inserts a new candidate. c.ID must be a pre-generated UUID. Returns the inserted row.
func (r *CandidateRepository) Create(ctx context.Context, c *models.Candidate) (*models.Candidate, error) {
	const q = `
		INSERT INTO candidates (
			id, company_id, user_id, full_name, email, phone,
			avatar_url, cv_file_id, parsed_cv_json, ai_cv_summary,
			source, status, tags, created_by
		) VALUES (
			$1::uuid, $2::uuid, $3::uuid, $4, $5, $6,
			$7, $8::uuid, $9, $10,
			$11, $12, $13, $14::uuid
		) RETURNING *`

	var inserted models.Candidate
	if err := r.db.QueryRowxContext(ctx, q,
		c.ID, c.CompanyID, c.UserID, c.FullName, c.Email, c.Phone,
		c.AvatarURL, c.CVFileID, c.ParsedCVJSON, c.AICVSummary,
		c.Source, c.Status, c.Tags, c.CreatedBy,
	).StructScan(&inserted); err != nil {
		return nil, fmt.Errorf("candidate create: %w", err)
	}
	return &inserted, nil
}

// Update applies a patch map to a candidate. Returns nil, nil when not found.
func (r *CandidateRepository) Update(ctx context.Context, companyID, candidateID string, patch map[string]any) (*models.Candidate, error) {
	if len(patch) == 0 {
		return r.GetByID(ctx, companyID, candidateID)
	}

	allowedCols := map[string]bool{
		"full_name": true, "phone": true, "avatar_url": true, "cv_file_id": true,
		"parsed_cv_json": true, "ai_cv_summary": true, "source": true,
		"status": true, "tags": true,
	}

	setClauses := make([]string, 0, len(patch)+1)
	args := make([]any, 0, len(patch)+2)
	idx := 1

	for col, val := range patch {
		if !allowedCols[col] {
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, idx))
		args = append(args, val)
		idx++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, companyID, candidateID)
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	args = append(args, candidateID, companyID)
	q := fmt.Sprintf(
		`UPDATE candidates SET %s
		 WHERE id = $%d::uuid AND company_id = $%d::uuid AND deleted_at IS NULL
		 RETURNING *`,
		strings.Join(setClauses, ", "), idx, idx+1,
	)

	var updated models.Candidate
	if err := r.db.QueryRowxContext(ctx, q, args...).StructScan(&updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("candidate update: %w", err)
	}
	return &updated, nil
}

// SoftDelete sets deleted_at = NOW() scoped to companyID.
func (r *CandidateRepository) SoftDelete(ctx context.Context, companyID, candidateID string) error {
	const q = `
		UPDATE candidates SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1::uuid
		  AND company_id = $2::uuid
		  AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, q, candidateID, companyID)
	if err != nil {
		return fmt.Errorf("candidate soft delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("candidate soft delete: not found or already deleted")
	}
	return nil
}

// GetJobCandidate returns a job_candidates row by jobID+candidateID. Returns nil, nil when not found.
func (r *CandidateRepository) GetJobCandidate(ctx context.Context, jobID, candidateID string) (*models.JobCandidate, error) {
	const q = `
		SELECT * FROM job_candidates
		WHERE job_id = $1::uuid AND candidate_id = $2::uuid`

	var jc models.JobCandidate
	if err := r.db.GetContext(ctx, &jc, q, jobID, candidateID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get job candidate: %w", err)
	}
	return &jc, nil
}

// CreateJobCandidate inserts a job_candidates row. jc.ID must be a pre-generated UUID.
func (r *CandidateRepository) CreateJobCandidate(ctx context.Context, jc *models.JobCandidate) (*models.JobCandidate, error) {
	const q = `
		INSERT INTO job_candidates (
			id, company_id, job_id, candidate_id,
			pipeline_status, fit_score, ai_match_json, applied_at, created_by
		) VALUES (
			$1::uuid, $2::uuid, $3::uuid, $4::uuid,
			$5, $6, $7, $8, $9::uuid
		) RETURNING *`

	var inserted models.JobCandidate
	if err := r.db.QueryRowxContext(ctx, q,
		jc.ID, jc.CompanyID, jc.JobID, jc.CandidateID,
		jc.PipelineStatus, jc.FitScore, jc.AIMatchJSON, jc.AppliedAt, jc.CreatedBy,
	).StructScan(&inserted); err != nil {
		return nil, fmt.Errorf("create job candidate: %w", err)
	}
	return &inserted, nil
}

// UpdateJobCandidateStatus updates pipeline_status for a job_candidates row.
func (r *CandidateRepository) UpdateJobCandidateStatus(ctx context.Context, companyID, jobID, candidateID, status string) error {
	const q = `
		UPDATE job_candidates SET pipeline_status = $1, updated_at = NOW()
		WHERE job_id = $2::uuid AND candidate_id = $3::uuid AND company_id = $4::uuid`

	res, err := r.db.ExecContext(ctx, q, status, jobID, candidateID, companyID)
	if err != nil {
		return fmt.Errorf("update job candidate status: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("update job candidate status: row not found or access denied")
	}
	return nil
}

// DeleteJobCandidate removes a candidate from a job (unassign).
// Hard-deletes the job_candidates row scoped to company.
func (r *CandidateRepository) DeleteJobCandidate(ctx context.Context, companyID, jobID, candidateID string) error {
	const q = `
		DELETE FROM job_candidates
		WHERE job_id = $1::uuid
		  AND candidate_id = $2::uuid
		  AND company_id = $3::uuid`

	res, err := r.db.ExecContext(ctx, q, jobID, candidateID, companyID)
	if err != nil {
		return fmt.Errorf("delete job candidate: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("delete job candidate: row not found or access denied")
	}
	return nil
}

// ListByJobID returns job_candidates rows (with pipeline info) for a given job scoped to companyID.
func (r *CandidateRepository) ListByJobID(ctx context.Context, companyID, jobID string, p pagination.Params) ([]models.JobCandidate, int, error) {
	const countQ = `
		SELECT COUNT(*) FROM job_candidates
		WHERE job_id = $1::uuid AND company_id = $2::uuid`

	var total int
	if err := r.db.QueryRowContext(ctx, countQ, jobID, companyID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("list by job id count: %w", err)
	}

	sortBy := "created_at"
	allowedSort := map[string]bool{"created_at": true, "updated_at": true, "pipeline_status": true, "fit_score": true}
	if allowedSort[p.SortBy] {
		sortBy = p.SortBy
	}
	sortDir := "ASC"
	if strings.ToUpper(p.SortDir) == "DESC" {
		sortDir = "DESC"
	}

	listQ := fmt.Sprintf(`
		SELECT * FROM job_candidates
		WHERE job_id = $1::uuid AND company_id = $2::uuid
		ORDER BY %s %s
		LIMIT $3 OFFSET $4`,
		sortBy, sortDir,
	)

	var rows []models.JobCandidate
	if err := r.db.SelectContext(ctx, &rows, listQ, jobID, companyID, p.PageSize, p.Offset()); err != nil {
		return nil, 0, fmt.Errorf("list by job id select: %w", err)
	}
	if rows == nil {
		rows = []models.JobCandidate{}
	}

	return rows, total, nil
}
