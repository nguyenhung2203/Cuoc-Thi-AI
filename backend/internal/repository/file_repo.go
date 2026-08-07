package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"backend/internal/models"
)

type FileRepository struct {
	db *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) *FileRepository {
	return &FileRepository{db: db}
}

// Create inserts a file metadata record. f.ID must be a pre-generated UUID.
func (r *FileRepository) Create(ctx context.Context, f *models.File) (*models.File, error) {
	const q = `
		INSERT INTO files (
			id, company_id, owner_user_id, original_name,
			storage_key, mime_type, size_bytes, file_type, checksum
		) VALUES (
			$1::uuid, $2::uuid, $3::uuid, $4,
			$5, $6, $7, $8, $9
		) RETURNING *`

	var inserted models.File
	if err := r.db.QueryRowxContext(ctx, q,
		f.ID, f.CompanyID, f.OwnerUserID, f.OriginalName,
		f.StorageKey, f.MimeType, f.SizeBytes, f.FileType, f.Checksum,
	).StructScan(&inserted); err != nil {
		return nil, fmt.Errorf("file create: %w", err)
	}
	return &inserted, nil
}

// GetByID returns file metadata by ID. Returns nil, nil when not found.
// Note: storage_key is populated on the returned struct but must not be
// forwarded to API responses (it carries json:"-").
func (r *FileRepository) GetByID(ctx context.Context, fileID string) (*models.File, error) {
	const q = `SELECT * FROM files WHERE id = $1::uuid`

	var f models.File
	if err := r.db.GetContext(ctx, &f, q, fileID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("file get by id: %w", err)
	}
	return &f, nil
}

// GetByIDAndOwner returns a file only when owner_user_id matches.
// Used for candidate CV access control. Returns nil, nil when not found or owner mismatch.
func (r *FileRepository) GetByIDAndOwner(ctx context.Context, fileID, ownerUserID string) (*models.File, error) {
	const q = `
		SELECT * FROM files
		WHERE id = $1::uuid AND owner_user_id = $2::uuid`

	var f models.File
	if err := r.db.GetContext(ctx, &f, q, fileID, ownerUserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("file get by id and owner: %w", err)
	}
	return &f, nil
}

// GetByIDAndCompany returns a file when it belongs to the company directly
// or is the CV attached to a candidate belonging to that company. Portal CVs
// are owner-scoped at upload time and become visible through the candidate
// relationship when that candidate applies to a company's job.
func (r *FileRepository) GetByIDAndCompany(ctx context.Context, fileID, companyID string) (*models.File, error) {
	const q = `
		SELECT f.* FROM files f
		WHERE f.id = $1::uuid
		  AND (
			f.company_id = $2::uuid
			OR EXISTS (
				SELECT 1 FROM candidates c
				WHERE c.cv_file_id = f.id
				  AND c.company_id = $2::uuid
				  AND c.deleted_at IS NULL
			)
		  )`

	var f models.File
	if err := r.db.GetContext(ctx, &f, q, fileID, companyID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("file get by id and company: %w", err)
	}
	return &f, nil
}

func (r *FileRepository) FindByChecksum(ctx context.Context, checksum string) (*models.File, error) {
	var file models.File
	err := r.db.GetContext(ctx, &file, "SELECT * FROM files WHERE checksum = $1 LIMIT 1", checksum)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// DeleteOwnedCV deletes a CV file row when it belongs to ownerUserID.
// Returns false when the file is missing or not owned by the user.
func (r *FileRepository) DeleteOwnedCV(ctx context.Context, fileID, ownerUserID string) (bool, error) {
	const q = `
		DELETE FROM files
		WHERE id = $1::uuid
		  AND owner_user_id = $2::uuid
		  AND file_type = 'cv'`
	res, err := r.db.ExecContext(ctx, q, fileID, ownerUserID)
	if err != nil {
		return false, fmt.Errorf("file delete owned cv: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// ListOwnedCVs returns CV file metadata owned by the user (newest first).
func (r *FileRepository) ListOwnedCVs(ctx context.Context, ownerUserID string) ([]models.File, error) {
	const q = `
		SELECT * FROM files
		WHERE owner_user_id = $1::uuid AND file_type = 'cv'
		ORDER BY created_at DESC`
	var files []models.File
	if err := r.db.SelectContext(ctx, &files, q, ownerUserID); err != nil {
		return nil, fmt.Errorf("list owned cvs: %w", err)
	}
	return files, nil
}

// SaveCVParse stores AI-parsed CV JSON on the owning file row.
func (r *FileRepository) SaveCVParse(ctx context.Context, fileID, ownerUserID, parsedJSON, summary string) (bool, error) {
	const q = `
		UPDATE files
		SET parsed_json = $1::jsonb,
		    ai_summary = $2,
		    parse_status = 'ready',
		    parse_error = NULL,
		    parsed_at = NOW()
		WHERE id = $3::uuid
		  AND owner_user_id = $4::uuid
		  AND file_type = 'cv'`
	res, err := r.db.ExecContext(ctx, q, parsedJSON, summary, fileID, ownerUserID)
	if err != nil {
		return false, fmt.Errorf("file save cv parse: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// UpdateCVParseStatus sets parse lifecycle fields on an owned CV file.
func (r *FileRepository) UpdateCVParseStatus(ctx context.Context, fileID, ownerUserID, status string, parseErr error) error {
	var errorText any
	if parseErr != nil {
		errorText = parseErr.Error()
	}
	const q = `
		UPDATE files
		SET parse_status = $1,
		    parse_error = $2,
		    parsed_at = CASE WHEN $1 IN ('ready', 'failed') THEN NOW() ELSE parsed_at END
		WHERE id = $3::uuid
		  AND owner_user_id = $4::uuid
		  AND file_type = 'cv'`
	_, err := r.db.ExecContext(ctx, q, status, errorText, fileID, ownerUserID)
	if err != nil {
		return fmt.Errorf("file update cv parse status: %w", err)
	}
	return nil
}

// GetLatestOwnedCVParse returns parse_status + parsed_json for the newest owned CV.
func (r *FileRepository) GetLatestOwnedCVParse(ctx context.Context, ownerUserID string) (status string, parsedJSON string, err error) {
	const q = `
		SELECT coalesce(parse_status, ''), coalesce(parsed_json::text, '')
		FROM files
		WHERE owner_user_id = $1::uuid AND file_type = 'cv'
		ORDER BY created_at DESC
		LIMIT 1`
	err = r.db.QueryRowContext(ctx, q, ownerUserID).Scan(&status, &parsedJSON)
	if err == sql.ErrNoRows {
		return "", "", nil
	}
	if err != nil {
		return "", "", fmt.Errorf("file get latest cv parse: %w", err)
	}
	return status, parsedJSON, nil
}
