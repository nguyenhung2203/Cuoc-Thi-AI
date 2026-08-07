package realtime

import (
	"context"
	"database/sql"
	"log"

	"backend/internal/models"
	"backend/internal/repository"

	"github.com/google/uuid"
)

// AuditLogger handles async writing of audit logs
type AuditLogger struct {
	logChan chan *models.AuditLog
	repo    *repository.AuditRepository
}

// NewAuditLogger initializes the logger and starts the background worker
func NewAuditLogger(repo *repository.AuditRepository) *AuditLogger {
	logger := &AuditLogger{
		logChan: make(chan *models.AuditLog, 1000),
		repo:    repo,
	}
	go logger.start()
	return logger
}

// start runs in a goroutine to process logs from the channel
func (l *AuditLogger) start() {
	for al := range l.logChan {
		err := l.repo.Insert(context.Background(), al)
		if err != nil {
			log.Printf("[audit] Failed to insert audit log: %v", err)
		}
	}
}

// LogEvent queues a new audit log event to be processed asynchronously
func (l *AuditLogger) LogEvent(action, actorID, actorRole, resourceType, resourceID, companyID, ip string, _ map[string]interface{}) {
	al := &models.AuditLog{
		ID:           uuid.New().String(),
		Action:       action,
		ResourceType: resourceType,
	}
	nullStr := func(s string) sql.NullString {
		if s == "" {
			return sql.NullString{}
		}
		return sql.NullString{String: s, Valid: true}
	}
	if actorID != "" {
		al.ActorUserID = nullStr(actorID)
	}
	if actorRole != "" {
		al.ActorRole = nullStr(actorRole)
	}
	if resourceID != "" {
		al.ResourceID = nullStr(resourceID)
	}
	if companyID != "" {
		al.CompanyID = nullStr(companyID)
	}
	if ip != "" {
		al.IPAddress = nullStr(ip)
	}

	select {
	case l.logChan <- al:
		// Successfully queued
	default:
		log.Printf("[audit] ERROR: Audit log channel is full, dropping event %s", action)
	}
}
