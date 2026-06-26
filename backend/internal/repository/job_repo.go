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

type JobRepository struct {
	db *sqlx.DB
}

func NewJobRepository(db *sqlx.DB) *JobRepository {
	return &JobRepository{db: db}
}

// JobListResult holds paginated jobs with total count.
type JobListResult struct {
	Jobs  []models.Job
	Total int
}

// List returns jobs for a company with optional status/keyword filters.
func (r *JobRepository) List(ctx context.Context, companyID string, status, keyword string, p pagination.Params) (*JobListResult, error) {
	args := []any{companyID}
	argIdx := 2

	where := []string{"company_id = $1::uuid", "deleted_at IS NULL"}

	if status != "" {
		where = append(where, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}
	if keyword != "" {
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+keyword+"%")
		argIdx++
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM jobs %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("job list count: %w", err)
	}

	sortBy := "created_at"
	allowedSort := map[string]bool{"created_at": true, "updated_at": true, "title": true, "status": true}
	if allowedSort[p.SortBy] {
		sortBy = p.SortBy
	}
	sortDir := "ASC"
	if strings.ToUpper(p.SortDir) == "DESC" {
		sortDir = "DESC"
	}

	listQuery := fmt.Sprintf(
		"SELECT * FROM jobs %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, sortBy, sortDir, argIdx, argIdx+1,
	)
	args = append(args, p.PageSize, p.Offset())

	var jobs []models.Job
	if err := r.db.SelectContext(ctx, &jobs, listQuery, args...); err != nil {
		return nil, fmt.Errorf("job list select: %w", err)
	}
	if jobs == nil {
		jobs = []models.Job{}
	}

	return &JobListResult{Jobs: jobs, Total: total}, nil
}

// GetByID returns a job scoped to companyID. Returns nil, nil when not found.
func (r *JobRepository) GetByID(ctx context.Context, companyID, jobID string) (*models.Job, error) {
	const q = `
		SELECT * FROM jobs
		WHERE id = $1::uuid
		  AND company_id = $2::uuid
		  AND deleted_at IS NULL`

	var job models.Job
	if err := r.db.GetContext(ctx, &job, q, jobID, companyID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("job get by id: %w", err)
	}
	return &job, nil
}

// Create inserts a new job. job.ID must be a pre-generated UUID. Returns the inserted row.
func (r *JobRepository) Create(ctx context.Context, job *models.Job) (*models.Job, error) {
	const q = `
		INSERT INTO jobs (
			id, company_id, title, department, level, location,
			employment_type, salary_min, salary_max, currency,
			description, requirements, benefits, status,
			ai_summary, ai_analysis_json, created_by
		) VALUES (
			$1::uuid, $2::uuid, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14,
			$15, $16, $17::uuid
		) RETURNING *`

	var inserted models.Job
	if err := r.db.QueryRowxContext(ctx, q,
		job.ID, job.CompanyID, job.Title, job.Department, job.Level, job.Location,
		job.EmploymentType, job.SalaryMin, job.SalaryMax, job.Currency,
		job.Description, job.Requirements, job.Benefits, job.Status,
		job.AISummary, job.AIAnalysisJSON, job.CreatedBy,
	).StructScan(&inserted); err != nil {
		return nil, fmt.Errorf("job create: %w", err)
	}
	return &inserted, nil
}

// Update applies non-nil pointer fields from a patch map. Only keys present in
// the map are included in the SET clause. Returns nil, nil when not found.
func (r *JobRepository) Update(ctx context.Context, companyID, jobID string, patch map[string]any) (*models.Job, error) {
	if len(patch) == 0 {
		return r.GetByID(ctx, companyID, jobID)
	}

	allowedCols := map[string]bool{
		"title": true, "department": true, "level": true, "location": true,
		"employment_type": true, "salary_min": true, "salary_max": true, "currency": true,
		"description": true, "requirements": true, "benefits": true, "status": true,
		"ai_summary": true, "ai_analysis_json": true,
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
		return r.GetByID(ctx, companyID, jobID)
	}

	setClauses = append(setClauses, "updated_at = NOW()")

	args = append(args, jobID, companyID)
	q := fmt.Sprintf(
		`UPDATE jobs SET %s
		 WHERE id = $%d::uuid AND company_id = $%d::uuid AND deleted_at IS NULL
		 RETURNING *`,
		strings.Join(setClauses, ", "), idx, idx+1,
	)

	var updated models.Job
	if err := r.db.QueryRowxContext(ctx, q, args...).StructScan(&updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("job update: %w", err)
	}
	return &updated, nil
}

// SoftDelete sets deleted_at = NOW() scoped to companyID.
func (r *JobRepository) SoftDelete(ctx context.Context, companyID, jobID string) error {
	const q = `
		UPDATE jobs SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1::uuid
		  AND company_id = $2::uuid
		  AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, q, jobID, companyID)
	if err != nil {
		return fmt.Errorf("job soft delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("job soft delete: not found or already deleted")
	}
	return nil
}

// CountCandidates returns how many job_candidates rows exist for a job.
func (r *JobRepository) CountCandidates(ctx context.Context, jobID string) (int, error) {
	const q = `SELECT COUNT(*) FROM job_candidates WHERE job_id = $1::uuid`
	var count int
	if err := r.db.QueryRowContext(ctx, q, jobID).Scan(&count); err != nil {
		return 0, fmt.Errorf("job count candidates: %w", err)
	}
	return count, nil
}

// CountInterviews returns how many interviews rows exist for a job.
// Returns 0 without error if the interviews table does not yet exist.
func (r *JobRepository) CountInterviews(ctx context.Context, jobID string) (int, error) {
	const existsQ = `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = 'interviews'
		)`
	var tableExists bool
	if err := r.db.QueryRowContext(ctx, existsQ).Scan(&tableExists); err != nil {
		return 0, fmt.Errorf("job count interviews (table check): %w", err)
	}
	if !tableExists {
		return 0, nil
	}

	const q = `SELECT COUNT(*) FROM interviews WHERE job_id = $1::uuid`
	var count int
	if err := r.db.QueryRowContext(ctx, q, jobID).Scan(&count); err != nil {
		return 0, fmt.Errorf("job count interviews: %w", err)
	}
	return count, nil
}

// SaveAIAnalysisWithTx saves the AI analysis results (Job update, Rubric, Questions) in a single transaction.
func (r *JobRepository) SaveAIAnalysisWithTx(
	ctx context.Context,
	companyID, jobID string,
	aiSummary string,
	aiAnalysisJSON []byte,
	rubric *models.Rubric,
	criteria []models.RubricCriteria,
	questions []models.QuestionBank,
) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// 1. Lock Job row and check if already analyzed
	var currentJSON []byte
	checkQ := `
		SELECT ai_analysis_json FROM jobs
		WHERE id = $1::uuid AND company_id = $2::uuid AND deleted_at IS NULL
		FOR UPDATE
	`
	err = tx.QueryRowContext(ctx, checkQ, jobID, companyID).Scan(&currentJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("job not found or deleted")
		}
		return fmt.Errorf("check job: %w", err)
	}

	if currentJSON != nil && string(currentJSON) != "null" && string(currentJSON) != "\"null\"" {
		return fmt.Errorf("CONFLICT: job has already been analyzed")
	}

	// 2. Update Job
	updateQ := `
		UPDATE jobs
		SET ai_summary = $1, ai_analysis_json = $2, updated_at = NOW()
		WHERE id = $3::uuid AND company_id = $4::uuid AND deleted_at IS NULL
	`
	_, err = tx.ExecContext(ctx, updateQ, sql.NullString{String: aiSummary, Valid: true}, aiAnalysisJSON, jobID, companyID)
	if err != nil {
		return fmt.Errorf("update job: %w", err)
	}

	// 3. Insert Rubric & Criteria
	if rubric != nil {
		rubricQuery := `
			INSERT INTO rubrics (id, company_id, job_id, name, description, total_weight, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err = tx.ExecContext(ctx, rubricQuery,
			rubric.ID, rubric.CompanyID, rubric.JobID, rubric.Name, rubric.Description, rubric.TotalWeight, rubric.CreatedBy,
		)
		if err != nil {
			return fmt.Errorf("insert rubric: %w", err)
		}

		if len(criteria) > 0 {
			criteriaQuery := `
				INSERT INTO rubric_criteria (
					id, rubric_id, name, description, weight, min_score, max_score, scoring_guide, order_index
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`
			stmt, err := tx.PrepareContext(ctx, criteriaQuery)
			if err != nil {
				return fmt.Errorf("prepare criteria statement: %w", err)
			}

			for _, c := range criteria {
				_, err := stmt.ExecContext(ctx,
					c.ID, rubric.ID, c.Name, c.Description, c.Weight, c.MinScore, c.MaxScore, c.ScoringGuide, c.OrderIndex,
				)
				if err != nil {
					stmt.Close()
					return fmt.Errorf("insert criteria: %w", err)
				}
			}
			stmt.Close()
		}
	}

	// 4. Insert Questions
	if len(questions) > 0 {
		qQuery := `
			INSERT INTO question_bank (
				company_id, job_id, created_by, question_text, question_type,
				skill_tags, level, expected_signals, is_ai_generated
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`
		stmt, err := tx.PrepareContext(ctx, qQuery)
		if err != nil {
			return fmt.Errorf("prepare question statement: %w", err)
		}

		for _, q := range questions {
			_, err := stmt.ExecContext(ctx,
				q.CompanyID, q.JobID, q.CreatedBy, q.QuestionText, q.QuestionType,
				q.SkillTags, q.Level, q.ExpectedSignals, q.IsAIGenerated,
			)
			if err != nil {
				stmt.Close()
				return fmt.Errorf("insert question: %w", err)
			}
		}
		stmt.Close()
	}

	return tx.Commit()
}
