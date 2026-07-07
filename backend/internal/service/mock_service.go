package service

import (
	"context"
	"database/sql"
	"encoding/json"


	"backend/internal/models"
	"backend/internal/repository"
)

type MockService struct {
	repo  *repository.MockInterviewRepository
	aiSvc *AIService
}

func NewMockService(repo *repository.MockInterviewRepository, aiSvc *AIService) *MockService {
	return &MockService{
		repo:  repo,
		aiSvc: aiSvc,
	}
}

func (s *MockService) CreateMockInterview(ctx context.Context, userID, role, level, cvFileID string) (*models.MockInterview, error) {
	mi := &models.MockInterview{
		UserID:      userID,
		TargetRole:  role,
		TargetLevel: sql.NullString{String: level, Valid: level != ""},
		CVFileID:    sql.NullString{String: cvFileID, Valid: cvFileID != ""},
		Status:      "in_progress",
	}

	err := s.repo.Create(ctx, mi)
	if err != nil {
		return nil, err
	}

	// Trigger first AI message
	aiResult, err := s.aiSvc.GenerateMockResponse(ctx, role, level, "Hello, I am ready to start the interview.", nil)
	if err == nil && aiResult != nil {
		s.repo.AddMessage(ctx, &models.MockInterviewMessage{
			MockInterviewID: mi.ID,
			SenderType:      "ai",
			Content:         aiResult.AIResponse,
			QuestionType:    sql.NullString{String: aiResult.QuestionType, Valid: true},
		})
	}

	return mi, nil
}

func (s *MockService) GetMessages(ctx context.Context, mockID string) ([]models.MockInterviewMessage, error) {
	return s.repo.GetMessages(ctx, mockID)
}

func (s *MockService) GetMockInterview(ctx context.Context, mockID string) (*models.MockInterview, error) {
	return s.repo.GetByID(ctx, mockID)
}

func (s *MockService) ProcessMessage(ctx context.Context, mockID, content string) error {
	mi, err := s.repo.GetByID(ctx, mockID)
	if err != nil {
		return err
	}

	// 1. Save candidate message
	err = s.repo.AddMessage(ctx, &models.MockInterviewMessage{
		MockInterviewID: mockID,
		SenderType:      "candidate",
		Content:         content,
	})
	if err != nil {
		return err
	}

	// 2. Fetch history
	history, err := s.repo.GetMessages(ctx, mockID)
	if err != nil {
		return err
	}

	// 3. Generate AI Response
	aiResult, err := s.aiSvc.GenerateMockResponse(ctx, mi.TargetRole, mi.TargetLevel.String, content, history)
	if err != nil {
		return err
	}

	// 4. Save AI Response
	scoreJSON, _ := json.Marshal(map[string]int{"score": aiResult.Score})
	return s.repo.AddMessage(ctx, &models.MockInterviewMessage{
		MockInterviewID: mockID,
		SenderType:      "ai",
		Content:         aiResult.AIResponse,
		QuestionType:    sql.NullString{String: aiResult.QuestionType, Valid: true},
		ScoreJSON:       sql.NullString{String: string(scoreJSON), Valid: true},
	})
}

func (s *MockService) FinishMockInterview(ctx context.Context, mockID string) error {
	history, err := s.repo.GetMessages(ctx, mockID)
	if err != nil {
		return err
	}

	totalScore := 0
	count := 0
	for _, m := range history {
		if m.SenderType == "ai" && m.ScoreJSON.Valid {
			var sj map[string]int
			if err := json.Unmarshal([]byte(m.ScoreJSON.String), &sj); err == nil {
				if score, ok := sj["score"]; ok {
					totalScore += score
					count++
				}
			}
		}
	}

	finalScore := 0.0
	if count > 0 {
		finalScore = float64(totalScore) / float64(count)
	}

	return s.repo.Finish(ctx, mockID, finalScore, "Interview completed successfully.")

	"fmt"
	"strings"
	"time"

	"backend/internal/models"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/logger"
	"backend/internal/pkg/utils"
	"backend/internal/repository"

	"github.com/google/uuid"
)

type MockService struct {
	mockRepo     *repository.MockRepository
	orchestrator *AIOrchestratorService
	promptSvc    *PromptService
}

func NewMockService(mockRepo *repository.MockRepository, orchestrator *AIOrchestratorService, promptSvc *PromptService) *MockService {
	return &MockService{
		mockRepo:     mockRepo,
		orchestrator: orchestrator,
		promptSvc:    promptSvc,
	}
}

type CreateMockRequest struct {
	TargetRole  string `json:"target_role"`
	TargetLevel string `json:"target_level"`
	CVFileID    string `json:"cv_file_id,omitempty"`
}

// getByIDAndCheckOwnership is a shared helper for ownership guard.
func (s *MockService) getByIDAndCheckOwnership(ctx context.Context, id, userID string) (*models.MockInterview, error) {
	m, err := s.mockRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, apierrors.NewNotFound("mock_interview")
	}
	if m.UserID != userID {
		return nil, apierrors.NewForbidden("access denied")
	}
	return m, nil
}

// K-S6-01: Create
func (s *MockService) Create(ctx context.Context, userID string, req CreateMockRequest) (*models.MockInterview, error) {
	m := &models.MockInterview{
		ID:         uuid.New().String(),
		UserID:     userID,
		TargetRole: req.TargetRole,
		Status:     "draft",
	}
	if req.TargetLevel != "" {
		m.TargetLevel = sql.NullString{String: req.TargetLevel, Valid: true}
	}
	if req.CVFileID != "" {
		m.CVFileID = sql.NullString{String: req.CVFileID, Valid: true}
	}
	if err := s.mockRepo.Create(ctx, m); err != nil {
		return nil, fmt.Errorf("failed to create mock interview: %w", err)
	}

	// Auto-generate first AI question in background
	go func(interviewID, targetRole, targetLevel, cvFileID string) {
		bgCtx := context.Background()
		firstQ, err := s.generateNextQuestion(bgCtx, interviewID, targetRole, targetLevel, cvFileID, nil)
		if err != nil {
			logger.Error("mock: failed to generate first question", "id", interviewID, "error", err)
			// Fallback question so interview is not empty
			firstQ = "Xin chào! Hãy giới thiệu ngắn gọn về bản thân và kinh nghiệm của bạn."
		}
		if err := s.mockRepo.CreateMessage(bgCtx, &models.MockInterviewMessage{
			ID:              uuid.New().String(),
			MockInterviewID: interviewID,
			SenderType:      "ai",
			Content:         firstQ,
		}); err != nil {
			logger.Error("mock: failed to save first question", "id", interviewID, "error", err)
		}
	}(m.ID, req.TargetRole, req.TargetLevel, req.CVFileID)

	return m, nil
}

// K-S6-01: List by user
func (s *MockService) ListByUser(ctx context.Context, userID string, page, pageSize int) ([]models.MockInterview, error) {
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.mockRepo.ListByUserID(ctx, userID, pageSize, offset)
}

// K-S6-01: Get detail
func (s *MockService) GetByID(ctx context.Context, id, userID string) (*models.MockInterview, error) {
	return s.getByIDAndCheckOwnership(ctx, id, userID)
}

// K-S6-01: Start — use atomic WHERE guard
func (s *MockService) Start(ctx context.Context, id, userID string) error {
	m, err := s.getByIDAndCheckOwnership(ctx, id, userID)
	if err != nil {
		return err
	}
	if m.Status != "draft" {
		return apierrors.NewValidation("status", []string{"can only start a draft session"})
	}
	// Atomic: only update if still draft
	ok, err := s.mockRepo.UpdateStatusIf(ctx, id, "active", "draft")
	if err != nil {
		return err
	}
	if !ok {
		return apierrors.NewValidation("status", []string{"session was already started"})
	}
	return s.mockRepo.UpdateStartedAt(ctx, id, time.Now())
}

// K-S6-01: End — use atomic WHERE guard
func (s *MockService) End(ctx context.Context, id, userID string) error {
	m, err := s.getByIDAndCheckOwnership(ctx, id, userID)
	if err != nil {
		return err
	}
	if m.Status != "active" {
		return apierrors.NewValidation("status", []string{"can only end an active session"})
	}
	// Atomic: only update if still active
	ok, err := s.mockRepo.UpdateStatusIf(ctx, id, "completed", "active")
	if err != nil {
		return err
	}
	if !ok {
		return apierrors.NewValidation("status", []string{"session was already ended"})
	}
	return s.mockRepo.UpdateEndedAt(ctx, id, time.Now())
}

// K-S6-02: Send message + AI response
func (s *MockService) SendMessage(ctx context.Context, id, userID, content string) (*models.MockInterviewMessage, *models.MockInterviewMessage, error) {
	m, err := s.getByIDAndCheckOwnership(ctx, id, userID)
	if err != nil {
		return nil, nil, err
	}
	if m.Status != "active" {
		return nil, nil, apierrors.NewValidation("status", []string{"session is not active"})
	}

	// Limit: max 20 candidate messages
	count, err := s.mockRepo.CountMessagesBySender(ctx, id, "candidate")
	if err != nil {
		return nil, nil, err
	}
	if count >= 20 {
		return nil, nil, apierrors.NewValidation("messages", []string{"maximum question limit reached (20)"})
	}

	// Use transaction-like DB batch: save user msg, load history, generate AI, save AI msg
	userMsg := &models.MockInterviewMessage{
		ID:              uuid.New().String(),
		MockInterviewID: id,
		SenderType:      "candidate",
		Content:         content,
	}
	if err := s.mockRepo.CreateMessage(ctx, userMsg); err != nil {
		return nil, nil, err
	}

	// Load full history for AI context
	messages, err := s.mockRepo.ListMessages(ctx, id)
	if err != nil {
		// If load fails, fallback: save generic AI question to keep conversation flowing
		logger.Error("mock: failed to list messages", "id", id, "error", err)
		aiMsg := &models.MockInterviewMessage{
			ID:              uuid.New().String(),
			MockInterviewID: id,
			SenderType:      "ai",
			Content:         "Cảm ơn câu trả lời của bạn. Hãy tiếp tục với câu hỏi tiếp theo.",
		}
		if err2 := s.mockRepo.CreateMessage(ctx, aiMsg); err2 != nil {
			return nil, nil, err2
		}
		return userMsg, aiMsg, nil
	}

	// AI call with timeout context
	aiCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	aiQuestion, err := s.generateNextQuestion(aiCtx, id, m.TargetRole, m.TargetLevel.String, m.CVFileID.String, messages)
	if err != nil {
		logger.Error("mock: AI question generation failed", "id", id, "error", err)
		aiQuestion = "Cảm ơn câu trả lời của bạn. Hãy tiếp tục với câu hỏi tiếp theo."
	}

	aiMsg := &models.MockInterviewMessage{
		ID:              uuid.New().String(),
		MockInterviewID: id,
		SenderType:      "ai",
		Content:         aiQuestion,
	}
	if err := s.mockRepo.CreateMessage(ctx, aiMsg); err != nil {
		return nil, nil, err
	}

	return userMsg, aiMsg, nil
}

// K-S6-03: Score answer
func (s *MockService) ScoreAnswer(ctx context.Context, mockID, messageID, userID string) error {
	return apierrors.NewValidation("not_implemented", []string{"scoring not yet implemented for mock interviews"})
}

// K-S6-04: Get report
func (s *MockService) GetReport(ctx context.Context, id, userID string) (*models.MockInterview, error) {
	return s.getByIDAndCheckOwnership(ctx, id, userID)
}

func (s *MockService) generateNextQuestion(ctx context.Context, mockID, targetRole, targetLevel, cvFileID string, prevMessages []models.MockInterviewMessage) (string, error) {
	variables := map[string]string{
		"TARGET_ROLE":  targetRole,
		"TARGET_LEVEL": targetLevel,
	}

	// Build conversation history
	var history []string
	for _, msg := range prevMessages {
		history = append(history, fmt.Sprintf("%s: %s", msg.SenderType, msg.Content))
	}
	if len(history) > 0 {
		variables["CONVERSATION"] = strings.Join(history, "\n")
	} else {
		variables["CONVERSATION"] = "Chưa có hội thoại. Hãy bắt đầu phỏng vấn."
	}

	// Use empty string as companyID — mock is user-scoped, not company-scoped.
	// CallAI falls back to system-wide prompt templates when companyID is empty.
	data, err := s.orchestrator.CallAI(ctx, "MOCK_QUESTION", "", variables)
	if err != nil {
		return "", err
	}
	result := string(data)
	cleanJSON := utils.CleanJSON(result)

	// Unmarshal to extract question_text
	var parsed struct {
		QuestionText string `json:"question_text"`
	}
	if err := json.Unmarshal([]byte(cleanJSON), &parsed); err != nil || parsed.QuestionText == "" {
		return result, nil
	}
	return parsed.QuestionText, nil
}
