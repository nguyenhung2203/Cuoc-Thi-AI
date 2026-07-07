package service

import (
	"context"


	"backend/internal/models"
	"backend/internal/pkg/errors"

	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"backend/internal/models"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/logger"
	"backend/internal/pkg/utils"

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

	reportRepo     repository.ReportRepository
	transcriptRepo *repository.TranscriptRepository
	scoreRepo      repository.ScoreRepository
	jobRepo        *repository.JobRepository
	interviewRepo  *repository.InterviewRepository
	notifRepo      *repository.NotificationRepository
	orchestrator   *AIOrchestratorService
}

func NewReportService(
	reportRepo repository.ReportRepository,
	transcriptRepo *repository.TranscriptRepository,
	scoreRepo repository.ScoreRepository,
	jobRepo *repository.JobRepository,
	interviewRepo *repository.InterviewRepository,
	notifRepo *repository.NotificationRepository,
	orchestrator *AIOrchestratorService,
) *ReportService {
	return &ReportService{
		reportRepo:     reportRepo,
		transcriptRepo: transcriptRepo,
		scoreRepo:      scoreRepo,
		jobRepo:        jobRepo,
		interviewRepo:  interviewRepo,
		notifRepo:      notifRepo,
		orchestrator:   orchestrator,
	}
}

// generateReportSync does the actual work of generating a report synchronously.
// It does NOT manage report_status transitions — caller must handle that.
func (s *ReportService) generateReportSync(ctx context.Context, companyID, interviewID, jobID, generatedBy string) error {
	// 1. Fetch transcript
	transcripts, err := s.transcriptRepo.ListByInterview(ctx, interviewID)
	if err != nil {
		return fmt.Errorf("failed to fetch transcripts: %w", err)
	}

	var transcriptSegments []string
	for _, t := range transcripts {
		content := t.Content
		if t.EditedContent.Valid && t.EditedContent.String != "" {
			content = t.EditedContent.String
		}
		transcriptSegments = append(transcriptSegments, fmt.Sprintf("%s: %s", t.SpeakerRole, content))
	}
	transcriptText := strings.Join(transcriptSegments, "\n")
	if transcriptText == "" {
		return apierrors.NewValidation("interview_id", []string{"interview has no transcript"})
	}
	// Prevent DoS
	transcriptText = utils.TruncateText(transcriptText, 50000)

	// 2. Fetch scores
	scores, err := s.scoreRepo.GetScoresByInterviewID(ctx, interviewID)
	if err != nil {
		return fmt.Errorf("failed to fetch scores: %w", err)
	}

	scoresBytes, marshalErr := json.Marshal(scores)
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal scores: %w", marshalErr)
	}
	scoresJSON := string(scoresBytes)

	// 3. Fetch job requirements
	job, err := s.jobRepo.GetByID(ctx, companyID, jobID)
	if err != nil {
		return fmt.Errorf("failed to fetch job: %w", err)
	}

	var jobRequirements string
	if job.Requirements.Valid {
		jobRequirements = job.Requirements.String
	} else {
		jobRequirements = job.Description
	}

	// 4. Call AI
	aiReport, err := s.orchestrator.GenerateReport(ctx, companyID, jobRequirements, transcriptText, scoresJSON)
	if err != nil {
		return fmt.Errorf("AI report generation failed: %w", err)
	}

	// 5. Serialize results
	strengthsBytes, _ := json.Marshal(aiReport.Strengths)
	weaknessesBytes, _ := json.Marshal(aiReport.Weaknesses)
	risksBytes, _ := json.Marshal(aiReport.Risks)
	evidenceBytes, _ := json.Marshal(aiReport.EvidenceJSON)
	fullJSONBytes, _ := json.Marshal(aiReport)

	dbReport := &models.InterviewReport{
		InterviewID:    interviewID,
		Summary:        aiReport.Summary,
		FinalScore:     sql.NullFloat64{Float64: aiReport.FinalScore, Valid: true},
		Recommendation: aiReport.Recommendation,
		Strengths:      models.JSONB(strengthsBytes),
		Weaknesses:     models.JSONB(weaknessesBytes),
		Risks:          models.JSONB(risksBytes),
		EvidenceJSON:   models.JSONB(evidenceBytes),
		ReportJSON:     models.JSONB(fullJSONBytes),
		GeneratedBy:    generatedBy,
	}
	if aiReport.AIReasoningSummary != "" {
		dbReport.AIReasoningSummary = sql.NullString{String: aiReport.AIReasoningSummary, Valid: true}
	}

	// 6. Persist
	return s.reportRepo.UpsertReport(ctx, dbReport)
}

// GenerateReport pulls data and sends to AI to generate the report asynchronously.
func (s *ReportService) GenerateReport(ctx context.Context, companyID, interviewID, jobID string, generatedBy string) (*models.InterviewReport, error) {
	// 0. Verify interview belongs to company (Prevent IDOR)
	interview, err := s.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return nil, err
	}
	if interview == nil {
		return nil, apierrors.NewNotFound("interview")
	}

	// Atomic: only proceed if not already generating
	ok, err := s.interviewRepo.UpdateReportStatusIf(ctx, interviewID, "generating", "")
	if err != nil {
		return nil, fmt.Errorf("failed to lock report generation: %w", err)
	}
	if !ok {
		// Try "pending" too, in case it was never set (first time)
		ok, err = s.interviewRepo.UpdateReportStatusIf(ctx, interviewID, "generating", "pending")
		if err != nil {
			return nil, fmt.Errorf("failed to lock report generation: %w", err)
		}
	}
	if !ok {
		return nil, apierrors.NewValidation("interview_id", []string{"report is currently generating"})
	}

	recruiterID := interview.RecruiterID.String

	// Spin up goroutine to process async with timeout
	go func(compID, intID, jobIDStr, genBy, recID string) {
		// Timeout để tránh goroutine leak khi AI treo
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*60e9) // 10 minutes
		defer cancel()

		// Protect from panics
		defer func() {
			if r := recover(); r != nil {
				_ = s.interviewRepo.UpdateReportStatus(bgCtx, intID, "failed")
				s.notifyReportFailed(bgCtx, intID, recID, "panic")
			}
		}()

		err := s.generateReportSync(bgCtx, compID, intID, jobIDStr, genBy)
		if err != nil {
			if !strings.Contains(err.Error(), "context deadline exceeded") {
				_ = s.interviewRepo.UpdateReportStatus(bgCtx, intID, "failed")
				s.notifyReportFailed(bgCtx, intID, recID, err.Error())
			} else {
				_ = s.interviewRepo.UpdateReportStatus(bgCtx, intID, "failed")
				s.notifyReportFailed(bgCtx, intID, recID, "timeout")
			}
			return
		}

		_ = s.interviewRepo.UpdateReportStatus(bgCtx, intID, "ready")

		// Send notification
		if recID != "" {
			_ = s.notifRepo.Create(bgCtx, &models.Notification{
				UserID:  recID,
				Title:   "Báo cáo AI đã hoàn thành",
				Message: "Báo cáo phỏng vấn đã được AI tổng hợp xong. Vui lòng xem kết quả.",
				Type:    "report_ready",
				Link:    "/interviews/" + intID + "/report",
			})
		}
	}(companyID, interviewID, jobID, generatedBy, recruiterID)

	// Return a stub report indicating generating status
	return &models.InterviewReport{
		InterviewID: interviewID,
		Summary:     "Báo cáo đang được xử lý bởi AI...",
	}, nil
}

func (s *ReportService) notifyReportFailed(ctx context.Context, interviewID, recruiterID, reason string) {
	if recruiterID == "" {
		return
	}
	if err := s.notifRepo.Create(ctx, &models.Notification{
		UserID:  recruiterID,
		Title:   "Báo cáo AI thất bại",
		Message: "Báo cáo phỏng vấn không thể tạo do lỗi AI. Vui lòng thử lại sau.",
		Type:    "report_failed",
		Link:    "/interviews/" + interviewID + "/report",
	}); err != nil {
		logger.Error("failed to create report_failed notification", "interview_id", interviewID, "error", err)
	}
}

func (s *ReportService) RetryReport(ctx context.Context, companyID, interviewID string, requestedBy string) (*models.InterviewReport, error) {
	interview, err := s.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return nil, err
	}
	if interview == nil {
		return nil, apierrors.NewNotFound("interview")
	}

	if interview.ReportStatus == "generating" {
		return nil, apierrors.NewValidation("interview_id", []string{"report is currently generating"})
	}
	if interview.ReportStatus != "failed" {
		return nil, apierrors.NewValidation("interview_id", []string{"report is not in failed status"})
	}

	// Atomic: set generating
	ok, err := s.interviewRepo.UpdateReportStatusIf(ctx, interviewID, "generating", "failed")
	if err != nil {
		return nil, fmt.Errorf("failed to lock report generation: %w", err)
	}
	if !ok {
		return nil, apierrors.NewValidation("interview_id", []string{"report status changed, retry aborted"})
	}

	// Sync: run immediately, not goroutine
	err = s.generateReportSync(ctx, companyID, interviewID, interview.JobID.String, requestedBy)
	if err != nil {
		_ = s.interviewRepo.UpdateReportStatus(ctx, interviewID, "failed")
		s.notifyReportFailed(ctx, interviewID, interview.RecruiterID.String, err.Error())
		return nil, fmt.Errorf("retry failed: %w", err)
	}

	_ = s.interviewRepo.UpdateReportStatus(ctx, interviewID, "ready")

	report, err := s.reportRepo.GetByInterviewID(ctx, interviewID)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (s *ReportService) GetReport(ctx context.Context, companyID, interviewID string) (*models.InterviewReport, error) {
	// 0. Verify interview belongs to company (Prevent IDOR)
	interview, err := s.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return nil, err
	}
	if interview == nil {
		return nil, apierrors.NewNotFound("interview")
	}

	if interview.ReportStatus == "generating" {
		return &models.InterviewReport{
			InterviewID: interviewID,
			Summary:     "Báo cáo đang được xử lý bởi AI, vui lòng đợi trong giây lát...",
		}, nil
	}

	report, err := s.reportRepo.GetByInterviewID(ctx, interviewID)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, apierrors.NewNotFound("report")
	}
	return report, nil
}

func (s *ReportService) OverrideDecision(ctx context.Context, companyID, interviewID string, decision string, comment string) error {
	// 0. Verify interview belongs to company (Prevent IDOR)
	interview, err := s.interviewRepo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return err
	}
	if interview == nil {
		return apierrors.NewNotFound("interview")
	}
	return s.reportRepo.UpdateRecruiterDecision(ctx, interviewID, decision, comment)

}
