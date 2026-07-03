package service

import (
	"context"

	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/repository"
)

type ReportService struct {
	repo          *repository.ReportRepository
	transcriptSvc *TranscriptService
	aiSvc         *AIService
}

func NewReportService(repo *repository.ReportRepository, transcriptSvc *TranscriptService, aiSvc *AIService) *ReportService {
	return &ReportService{
		repo:          repo,
		transcriptSvc: transcriptSvc,
		aiSvc:         aiSvc,
	}
}

// GetOrGenerateReport fetches the report. If not exists, it generates it using AI.
func (s *ReportService) GetOrGenerateReport(ctx context.Context, interviewID, companyID string) (*models.InterviewReport, error) {
	report, err := s.repo.GetByInterviewID(ctx, interviewID)
	if err == nil && report != nil {
		return report, nil
	}

	// Not found, generate it
	transcripts, err := s.transcriptSvc.ListTranscripts(ctx, interviewID, companyID)
	if err != nil {
		return nil, errors.NewInternal("failed to fetch transcripts for report")
	}

	if len(transcripts) == 0 {
		return nil, errors.NewBadRequest("no transcripts available to generate report")
	}

	// Generate via AI
	newReport, err := s.aiSvc.GenerateInterviewReport(ctx, interviewID, transcripts)
	if err != nil {
		return nil, errors.NewInternal("failed to generate AI report: " + err.Error())
	}

	// Save to DB
	err = s.repo.Create(ctx, newReport)
	if err != nil {
		return nil, errors.NewInternal("failed to save generated report")
	}

	return newReport, nil
}

func (s *ReportService) UpdateDecision(ctx context.Context, interviewID, companyID, decision, comment string) error {
	// First ensure interview belongs to company (authorization)
	_, err := s.transcriptSvc.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return errors.NewNotFound("interview not found or access denied")
	}

	err = s.repo.UpdateDecision(ctx, interviewID, decision, comment)
	if err != nil {
		return errors.NewInternal("failed to update decision")
	}
	return nil
}
