package service

import (
	"context"
	"strings"
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
	return &TranscriptService{repo: repo, interviewRepo: interviewRepo}
}

type PushTranscriptRequest struct {
	ParticipantID string  `json:"participant_id"`
	SpeakerType   string  `json:"speaker_type"`
	SpeakerName   string  `json:"speaker_name"`
	Content       string  `json:"content"`
	Language      string  `json:"language"`
	StartTimeMs   int64   `json:"start_time_ms"`
	EndTimeMs     int64   `json:"end_time_ms"`
	Confidence    float64 `json:"confidence"`
	Source        string  `json:"source"`
	IsFinal       *bool   `json:"is_final"`
}

func (s *TranscriptService) PushTranscript(ctx context.Context, interviewID, companyID string, req PushTranscriptRequest) (*models.InterviewTranscript, error) {
	if strings.TrimSpace(req.Content) == "" {
		return nil, errors.NewValidation("content", []string{"content is required"})
	}
	if !validSpeakerType(req.SpeakerType) {
		return nil, errors.NewValidation("speaker_type", []string{"must be recruiter, candidate, ai, or system"})
	}
	if req.Source == "" {
		req.Source = "manual"
	}
	if !validTranscriptSource(req.Source) {
		return nil, errors.NewValidation("source", []string{"must be audio, chat, manual, or ai"})
	}

	_, err := s.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return nil, errors.NewNotFound("interview not found or access denied")
	}

	now := time.Now().UTC()
	isFinal := true
	if req.IsFinal != nil {
		isFinal = *req.IsFinal
	}
	transcript := &models.InterviewTranscript{
		ID:          uuid.NewString(),
		InterviewID: interviewID,
		SpeakerType: req.SpeakerType,
		Content:     req.Content,
		Source:      req.Source,
		IsFinal:     isFinal,
		CreatedAt:   now,
	}
	if req.ParticipantID != "" {
		transcript.ParticipantID = &req.ParticipantID
	}
	if req.SpeakerName != "" {
		transcript.SpeakerName = &req.SpeakerName
	}
	if req.Language != "" {
		transcript.Language = &req.Language
	}
	if req.StartTimeMs != 0 {
		transcript.StartTimeMs = &req.StartTimeMs
	}
	if req.EndTimeMs != 0 {
		transcript.EndTimeMs = &req.EndTimeMs
	}
	if req.Confidence != 0 {
		transcript.Confidence = &req.Confidence
	}

	if err := s.repo.Create(ctx, transcript); err != nil {
		return nil, err
	}
	return transcript, nil
}

func validSpeakerType(value string) bool {
	switch value {
	case "recruiter", "candidate", "ai", "system":
		return true
	default:
		return false
	}
}

func validTranscriptSource(value string) bool {
	switch value {
	case "audio", "chat", "manual", "ai":
		return true
	default:
		return false
	}
}

func (s *TranscriptService) ListTranscripts(ctx context.Context, interviewID, companyID string) ([]models.InterviewTranscript, error) {
	_, err := s.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return nil, errors.NewNotFound("interview not found or access denied")
	}
	return s.repo.ListByInterview(ctx, interviewID)
}

func (s *TranscriptService) EditTranscript(ctx context.Context, transcriptID, interviewID, companyID, editedContent, editedBy string) error {
	if strings.TrimSpace(editedContent) == "" {
		return errors.NewValidation("edited_content", []string{"edited_content is required"})
	}
	_, err := s.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return errors.NewNotFound("interview not found or access denied")
	}
	return s.repo.UpdateEditedContent(ctx, transcriptID, interviewID, editedContent, editedBy)
}
