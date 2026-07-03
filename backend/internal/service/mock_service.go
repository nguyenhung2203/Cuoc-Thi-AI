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
}
