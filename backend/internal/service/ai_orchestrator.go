package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/sony/gobreaker/v2"

	"backend/internal/models"
	apierrors "backend/internal/pkg/errors"
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
			Timeout: 30 * time.Second,
		},
		breaker:      gobreaker.NewCircuitBreaker[[]byte](cbSettings),
		aiServiceURL: aiServiceURL,
	}
}

// CallAI orchestrates rendering the prompt, calling the Python service with retry + circuit breaker, and logging.
func (s *AIOrchestratorService) CallAI(ctx context.Context, templateName, companyID string, variables map[string]string) ([]byte, error) {
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
	// In reality you would parse tmpl.Params JSONB to get exact values.

	payloadBytes, _ := json.Marshal(payload)

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
	if parseErr := json.Unmarshal(responseBody, &aiResp); parseErr != nil {
		logEntry.Status = "failed"
		logEntry.Error = sql.NullString{String: parseErr.Error(), Valid: true}
		s.logSvc.LogAsync(logEntry)
		return nil, fmt.Errorf("parse error: %w", parseErr)
	}

	logEntry.Status = "success"
	logEntry.OutputJSON = models.JSONB(responseBody)
	s.logSvc.LogAsync(logEntry)

	if aiResp.InsufficientData {
		return nil, fmt.Errorf("insufficient data")
	}

	return aiResp.Data, nil
}

