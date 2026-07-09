package service

import (
	"context"
	"database/sql"
	"encoding/json"

	"backend/internal/models"
	"backend/internal/repository"

	"github.com/google/uuid"
)

type AuditService struct {
	auditRepo *repository.AuditRepository
}

func NewAuditService(auditRepo *repository.AuditRepository) *AuditService {
	return &AuditService{auditRepo: auditRepo}
}

type AuditLogInput struct {
	CompanyID    string
	ActorUserID  string
	ActorRole    string
	Action       string
	ResourceType string
	ResourceID   string
	BeforeData   interface{}
	AfterData    interface{}
	IPAddress    string
	UserAgent    string
}

// LogAction records an audit log entry with before/after snapshots serialized to JSON.
func (s *AuditService) LogAction(ctx context.Context, in AuditLogInput) error {
	al := &models.AuditLog{
		ID:           uuid.New().String(),
		Action:       in.Action,
		ResourceType: in.ResourceType,
	}
	if in.CompanyID != "" {
		al.CompanyID = sql.NullString{String: in.CompanyID, Valid: true}
	}
	if in.ActorUserID != "" {
		al.ActorUserID = sql.NullString{String: in.ActorUserID, Valid: true}
	}
	if in.ActorRole != "" {
		al.ActorRole = sql.NullString{String: in.ActorRole, Valid: true}
	}
	if in.ResourceID != "" {
		al.ResourceID = sql.NullString{String: in.ResourceID, Valid: true}
	}
	if in.IPAddress != "" {
		al.IPAddress = sql.NullString{String: in.IPAddress, Valid: true}
	}
	if in.UserAgent != "" {
		al.UserAgent = sql.NullString{String: in.UserAgent, Valid: true}
	}
	if in.BeforeData != nil {
		b, _ := json.Marshal(in.BeforeData)
		al.BeforeJSON = models.JSONB(b)
	}
	if in.AfterData != nil {
		b, _ := json.Marshal(in.AfterData)
		al.AfterJSON = models.JSONB(b)
	}

	return s.auditRepo.Insert(ctx, al)
}

// ListByCompany returns audit logs for a company with optional filters.
func (s *AuditService) ListByCompany(ctx context.Context, companyID, resourceType, actorUserID string, page, pageSize int) ([]models.AuditLog, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.auditRepo.ListByCompany(ctx, companyID, resourceType, actorUserID, pageSize, offset)
}
