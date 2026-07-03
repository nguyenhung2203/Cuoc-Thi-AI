package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/repository"
)

type TranscriptService struct {
	repo          *repository.TranscriptRepository
	interviewRepo *repository.InterviewRepository
}

func NewTranscriptService(repo *repository.TranscriptRepository, interviewRepo *repository.InterviewRepository) *TranscriptService {
	return &TranscriptService{
		repo:          repo,
		interviewRepo: interviewRepo,
	}
}

type PushTranscriptRequest struct {
	SpeakerID   string  `json:"speaker_id"`
	SpeakerRole string  `json:"speaker_role"`
	StartTime   float64 `json:"start_time"`
	EndTime     float64 `json:"end_time"`
	Content     string  `json:"content"`
	IsFinal     bool    `json:"is_final"`
	Language    string  `json:"language"`
}

func (s *TranscriptService) PushTranscript(ctx context.Context, interviewID, companyID string, req PushTranscriptRequest) (*models.InterviewTranscript, error) {
	_, err := s.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return nil, errors.NewNotFound("interview not found or access denied")
	}

	t := &models.InterviewTranscript{
		ID:          uuid.NewString(),
		InterviewID: interviewID,
		SpeakerID:   req.SpeakerID,
		SpeakerRole: req.SpeakerRole,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Content:     req.Content,
		IsFinal:     req.IsFinal,
		Language:    req.Language,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.repo.Create(ctx, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TranscriptService) ListTranscripts(ctx context.Context, interviewID, companyID string) ([]models.InterviewTranscript, error) {
	_, err := s.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return nil, errors.NewNotFound("interview not found or access denied")
	}
	return s.repo.ListByInterview(ctx, interviewID)
}

func (s *TranscriptService) EditTranscript(ctx context.Context, transcriptID, interviewID, companyID, editedContent, editedBy string) error {
	_, err := s.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return errors.NewNotFound("interview not found or access denied")
	}
	return s.repo.UpdateEditedContent(ctx, transcriptID, interviewID, editedContent, editedBy)
}
