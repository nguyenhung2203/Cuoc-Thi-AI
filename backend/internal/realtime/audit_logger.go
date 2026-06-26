package realtime

import (
	"log"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
)

// AuditLogger handles async writing of audit logs
type AuditLogger struct {
	logChan chan *models.AuditLog
	repo    *repository.AuditRepository
}

// NewAuditLogger initializes the logger and starts the background worker
func NewAuditLogger(repo *repository.AuditRepository) *AuditLogger {
	logger := &AuditLogger{
		logChan: make(chan *models.AuditLog, 1000), // Buffer to avoid blocking
		repo:    repo,
	}
	go logger.start()
	return logger
}

// start runs in a goroutine to process logs from the channel
func (l *AuditLogger) start() {
	for al := range l.logChan {
		err := l.repo.Insert(al)
		if err != nil {
			log.Printf("[audit] Failed to insert audit log: %v", err)
		}
	}
}

// LogEvent queues a new audit log event to be processed asynchronously
func (l *AuditLogger) LogEvent(action, actorID, actorType, resourceType, resourceID, companyID, ip string, metadata map[string]interface{}) {
	al := &models.AuditLog{
		Action:       action,
		ActorID:      actorID,
		ActorType:    actorType,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		CompanyID:    companyID,
		Metadata:     metadata,
		IPAddress:    ip,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	select {
	case l.logChan <- al:
		// Successfully queued
	default:
		// Channel is full, log error and drop to avoid blocking realtime
		log.Printf("[audit] ERROR: Audit log channel is full, dropping event %s", action)
	}
}
