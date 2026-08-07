package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/sony/gobreaker/v2"

	"backend/internal/ai"
	"backend/internal/models"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/utils"
)

// aiServiceError normalises an AI failure for handlers: an existing *AppError
// (circuit-breaker 502, insufficient-data 422) passes through unchanged;
// anything else (timeouts, transport errors, unreadable responses) becomes a
// 502 AI_SERVICE_ERROR so callers never surface a raw 500 for an upstream
// AI outage.
func aiServiceError(err error, msg string) error {
	if err == nil {
		return nil
	}
	if appErr, ok := apierrors.IsAppError(err); ok {
		return appErr
	}
	return apierrors.NewAIServiceError(msg + ": " + err.Error())
}

// StandardAIResponse is the expected structure returned from the Python ai-service.
type StandardAIResponse struct {
	Data             json.RawMessage `json:"data"`
	Evidence         string          `json:"evidence"`
	Confidence       float64         `json:"confidence"`
	InsufficientData bool            `json:"insufficient_data"`
	TokensIn         int             `json:"tokens_in"`
	TokensOut        int             `json:"tokens_out"`
	Model            string          `json:"model"`
}

// AIPayload is what we send to the Python ai-service.
type AIPayload struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
}

type AIOrchestratorService struct {
	promptSvc    *PromptService
	logSvc       *AILogService
	httpClient   *http.Client
	breaker      *gobreaker.CircuitBreaker[[]byte]
	aiServiceURL string
}

func NewAIOrchestratorService(promptSvc *PromptService, logSvc *AILogService, aiServiceURL string) *AIOrchestratorService {
	cbSettings := gobreaker.Settings{
		Name:        "AI-Service",
		MaxRequests: 3,
		Interval:    30 * time.Second,
		Timeout:     60 * time.Second, // Tripped state lasts for 60 seconds
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip if there are more than 5 consecutive failures
			return counts.ConsecutiveFailures > 5
		},
	}

	return &AIOrchestratorService{
		promptSvc: promptSvc,
		logSvc:    logSvc,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		breaker:      gobreaker.NewCircuitBreaker[[]byte](cbSettings),
		aiServiceURL: aiServiceURL,
	}
}

// CallAIWithFullResponse orchestrates rendering the prompt, calling the Python service, and returns the full StandardAIResponse.
func (s *AIOrchestratorService) CallAIWithFullResponse(ctx context.Context, templateName, companyID string, variables map[string]string) (*StandardAIResponse, error) {
	startTime := time.Now()

	// 1. Load Prompt Template
	tmpl, err := s.promptSvc.LoadTemplate(ctx, templateName, companyID, 0)
	if err != nil {
		return nil, err
	}

	// 2. Render Prompt
	renderedPrompt := s.promptSvc.Render(tmpl.Content, variables)

	// Prepare request payload
	payload := AIPayload{
		Model:  tmpl.Model,
		Prompt: renderedPrompt,
	}

	// Extract temp/max_tokens from params if present (simplified logic)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 3. Setup Retry and Circuit Breaker
	var responseBody []byte

	operation := func() ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.aiServiceURL+"/api/v1/generate", bytes.NewReader(payloadBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 {
			// Server error, should trigger retry
			return nil, fmt.Errorf("AI service returned 5xx status: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			// e.g. 400 Bad Request, usually don't retry, but for simplicity returning error
			return nil, backoff.Permanent(fmt.Errorf("AI service returned %d: %s", resp.StatusCode, string(body)))
		}

		return body, nil
	}

	// Exponential backoff configuration
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 1 * time.Second
	bo.MaxElapsedTime = 15 * time.Second // Stop trying after 15s

	// Wrap operation in Circuit Breaker
	cbOperation := func() ([]byte, error) {
		return s.breaker.Execute(func() ([]byte, error) {
			return operation()
		})
	}

	// Execute with Retry
	responseBody, err = backoff.RetryWithData(cbOperation, backoff.WithContext(bo, ctx))

	latency := time.Since(startTime).Milliseconds()

	// 4. Log the Request Async
	logEntry := &models.AIRequestLog{
		CompanyID:       companyID,
		Provider:        "gemini",
		Model:           sql.NullString{String: tmpl.Model, Valid: tmpl.Model != ""},
		Operation:       sql.NullString{String: templateName, Valid: true},
		TemplateID:      sql.NullString{String: tmpl.ID, Valid: true},
		TemplateVersion: sql.NullInt32{Int32: int32(tmpl.Version), Valid: true},
		InputJSON:       models.JSONB(payloadBytes),
		LatencyMs:       sql.NullInt32{Int32: int32(latency), Valid: true},
		CreatedAt:       time.Now(),
	}

	if err != nil {
		logEntry.Status = "failed"
		logEntry.Error = sql.NullString{String: err.Error(), Valid: true}
		s.logSvc.LogAsync(logEntry)

		// Check if circuit breaker is open
		if err == gobreaker.ErrOpenState {
			return nil, &apierrors.AppError{
				Code:       apierrors.AI_SERVICE_ERROR,
				Message:    "AI service is currently unavailable",
				HTTPStatus: http.StatusBadGateway,
			}
		}

		return nil, aiServiceError(err, "AI service call failed")
	}

	// 5. Parse Response
	var aiResp StandardAIResponse
	cleanJSON := utils.CleanJSON(string(responseBody))
	if parseErr := json.Unmarshal([]byte(cleanJSON), &aiResp); parseErr != nil {
		logEntry.Status = "failed"
		logEntry.Error = sql.NullString{String: parseErr.Error(), Valid: true}
		s.logSvc.LogAsync(logEntry)
		return nil, aiServiceError(parseErr, "AI service returned an unreadable response")
	}

	logEntry.Status = "success"
	logEntry.OutputJSON = models.JSONB(responseBody)
	// Persist provider-reported usage, including legitimate zero values.
	logEntry.TokensIn = sql.NullInt32{Int32: int32(aiResp.TokensIn), Valid: true}
	logEntry.TokensOut = sql.NullInt32{Int32: int32(aiResp.TokensOut), Valid: true}
	logEntry.TotalTokens = sql.NullInt32{Int32: int32(aiResp.TokensIn + aiResp.TokensOut), Valid: true}
	if aiResp.Model != "" {
		logEntry.Model = sql.NullString{String: aiResp.Model, Valid: true}
	}
	// Cost estimation: Gemini 2.5 Flash ~$0.30/1M input, $1.25/1M output
	cost := float64(aiResp.TokensIn)*0.30/1_000_000 + float64(aiResp.TokensOut)*1.25/1_000_000
	if cost > 0 {
		logEntry.Cost = sql.NullFloat64{Float64: cost, Valid: true}
	}
	s.logSvc.LogAsync(logEntry)

	if aiResp.InsufficientData {
		return nil, apierrors.NewValidation("insufficient data", []string{"AI needs more information to process this request"})
	}

	return &aiResp, nil
}

// GenerateStructuredQuestions calls the typed question-generation endpoint.
func (s *AIOrchestratorService) GenerateStructuredQuestions(ctx context.Context, companyID string, request ai.StructuredQuestionRequest) (*ai.QuestionGenerationResult, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, aiServiceError(err, "failed to encode question request")
	}
	url := s.aiServiceURL + "/api/v1/questions/generate"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, aiServiceError(err, "failed to create question request")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, aiServiceError(err, "AI question request failed")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, aiServiceError(err, "failed to read question response")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apierrors.NewAIServiceError(fmt.Sprintf("AI question endpoint returned %d: %s", resp.StatusCode, string(body)))
	}
	var result ai.QuestionGenerationResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, aiServiceError(err, "AI question response was invalid")
	}
	if len(result.Questions) == 0 {
		return nil, apierrors.NewAIServiceError("AI question response contained no questions")
	}
	return &result, nil
}

// GenerateFollowUp calls the typed adaptive follow-up endpoint.
func (s *AIOrchestratorService) GenerateFollowUp(ctx context.Context, companyID string, request ai.StructuredQuestionRequest) (*ai.QuestionGenerationResult, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, aiServiceError(err, "failed to encode follow-up request")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.aiServiceURL+"/api/v1/questions/follow-up", bytes.NewReader(payload))
	if err != nil {
		return nil, aiServiceError(err, "failed to create follow-up request")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, aiServiceError(err, "AI follow-up request failed")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, aiServiceError(err, "failed to read follow-up response")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apierrors.NewAIServiceError(fmt.Sprintf("AI follow-up endpoint returned %d: %s", resp.StatusCode, string(body)))
	}
	var result ai.QuestionGenerationResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, aiServiceError(err, "AI follow-up response was invalid")
	}
	if len(result.Questions) != 1 {
		return nil, apierrors.NewAIServiceError("AI follow-up response must contain one question")
	}
	return &result, nil
}

// CallAI orchestrates rendering the prompt, calling the Python service with retry + circuit breaker, and returning only Data.
func (s *AIOrchestratorService) CallAI(ctx context.Context, templateName, companyID string, variables map[string]string) ([]byte, error) {
	resp, err := s.CallAIWithFullResponse(ctx, templateName, companyID, variables)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

type AIScoreData struct {
	Score             float64  `json:"score"`
	AIComment         string   `json:"ai_comment"`
	Communication     float64  `json:"communication"`
	Tone              float64  `json:"tone"`
	Personality       float64  `json:"personality"`
	Strengths         []string `json:"strengths"`
	Weaknesses        []string `json:"weaknesses"`
	ImprovementAdvice []string `json:"improvement_advice"`
}

type ScoreResult struct {
	Score             float64
	Evidence          string
	AIComment         string
	Confidence        float64
	Communication     float64
	Tone              float64
	Personality       float64
	Strengths         []string
	Weaknesses        []string
	ImprovementAdvice []string
}

func NormalizeCandidateLevel(level string) string {
	level = strings.ToLower(strings.TrimSpace(level))
	switch level {
	case "intern", "fresher", "entry", "entry-level":
		return "fresher"
	case "junior", "jr":
		return "junior"
	case "mid", "middle", "mid-level", "medior":
		return "mid"
	default:
		return "junior"
	}
}

func (s *AIOrchestratorService) ScoreAnswer(ctx context.Context, companyID string, transcriptText string, criterionName, criterionDesc string, minScore, maxScore int, scoringGuide string) (*ScoreResult, error) {
	return s.ScoreAnswerForLevel(ctx, companyID, transcriptText, criterionName, criterionDesc, minScore, maxScore, scoringGuide, "junior")
}

func (s *AIOrchestratorService) ScoreAnswerForLevel(ctx context.Context, companyID string, transcriptText string, criterionName, criterionDesc string, minScore, maxScore int, scoringGuide, candidateLevel string) (*ScoreResult, error) {
	variables := map[string]string{
		"transcript":      transcriptText,
		"criterion_name":  criterionName,
		"criterion_desc":  criterionDesc,
		"min_score":       fmt.Sprintf("%d", minScore),
		"max_score":       fmt.Sprintf("%d", maxScore),
		"candidate_level": NormalizeCandidateLevel(candidateLevel),
		"scoring_guide":   scoringGuide,
	}

	fullResp, err := s.CallAIWithFullResponse(ctx, "score_answer", companyID, variables)
	if err != nil {
		return nil, err
	}

	var scoreData AIScoreData
	cleanJSON := utils.CleanJSON(string(fullResp.Data))
	if err := json.Unmarshal([]byte(cleanJSON), &scoreData); err != nil {
		return nil, aiServiceError(err, "AI returned unreadable score data")
	}

	return &ScoreResult{
		Score:             scoreData.Score,
		AIComment:         scoreData.AIComment,
		Evidence:          fullResp.Evidence,
		Confidence:        fullResp.Confidence,
		Communication:     scoreData.Communication,
		Tone:              scoreData.Tone,
		Personality:       scoreData.Personality,
		Strengths:         scoreData.Strengths,
		Weaknesses:        scoreData.Weaknesses,
		ImprovementAdvice: scoreData.ImprovementAdvice,
	}, nil
}

// GenerateReport calls AI to generate a final interview report based on transcript, scores, and job requirements.
func (s *AIOrchestratorService) GenerateReport(ctx context.Context, companyID string, jobRequirements string, transcript string, scoresJSON string) (*ai.ReportGenerationResult, error) {
	variables := map[string]string{
		"job_requirements": jobRequirements,
		"transcript":       transcript,
		"scores":           scoresJSON,
	}

	fullResp, err := s.CallAIWithFullResponse(ctx, "generate_report", companyID, variables)
	if err != nil {
		return nil, err
	}

	var result ai.ReportGenerationResult
	cleanJSON := utils.CleanJSON(string(fullResp.Data))
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		return nil, aiServiceError(err, "AI returned an unreadable report")
	}

	return &result, nil
}
