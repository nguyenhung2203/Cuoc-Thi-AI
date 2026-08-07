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

// LogAsync queues an AI usage record. If the buffer is saturated, write it
// synchronously instead of silently losing billing/audit data.
func (s *AILogService) LogAsync(logEntry *models.AIRequestLog) {
	select {
	case s.logChan <- logEntry:
	default:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.repo.Create(ctx, logEntry); err != nil {
			log.Printf("[AILogService] failed to save saturated-queue AI log: %v", err)
		}
	}
}

// List returns AI request logs matching the filter (admin debugging / cost audit).
func (s *AILogService) List(ctx context.Context, f repository.AILogFilter) ([]models.AIRequestLog, error) {
	return s.repo.List(ctx, f)
}
