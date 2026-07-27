package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"backend/internal/ai"
	"backend/internal/dto/response"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/utils"
	"backend/internal/repository"
)

// SuggestionService handles AI follow-up suggestion generation.
type SuggestionService struct {
	orchestrator    *AIOrchestratorService
	suggestAnalyzer *ai.SuggestFollowUpAnalyzer
	transcriptRepo  *repository.TranscriptRepository
	interviewRepo   *repository.InterviewRepository
	jobRepo         *repository.JobRepository
}

func NewSuggestionService(
	orchestrator *AIOrchestratorService,
	transcriptRepo *repository.TranscriptRepository,
	interviewRepo *repository.InterviewRepository,
	jobRepo *repository.JobRepository,
) *SuggestionService {
	return &SuggestionService{
		orchestrator:    orchestrator,
		suggestAnalyzer: ai.NewSuggestFollowUpAnalyzer(orchestrator),
		transcriptRepo:  transcriptRepo,
		interviewRepo:   interviewRepo,
		jobRepo:         jobRepo,
	}
}

// SuggestFollowUp generates an AI follow-up question based on recent transcript context.
func (s *SuggestionService) SuggestFollowUp(
	ctx context.Context,
	companyID, interviewID, lastTranscriptID, focus string,
) (*response.SuggestFollowUpResponse, error) {
	// 1. Verify interview belongs to company (Prevent IDOR)
	interview, err := s.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return nil, err
	}
	if interview == nil {
		return nil, apierrors.NewNotFound("interview")
	}

	// 2. Fetch recent transcripts (latest 20 segments)
	transcripts, err := s.transcriptRepo.ListByInterview(ctx, interviewID)
	if err != nil {
		return nil, apierrors.NewInternal("failed to fetch transcripts")
	}

	// Build a transcript window: most recent N segments
	windowSize := 10
	if len(transcripts) > windowSize {
		transcripts = transcripts[len(transcripts)-windowSize:]
	}

	var transcriptLines []string
	var windowItems []response.TranscriptItem
	for _, t := range transcripts {
		content := t.Content
		if t.EditedContent != nil && *t.EditedContent != "" {
			content = *t.EditedContent
		}
		transcriptLines = append(transcriptLines, fmt.Sprintf("%s: %s", t.SpeakerType, content))
		windowItems = append(windowItems, response.TranscriptItem{
			ID:      t.ID,
			Speaker: t.SpeakerType,
			Content: content,
			IsFinal: t.IsFinal,
		})
	}

	transcriptWindow := strings.Join(transcriptLines, "\n")
	transcriptWindow = utils.TruncateText(transcriptWindow, 15000)

	// 3. Fetch job context (title + description) for better suggestions
	jobContext := ""
	if interview.JobID.Valid {
		job, err := s.jobRepo.GetByID(ctx, companyID, interview.JobID.String)
		if err == nil && job != nil {
			jobContext = utils.TruncateText(fmt.Sprintf("%s: %s", job.Title, job.Description), 5000)
		}
	}

	// 4. Call AI orchestrator
	result, err := s.suggestAnalyzer.SuggestFollowUp(ctx, transcriptWindow, jobContext, focus, companyID)
	if err != nil {
		return nil, fmt.Errorf("AI suggestion failed: %w", err)
	}

	if result.InsufficientData {
		return nil, apierrors.NewValidation("insufficient data", []string{"AI needs more conversation context to suggest a follow-up"})
	}

	now := time.Now()
	expiresAt := now.Add(5 * time.Minute) // suggestions expire after 5 min

	resp := &response.SuggestFollowUpResponse{
		SuggestedQuestion: result.SuggestedQuestion,
		Reason:            result.Reason,
		TargetSkill:       result.TargetSkill,
		Priority:          result.Priority,
		Confidence:        result.Confidence,
		CreatedAt:         now,
		ExpiresAt:         &expiresAt,
		TranscriptWindow:  windowItems,
	}

	_ = uuid.New() // seeded for future ai_suggestions table insert if needed

	return resp, nil
}
