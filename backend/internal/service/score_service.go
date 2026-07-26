package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"backend/internal/models"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/utils"
	"backend/internal/repository"
)

type ScoreService struct {
	scoreRepo      repository.ScoreRepository
	transcriptRepo *repository.TranscriptRepository
	rubricRepo     repository.RubricRepository
	interviewRepo  *repository.InterviewRepository
	orchestrator   *AIOrchestratorService
}

func NewScoreService(
	scoreRepo repository.ScoreRepository,
	transcriptRepo *repository.TranscriptRepository,
	rubricRepo repository.RubricRepository,
	interviewRepo *repository.InterviewRepository,
	orchestrator *AIOrchestratorService,
) *ScoreService {
	return &ScoreService{
		scoreRepo:      scoreRepo,
		transcriptRepo: transcriptRepo,
		rubricRepo:     rubricRepo,
		interviewRepo:  interviewRepo,
		orchestrator:   orchestrator,
	}
}

func (s *ScoreService) ScoreAnswer(ctx context.Context, companyID, interviewID string, transcriptIDs []string, criterionIDs []string) ([]models.InterviewScore, error) {
	if len(transcriptIDs) == 0 || len(criterionIDs) == 0 {
		return nil, apierrors.NewValidation("score-answer", []string{"transcript_ids and criterion_ids are required"})
	}

	// 0. Verify interview belongs to company (Prevent IDOR)
	interview, err := s.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return nil, err
	}
	if interview == nil {
		return nil, apierrors.NewNotFound("interview")
	}

	// 1. Fetch transcripts
	allTranscripts, err := s.transcriptRepo.ListByInterview(ctx, interviewID)
	if err != nil {
		return nil, err
	}

	// Filter transcripts to build the text
	var transcriptSegments []string
	for _, id := range transcriptIDs {
		for _, t := range allTranscripts {
			if t.ID == id {
				content := t.Content
				if t.EditedContent != nil && *t.EditedContent != "" {
					content = *t.EditedContent
				}
				transcriptSegments = append(transcriptSegments, fmt.Sprintf("%s: %s", t.SpeakerType, content))
				break
			}
		}
	}

	if len(transcriptSegments) == 0 {
		return nil, apierrors.NewValidation("transcript_ids", []string{"no matching transcripts found in this interview"})
	}

	transcriptText := strings.Join(transcriptSegments, "\n")

	// Prevent DoS
	transcriptText = utils.TruncateText(transcriptText, 30000)

	// 2. Fetch all criteria
	criteria, err := s.rubricRepo.GetCriteriaByIDs(ctx, criterionIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch criteria: %w", err)
	}
	if len(criteria) == 0 {
		return nil, apierrors.NewValidation("criterion_ids", []string{"none of the provided criterion_ids were found"})
	}

	var scores []models.InterviewScore

	// 3. Score each criterion
	for _, criterion := range criteria {
		// Prepare scoring guide
		var scoringGuideStr string
		if len(criterion.ScoringGuide) > 0 {
			scoringGuideStr = string(criterion.ScoringGuide)
		} else {
			scoringGuideStr = "Not provided"
		}

		var criterionDesc string
		if criterion.Description.Valid {
			criterionDesc = criterion.Description.String
		} else {
			criterionDesc = criterion.Name
		}

		result, err := s.orchestrator.ScoreAnswer(
			ctx,
			companyID,
			transcriptText,
			criterion.Name,
			criterionDesc,
			criterion.MinScore,
			criterion.MaxScore,
			scoringGuideStr,
		)

		status := "scored"
		var scoreVal sql.NullFloat64
		var weightedScoreVal sql.NullFloat64

		if err != nil {
			status = "error"
			if appErr, ok := err.(*apierrors.AppError); ok && appErr.Code == apierrors.VALIDATION_ERROR {
				status = "insufficient_evidence"
			}
			if result == nil {
				result = &ScoreResult{}
			}
			scoreVal = sql.NullFloat64{Valid: false}
			weightedScoreVal = sql.NullFloat64{Valid: false}
		} else {
			weightedScore := (result.Score / float64(criterion.MaxScore)) * criterion.Weight
			scoreVal = sql.NullFloat64{Float64: result.Score, Valid: true}
			weightedScoreVal = sql.NullFloat64{Float64: weightedScore, Valid: true}
		}

		scoreRecord := models.InterviewScore{
			InterviewID:       interviewID,
			RubricCriterionID: sql.NullString{String: criterion.ID, Valid: true},
			CriterionName:     criterion.Name,
			Score:             scoreVal,
			MaxScore:          float64(criterion.MaxScore),
			Weight:            criterion.Weight,
			WeightedScore:     weightedScoreVal,
			Evidence:          sql.NullString{String: result.Evidence, Valid: true},
			AIComment:         sql.NullString{String: result.AIComment, Valid: true},
			Confidence:        sql.NullFloat64{Float64: result.Confidence, Valid: true},
			Status:            status,
			ScoredBy:          "ai",
		}

		// 4. Save to DB
		if err := s.scoreRepo.SaveScore(ctx, &scoreRecord); err != nil {
			return nil, fmt.Errorf("failed to save score for criterion %s: %w", criterion.Name, err)
		}

		scores = append(scores, scoreRecord)
	}

	return scores, nil
}
