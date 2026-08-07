package service

import (
	"context"
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
	strPtr := func(s string) *string { return &s }
	if in.CompanyID != "" {
		al.CompanyID = strPtr(in.CompanyID)
	}
	if in.ActorUserID != "" {
		al.ActorUserID = strPtr(in.ActorUserID)
	}
	if in.ActorRole != "" {
		al.ActorRole = strPtr(in.ActorRole)
	}
	if in.ResourceID != "" {
		al.ResourceID = strPtr(in.ResourceID)
	}
	if in.IPAddress != "" {
		al.IPAddress = strPtr(in.IPAddress)
	}
	if in.UserAgent != "" {
		al.UserAgent = strPtr(in.UserAgent)
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
	if offset < 0 {
		offset = 0
	}
	return s.auditRepo.ListByCompany(ctx, companyID, resourceType, actorUserID, pageSize, offset)
}

func (s *AuditService) ListAll(ctx context.Context, page, pageSize int) ([]models.AuditLog, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	return s.auditRepo.ListAll(ctx, pageSize, offset)
}
