package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"database/sql"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"

	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/repository"
)

type FileService struct {
	fileRepo *repository.FileRepository
	baseURL  string
}

func NewFileService(fileRepo *repository.FileRepository, baseURL string) *FileService {
	return &FileService{fileRepo: fileRepo, baseURL: baseURL}
}

// PublicURL builds a browser-reachable URL for a stored file key.
func (s *FileService) PublicURL(storageKey string) string {
	return strings.TrimRight(s.baseURL, "/") + "/uploads/" + storageKey
}

func (s *FileService) ProcessUpload(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID, companyID, fileType string) (*models.File, error) {
	// Calc SHA-256
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, errors.NewInternal("failed to hash file")
	}
	checksum := hex.EncodeToString(hasher.Sum(nil))
	file.Seek(0, 0)

	// Check if exists in DB
	existing, err := s.fileRepo.FindByChecksum(ctx, checksum)
	if err == nil && existing != nil {
		// Also verify it actually exists on disk
		uploadDir := filepath.Join(".", "uploads")
		existingPath := filepath.Join(uploadDir, existing.StorageKey)
		if _, err := os.Stat(existingPath); err == nil {
			return s.SaveMetadata(ctx, companyID, userID, header.Filename, existing.StorageKey, existing.MimeType, header.Size, fileType, checksum)
		}
	}

	uploadDir := filepath.Join(".", "uploads")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, errors.NewInternal("failed to create upload directory")
	}

	storageKey := uuid.NewString() + filepath.Ext(header.Filename)
	filePath := filepath.Join(uploadDir, storageKey)
	
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, errors.NewInternal("failed to save file")
	}
	defer dst.Close()
	
	if _, err := io.Copy(dst, file); err != nil {
		return nil, errors.NewInternal("failed to write file")
	}

	// Validate magic bytes
	file.Seek(0, 0)
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	mime := mimetype.Detect(buf[:n])
	file.Seek(0, 0)

	allowedMIMEs := map[string]bool{
		"application/pdf": true,
		"image/jpeg":      true,
		"image/png":       true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}

	if !allowedMIMEs[mime.String()] {
		os.Remove(filePath) // clean up invalid file
		return nil, errors.NewBadRequest("file type not allowed: " + mime.String())
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return s.SaveMetadata(ctx, companyID, userID, header.Filename, storageKey, mimeType, header.Size, fileType, checksum)
}

func (s *FileService) SaveMetadata(
	ctx context.Context,
	companyID, ownerUserID, originalName, storageKey, mimeType string,
	sizeBytes int64,
	fileType string,
	checksum string,
) (*models.File, error) {
	f := &models.File{
		ID:           uuid.NewString(),
		CompanyID:    sql.NullString{String: companyID, Valid: companyID != ""},
		OwnerUserID:  sql.NullString{String: ownerUserID, Valid: ownerUserID != ""},
		OriginalName: originalName,
		StorageKey:   storageKey,
		MimeType:     mimeType,
		SizeBytes:    sizeBytes,
		FileType:     fileType,
		Checksum:     sql.NullString{String: checksum, Valid: checksum != ""},
		CreatedAt:    time.Now(),
	}

	created, err := s.fileRepo.Create(ctx, f)
	if err != nil {
		return nil, errors.NewInternal("failed to save file metadata")
	}
	return created, nil
}

func (s *FileService) GetSignedURL(
	ctx context.Context,
	fileID, requestingUserID, companyID string,
	isRecruiter bool,
	isAdmin bool,
) (url string, expiresAt time.Time, err error) {
	var file *models.File

	switch {
	case isAdmin:
		// Admin can access any file (e.g. recruiter verification documents,
		// which are owner-uploaded and not scoped to a company).
		file, err = s.fileRepo.GetByID(ctx, fileID)
	case isRecruiter:
		file, err = s.fileRepo.GetByIDAndCompany(ctx, fileID, companyID)
	default:
		file, err = s.fileRepo.GetByIDAndOwner(ctx, fileID, requestingUserID)
	}

	if err != nil || file == nil {
		return "", time.Time{}, errors.NewNotFound("file not found or access denied")
	}

	expiresAt = time.Now().Add(15 * time.Minute)
	return s.PublicURL(file.StorageKey), expiresAt, nil
}
