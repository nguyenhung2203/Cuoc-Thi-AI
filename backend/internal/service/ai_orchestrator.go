package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/sony/gobreaker/v2"

	"backend/internal/ai"
	"backend/internal/models"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/utils"
)

// StandardAIResponse is the expected structure returned from the Python ai-service.
type StandardAIResponse struct {
	Data             json.RawMessage `json:"data"`
	Evidence         string          `json:"evidence"`
	Confidence       float64         `json:"confidence"`
	InsufficientData bool            `json:"insufficient_data"`
}

// AIPayload is what we send to the Python ai-service.
type AIPayload struct {
	Model       string          `json:"model"`
	Prompt      string          `json:"prompt"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}

type AIOrchestratorService struct {
	promptSvc  *PromptService
	logSvc     *AILogService
	httpClient *http.Client
	breaker    *gobreaker.CircuitBreaker[[]byte]
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
			Timeout: 120 * time.Second,
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
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("GEMINI_API_KEY is not set")
		}

		model := strings.TrimSpace(tmpl.Model)
		if model == "" || strings.Contains(model, "1.5") || !strings.HasPrefix(model, "gemini") {
			model = "gemini-2.5-flash"
		}

		url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)
		
		reqBody := map[string]interface{}{
			"contents": []map[string]interface{}{
				{
					"parts": []map[string]interface{}{
						{"text": renderedPrompt},
					},
				},
			},
			"generationConfig": map[string]interface{}{
				"temperature": 0.7,
				"maxOutputTokens": 8192,
			},
		}
		reqBytes, _ := json.Marshal(reqBody)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBytes))
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
			return nil, fmt.Errorf("Gemini API returned 5xx status: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, backoff.Permanent(fmt.Errorf("Gemini API returned %d: %s", resp.StatusCode, string(body)))
		}

		// Parse Gemini response
		var geminiResp struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
		}
		if err := json.Unmarshal(body, &geminiResp); err != nil {
			return nil, err
		}
		
		if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
			return nil, fmt.Errorf("empty response from Gemini")
		}
		
		text := geminiResp.Candidates[0].Content.Parts[0].Text
		cleanText := utils.CleanJSON(text)
		
		// Ensure cleanText is a valid JSON object or wrap it
		if !json.Valid([]byte(cleanText)) {
			cleanText = `{"raw_text": ` + fmt.Sprintf("%q", cleanText) + `}`
		}

		// Construct StandardAIResponse JSON bytes
		stdResp := StandardAIResponse{
			Data: json.RawMessage(cleanText),
			Confidence: 0.9,
		}
		return json.Marshal(stdResp)
	}

	// Exponential backoff configuration
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 1 * time.Second
	bo.MaxElapsedTime = 120 * time.Second // Stop trying after 120s

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
		
		return nil, err
	}

	// 5. Parse Response
	var aiResp StandardAIResponse
	cleanJSON := utils.CleanJSON(string(responseBody))
	if parseErr := json.Unmarshal([]byte(cleanJSON), &aiResp); parseErr != nil {
		logEntry.Status = "failed"
		logEntry.Error = sql.NullString{String: parseErr.Error(), Valid: true}
		s.logSvc.LogAsync(logEntry)
		return nil, fmt.Errorf("parse error: %w", parseErr)
	}

	logEntry.Status = "success"
	logEntry.OutputJSON = models.JSONB(responseBody)
	s.logSvc.LogAsync(logEntry)

	if aiResp.InsufficientData {
		return nil, apierrors.NewValidation("insufficient data", []string{"AI needs more information to process this request"})
	}

	return &aiResp, nil
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
	Score     float64 `json:"score"`
	AIComment string  `json:"ai_comment"`
}

type ScoreResult struct {
	Score      float64
	Evidence   string
	AIComment  string
	Confidence float64
}

// ScoreAnswer calls AI to score a candidate's answer based on a rubric criterion
func (s *AIOrchestratorService) ScoreAnswer(ctx context.Context, companyID string, transcriptText string, criterionName, criterionDesc string, minScore, maxScore int, scoringGuide string) (*ScoreResult, error) {
	variables := map[string]string{
		"TRANSCRIPT":     transcriptText,
		"CRITERION_NAME": criterionName,
		"CRITERION_DESC": criterionDesc,
		"MIN_SCORE":      fmt.Sprintf("%d", minScore),
		"MAX_SCORE":      fmt.Sprintf("%d", maxScore),
		"SCORING_GUIDE":  scoringGuide,
	}

	fullResp, err := s.CallAIWithFullResponse(ctx, "SCORE_ANSWER", companyID, variables)
	if err != nil {
		return nil, err
	}

	var scoreData AIScoreData
	cleanJSON := utils.CleanJSON(string(fullResp.Data))
	if err := json.Unmarshal([]byte(cleanJSON), &scoreData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal score data: %w", err)
	}

	return &ScoreResult{
		Score:      scoreData.Score,
		AIComment:  scoreData.AIComment,
		Evidence:   fullResp.Evidence,
		Confidence: fullResp.Confidence,
	}, nil
}

// GenerateReport calls AI to generate a final interview report based on transcript, scores, and job requirements.
func (s *AIOrchestratorService) GenerateReport(ctx context.Context, companyID string, jobRequirements string, transcript string, scoresJSON string) (*ai.ReportGenerationResult, error) {
	variables := map[string]string{
		"JOB_REQUIREMENTS": jobRequirements,
		"TRANSCRIPT":       transcript,
		"SCORES":           scoresJSON,
	}

	fullResp, err := s.CallAIWithFullResponse(ctx, "GENERATE_REPORT", companyID, variables)
	if err != nil {
		return nil, err
	}

	var result ai.ReportGenerationResult
	cleanJSON := utils.CleanJSON(string(fullResp.Data))
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal report generation result: %w", err)
	}

	return &result, nil
}
