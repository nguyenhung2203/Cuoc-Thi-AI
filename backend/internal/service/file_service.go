package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"io"
	"mime/multipart"
	"regexp"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"

	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/repository"
	"backend/internal/storage"
)

// allowedMIMEs is the upload whitelist, checked against the DETECTED type
// (magic bytes), never the client-supplied Content-Type header.
var allowedMIMEs = map[string]bool{
	"application/pdf":    true,
	"image/jpeg":         true,
	"image/png":          true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

// mimeExt maps detected MIME types to the sanitized storage-key extension.
// Extensions never come from the client filename — that produced keys like
// "cv.p df" and could smuggle unexpected types.
var mimeExt = map[string]string{
	"application/pdf":    ".pdf",
	"image/jpeg":         ".jpg",
	"image/png":          ".png",
	"application/msword": ".doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
}

var extPattern = regexp.MustCompile(`^\.[a-z0-9]{1,5}$`)

// sanitizeExt returns ext lowercased when it is a plain dot-extension
// ([a-z0-9]{1,5}); anything else becomes empty.
func sanitizeExt(ext string) string {
	ext = strings.ToLower(ext)
	if extPattern.MatchString(ext) {
		return ext
	}
	return ""
}

// isAllowedMIME reports whether a detected MIME type may be stored.
func isAllowedMIME(mime string) bool { return allowedMIMEs[mime] }

type FileService struct {
	fileRepo *repository.FileRepository
	baseURL  string
	store    *storage.LocalStore
	signer   *storage.Signer
	maxBytes int64
}

func NewFileService(fileRepo *repository.FileRepository, baseURL string, store *storage.LocalStore, signer *storage.Signer, maxBytes int64) *FileService {
	return &FileService{fileRepo: fileRepo, baseURL: baseURL, store: store, signer: signer, maxBytes: maxBytes}
}

// MaxUploadBytes exposes the configured upload cap to handlers.
func (s *FileService) MaxUploadBytes() int64 {
	if s.maxBytes <= 0 {
		return 10 << 20
	}
	return s.maxBytes
}

// SignedURL builds a time-limited, HMAC-signed download URL for a stored file
// key. This replaced PublicURL: an unsigned /uploads/ URL must never escape.
func (s *FileService) SignedURL(storageKey string) (string, time.Time) {
	return s.signer.SignURL(s.baseURL, storageKey)
}

func (s *FileService) ProcessUpload(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID, companyID, fileType string) (*models.File, error) {
	// Defense-in-depth size guard; the handler-level MaxBytesReader is primary.
	if max := s.MaxUploadBytes(); header.Size > max {
		return nil, errors.NewPayloadTooLarge("file exceeds the maximum allowed size")
	}

	// 1. Detect the real MIME type from magic bytes BEFORE any disk write, so
	// rejected uploads never touch the filesystem.
	head := make([]byte, 512)
	n, _ := io.ReadFull(file, head)
	mime := mimetype.Detect(head[:n])
	if !isAllowedMIME(mime.String()) {
		return nil, errors.NewValidation("file", []string{
			"file type not allowed: " + mime.String(),
			"supported formats: PDF, JPEG, PNG, DOC, DOCX",
		})
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, errors.NewInternal("failed to read uploaded file")
	}

	// 2. Checksum (integrity metadata only — per-file storage, no dedupe:
	// sharing one blob across tenants by checksum let company A mint valid
	// signed URLs to company B's stored object).
	hasher := sha256.New()
	written, err := io.Copy(hasher, file)
	if err != nil {
		return nil, errors.NewInternal("failed to hash file")
	}
	checksum := hex.EncodeToString(hasher.Sum(nil))
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, errors.NewInternal("failed to read uploaded file")
	}

	// 3. Write to a .part key, then atomically publish via rename so a
	// half-written file can never be served.
	ext := sanitizeExt(mimeExt[mime.String()])
	storageKey := uuid.NewString() + ext
	tmpKey := storageKey + ".part"

	dst, err := s.store.Create(tmpKey)
	if err != nil {
		return nil, errors.NewInternal("failed to save file")
	}
	if _, err := io.Copy(dst, file); err != nil {
		dst.Close()
		_ = s.store.Remove(tmpKey)
		return nil, errors.NewInternal("failed to write file")
	}
	if err := dst.Sync(); err != nil {
		dst.Close()
		_ = s.store.Remove(tmpKey)
		return nil, errors.NewInternal("failed to write file")
	}
	if err := dst.Close(); err != nil {
		_ = s.store.Remove(tmpKey)
		return nil, errors.NewInternal("failed to write file")
	}
	if err := s.store.Rename(tmpKey, storageKey); err != nil {
		_ = s.store.Remove(tmpKey)
		return nil, errors.NewInternal("failed to save file")
	}

	// Persist the DETECTED type and the actual byte count — both are what the
	// download handler will serve, not what the client claimed.
	return s.SaveMetadata(ctx, companyID, userID, header.Filename, storageKey, mime.String(), written, fileType, checksum)
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

	url, expiresAt = s.SignedURL(file.StorageKey)
	return url, expiresAt, nil
}
