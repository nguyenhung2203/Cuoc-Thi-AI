package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"backend/internal/ai"
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
		if err := s.mockRepo.UpdateAIState(bgCtx, interviewID, "question", "processing", nil); err != nil {
			logger.Error("mock: failed to mark question processing", "id", interviewID, "error", err)
			return
		}
		firstQ, err := s.generateNextQuestion(bgCtx, interviewID, targetRole, targetLevel, cvFileID, nil)
		if err != nil {
			_ = s.mockRepo.UpdateAIState(bgCtx, interviewID, "question", "failed", err)
			logger.Error("mock: failed to generate first question", "id", interviewID, "error", err)
			return
		}
		if err := s.mockRepo.CreateMessage(bgCtx, &models.MockInterviewMessage{
			ID:              uuid.New().String(),
			MockInterviewID: interviewID,
			SenderType:      "ai",
			Content:         firstQ,
		}); err != nil {
			_ = s.mockRepo.UpdateAIState(bgCtx, interviewID, "question", "failed", err)
			logger.Error("mock: failed to save first question", "id", interviewID, "error", err)
			return
		}
		_ = s.mockRepo.UpdateAIState(bgCtx, interviewID, "question", "ready", nil)
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
	if err := s.mockRepo.UpdateEndedAt(ctx, id, time.Now()); err != nil {
		return err
	}

	// Scoring is part of completing an AI mock interview. Persist its lifecycle
	// separately so a provider failure remains visible and retryable.
	scoreCtx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	if err := s.mockRepo.UpdateAIState(scoreCtx, id, "scoring", "processing", nil); err != nil {
		return err
	}
	if err := s.scoreSession(scoreCtx, m); err != nil {
		_ = s.mockRepo.UpdateAIState(context.Background(), id, "scoring", "failed", err)
		logger.Error("mock: session scoring failed", "id", id, "error", err)
		return fmt.Errorf("failed to score mock interview: %w", err)
	}
	if err := s.mockRepo.UpdateAIState(scoreCtx, id, "scoring", "ready", nil); err != nil {
		return err
	}
	return nil
}

// mockFeedbackResult mirrors the mock_feedback prompt template output.
type mockFeedbackResult struct {
	Score              float64  `json:"score"`
	MaxScore           float64  `json:"max_score"`
	Strengths          []string `json:"strengths"`
	Improvements       []string `json:"improvements"`
	CommunicationScore float64  `json:"communication_score"`
	ToneScore          float64  `json:"tone_score"`
	PersonalityScore   float64  `json:"personality_score"`
	CoachComment       string   `json:"coach_comment"`
}

// scoreJSONShape is what MockResultPage.vue reads off each AI message.
type scoreJSONShape struct {
	Score         float64  `json:"score"`
	GoodPoints    []string `json:"good_points"`
	ImprovePoints []string `json:"improve_points"`
}

// scoreSession scores each (AI question -> candidate answer) pair, writes a
// normalized /10 score_json onto the AI message that follows each answer, and
// persists final_score + feedback_json on the session.
func (s *MockService) scoreSession(ctx context.Context, m *models.MockInterview) error {
	messages, err := s.mockRepo.ListMessages(ctx, m.ID)
	if err != nil {
		return err
	}

	var total float64
	var scored int
	var comments []string

	for i, msg := range messages {
		if msg.SenderType != "candidate" {
			continue
		}
		// The question is the most recent AI message before this answer.
		question := ""
		for j := i - 1; j >= 0; j-- {
			if messages[j].SenderType == "ai" {
				question = messages[j].Content
				break
			}
		}
		if question == "" {
			continue
		}

		fb, ferr := s.scoreAnswer(ctx, question, msg.Content, m.TargetRole)
		if ferr != nil {
			logger.Error("mock: failed to score answer", "id", m.ID, "error", ferr)
			continue
		}

		// Normalize to /10.
		max := fb.MaxScore
		if max <= 0 {
			max = 5
		}
		norm := fb.Score / max * 10
		if norm > 10 {
			norm = 10
		}
		if norm < 0 {
			norm = 0
		}

		sj := scoreJSONShape{
			Score:         norm,
			GoodPoints:    fb.Strengths,
			ImprovePoints: fb.Improvements,
		}
		sjBytes, merr := json.Marshal(sj)
		if merr != nil {
			continue
		}

		// Target AI message = the AI message that follows this answer (matches
		// the frontend's nextAi lookup). If none (e.g. last voice turn), create
		// one carrying the coach comment.
		targetAIID := ""
		for j := i + 1; j < len(messages); j++ {
			if messages[j].SenderType == "ai" {
				targetAIID = messages[j].ID
				break
			}
		}
		if targetAIID == "" {
			newMsg := &models.MockInterviewMessage{
				ID:              uuid.New().String(),
				MockInterviewID: m.ID,
				SenderType:      "ai",
				Content:         fb.CoachComment,
				ScoreJSON:       models.JSONB(sjBytes),
			}
			if cerr := s.mockRepo.CreateMessage(ctx, newMsg); cerr != nil {
				logger.Error("mock: failed to create feedback message", "id", m.ID, "error", cerr)
			}
		} else if uerr := s.mockRepo.UpdateMessageScore(ctx, targetAIID, models.JSONB(sjBytes)); uerr != nil {
			logger.Error("mock: failed to save message score", "id", m.ID, "error", uerr)
		}

		total += norm
		scored++
		if fb.CoachComment != "" {
			comments = append(comments, fb.CoachComment)
		}
	}

	if scored == 0 {
		return fmt.Errorf("no interview answers received a valid AI score")
	}

	finalScore := sql.NullFloat64{Float64: total / float64(scored), Valid: true}
	feedback := strings.Join(comments, "\n\n")
	feedbackBytes, _ := json.Marshal(feedback)
	return s.mockRepo.SaveFeedback(ctx, m.ID, models.JSONB(feedbackBytes), finalScore)
}

// scoreAnswer scores a single Q&A pair via the mock_feedback prompt template.
func (s *MockService) scoreAnswer(ctx context.Context, question, answer, targetRole string) (*mockFeedbackResult, error) {
	variables := map[string]string{
		"question":    question,
		"answer":      answer,
		"target_role": targetRole,
	}
	data, err := s.orchestrator.CallAI(ctx, "mock_feedback", "", variables)
	if err != nil {
		return nil, err
	}
	var fb mockFeedbackResult
	cleanJSON := utils.CleanJSON(string(data))
	if err := json.Unmarshal([]byte(cleanJSON), &fb); err != nil {
		return nil, fmt.Errorf("failed to parse mock feedback: %w", err)
	}
	return &fb, nil
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
		logger.Error("mock: failed to list messages", "id", id, "error", err)
		return userMsg, nil, fmt.Errorf("failed to load interview history: %w", err)
	}

	// AI call with timeout context
	aiCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := s.mockRepo.UpdateAIState(aiCtx, id, "question", "processing", nil); err != nil {
		return userMsg, nil, err
	}
	aiQuestion, err := s.generateNextQuestion(aiCtx, id, m.TargetRole, m.TargetLevel.String, m.CVFileID.String, messages)
	if err != nil {
		_ = s.mockRepo.UpdateAIState(context.Background(), id, "question", "failed", err)
		logger.Error("mock: AI question generation failed", "id", id, "error", err)
		return userMsg, nil, fmt.Errorf("failed to generate AI question: %w", err)
	}

	aiMsg := &models.MockInterviewMessage{
		ID:              uuid.New().String(),
		MockInterviewID: id,
		SenderType:      "ai",
		Content:         aiQuestion,
	}
	if err := s.mockRepo.CreateMessage(ctx, aiMsg); err != nil {
		_ = s.mockRepo.UpdateAIState(context.Background(), id, "question", "failed", err)
		return nil, nil, err
	}
	if err := s.mockRepo.UpdateAIState(ctx, id, "question", "ready", nil); err != nil {
		return userMsg, aiMsg, err
	}

	return userMsg, aiMsg, nil
}

// LiveTurn is one spoken exchange captured from a Gemini Live mock session.
type LiveTurn struct {
	Role string // "ai" or "candidate"
	Text string
}

// SaveLiveTranscript persists a spoken mock conversation as messages. Called
// once when a voice mock session ends so the report has conversation material.
func (s *MockService) SaveLiveTranscript(ctx context.Context, id, userID string, turns []LiveTurn) error {
	if _, err := s.getByIDAndCheckOwnership(ctx, id, userID); err != nil {
		return err
	}
	for _, t := range turns {
		sender := "candidate"
		if t.Role == "ai" {
			sender = "ai"
		}
		if err := s.mockRepo.CreateMessage(ctx, &models.MockInterviewMessage{
			ID:              uuid.New().String(),
			MockInterviewID: id,
			SenderType:      sender,
			Content:         t.Text,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *MockService) GetMessages(ctx context.Context, id, userID string) ([]models.MockInterviewMessage, error) {
	_, err := s.getByIDAndCheckOwnership(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return s.mockRepo.ListMessages(ctx, id)
}

// K-S6-04: Get report
func (s *MockService) GetReport(ctx context.Context, id, userID string) (*models.MockInterview, error) {
	return s.getByIDAndCheckOwnership(ctx, id, userID)
}

func (s *MockService) generateNextQuestion(ctx context.Context, mockID, targetRole, targetLevel, cvFileID string, prevMessages []models.MockInterviewMessage) (string, error) {
	previousQuestions := make([]string, 0, len(prevMessages))
	recentAnswer := ""
	for _, msg := range prevMessages {
		if msg.SenderType == "ai" {
			previousQuestions = append(previousQuestions, utils.TruncateText(msg.Content, 1000))
		} else if msg.SenderType == "candidate" {
			recentAnswer = utils.TruncateText(msg.Content, 5000)
		}
	}
	level := strings.ToLower(strings.TrimSpace(targetLevel))
	if level == "mid" {
		level = "middle"
	}
	if level == "" {
		level = "unknown"
	}
	result, err := s.orchestrator.GenerateStructuredQuestions(ctx, "", ai.StructuredQuestionRequest{
		JobTitle: targetRole, Level: level, Mode: "mock", Language: "vi",
		QuestionCount: 1, PreviousQuestions: previousQuestions, RecentAnswer: recentAnswer,
	})
	if err != nil {
		return "", err
	}
	if result == nil || len(result.Questions) == 0 || strings.TrimSpace(result.Questions[0].QuestionText) == "" {
		return "", fmt.Errorf("AI question response is missing question_text")
	}
	return result.Questions[0].QuestionText, nil
}
