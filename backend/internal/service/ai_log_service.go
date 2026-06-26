package service

import (
	"context"
	"log"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
)

type AILogService struct {
	repo    *repository.AILogRepository
	logChan chan *models.AIRequestLog
}

// NewAILogService starts a background worker that listens for AI logs and inserts them.
func NewAILogService(repo *repository.AILogRepository) *AILogService {
	s := &AILogService{
		repo:    repo,
		logChan: make(chan *models.AIRequestLog, 1000), // Buffer size of 1000
	}
	go s.worker()
	return s
}

// worker processes the queue and saves logs to the DB asynchronously.
func (s *AILogService) worker() {
	for logEntry := range s.logChan {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := s.repo.Create(ctx, logEntry)
		if err != nil {
			log.Printf("[AILogService] failed to save AI log: %v", err)
		}
		cancel()
	}
}

// LogAsync sends the log entry to the background worker. It doesn't block the caller.
func (s *AILogService) LogAsync(logEntry *models.AIRequestLog) {
	select {
	case s.logChan <- logEntry:
		// Successfully queued
	default:
		// Queue is full, drop the log to prevent blocking the API
		log.Printf("[AILogService] WARNING: log queue is full, dropping AI log for template %s", logEntry.TemplateID.String)
	}
}
