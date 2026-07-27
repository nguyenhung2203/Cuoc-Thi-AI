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
