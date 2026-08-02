package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
	"google.golang.org/genai"

	"backend/internal/models"
	"backend/internal/repository"
)

type AIService struct {
	fileRepo       *repository.FileRepository
	candidateRepo  *repository.CandidateRepository
	configProvider AIRuntimeConfigProvider
	uploadDir      string
	logSvc         *AILogService
}

func NewAIService(
	fileRepo *repository.FileRepository,
	candidateRepo *repository.CandidateRepository,
	configProvider AIRuntimeConfigProvider,
	uploadDir string,
) *AIService {
	rand.Seed(time.Now().UnixNano())
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	return &AIService{
		fileRepo: fileRepo, candidateRepo: candidateRepo,
		configProvider: configProvider, uploadDir: uploadDir,
	}
}

// SetLogService wires the AI log service so CV parse tokens are recorded.
func (s *AIService) SetLogService(logSvc *AILogService) {
	s.logSvc = logSvc
}

// generateGeminiJSON sends a prompt to Gemini (new google.golang.org/genai SDK)
// requesting a JSON response, rotating through the configured API keys with
// retry on failure. Returns the cleaned JSON text, tokens_in, tokens_out.
func (s *AIService) generateGeminiJSON(ctx context.Context, prompt string) (text string, tokensIn, tokensOut int32, err error) {
	if s.configProvider == nil {
		return "", 0, 0, fmt.Errorf("AI runtime configuration is not available")
	}
	cfg, err := s.configProvider.GetRuntimeConfig(ctx)
	if err != nil {
		return "", 0, 0, fmt.Errorf("load AI runtime configuration: %w", err)
	}
	if len(cfg.APIKeys) == 0 {
		return "", 0, 0, fmt.Errorf("GEMINI_API_KEY is not configured")
	}
	model := firstNonEmpty(cfg.TextModel, DefaultTextModel)

	var lastErr error
	startIndex := rand.Intn(len(cfg.APIKeys))
	for i := 0; i < len(cfg.APIKeys); i++ {
		idx := (startIndex + i) % len(cfg.APIKeys)
		key := cfg.APIKeys[idx]

		client, cerr := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  key,
			Backend: genai.BackendGeminiAPI,
		})
		if cerr != nil {
			lastErr = fmt.Errorf("failed to create gemini client: %w", cerr)
			continue
		}

		resp, gerr := client.Models.GenerateContent(ctx, model,
			genai.Text(prompt),
			&genai.GenerateContentConfig{ResponseMIMEType: "application/json"},
		)
		if gerr != nil {
			lastErr = fmt.Errorf("gemini generation failed: %w", gerr)
			log.Printf("[WARN] AI Key %d failed, retrying with next key... error: %v", idx, gerr)
			continue
		}

		rawText := resp.Text()
		if rawText == "" {
			lastErr = fmt.Errorf("empty response from gemini")
			continue
		}

		// Extract token counts from usage metadata
		var tIn, tOut int32
		if resp.UsageMetadata != nil {
			tIn = resp.UsageMetadata.PromptTokenCount
			tOut = resp.UsageMetadata.CandidatesTokenCount
		}

		clean := strings.TrimPrefix(rawText, "```json\n")
		clean = strings.TrimPrefix(clean, "```json")
		clean = strings.TrimSuffix(clean, "\n```")
		clean = strings.TrimSuffix(clean, "```")
		return strings.TrimSpace(clean), tIn, tOut, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("all api keys failed")
	}
	return "", 0, 0, lastErr
}

// ExtractTextFromPDF reads a local PDF file and extracts its text.
func ExtractTextFromPDF(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// Trim any leading whitespace or newlines which can break PDF parsers
	data = bytes.TrimLeft(data, " \t\r\n")

	reader := bytes.NewReader(data)
	r, err := pdf.NewReader(reader, int64(len(data)))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

type AIExtractionResult struct {
	Skills     []string `json:"skills"`
	Experience string   `json:"experience"`
	Education  string   `json:"education"`
}

// ExtractAndParseCV reads a CV PDF from disk (by storage key), asks the LLM to
// extract structured info, and returns the raw JSON plus a short summary.
// Reusable by both the recruiter candidate flow and the candidate portal.
func (s *AIService) ExtractAndParseCV(ctx context.Context, storageKey string) (string, string, error) {
	start := time.Now()
	filePath := filepath.Join(s.uploadDir, storageKey)
	text, err := ExtractTextFromPDF(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read PDF: %w", err)
	}

	prompt := fmt.Sprintf(`Analyze the following resume text and extract the key information into a JSON format with exactly three fields:
1. "skills": A list of strings representing the technical and soft skills.
2. "experience": A short summary string of their work experience (e.g., "3 years Backend Developer").
3. "education": A short summary string of their education (e.g., "BS Computer Science").

Resume Text:
%s`, text)

	cleanJSON, tokensIn, tokensOut, err := s.generateGeminiJSON(ctx, prompt)
	if err != nil {
		return "", "", err
	}

	// Log token usage to ai_request_logs (best-effort, non-blocking)
	if s.logSvc != nil {
		latency := time.Since(start).Milliseconds()
		cost := float64(tokensIn)*0.30/1_000_000 + float64(tokensOut)*1.25/1_000_000
		s.logSvc.LogAsync(&models.AIRequestLog{
			CompanyID: "", Provider: "gemini", Model: sql.NullString{String: DefaultTextModel, Valid: true},
			Operation: sql.NullString{String: "cv_extract", Valid: true}, InputJSON: models.JSONB(`{"template":"cv_extract"}`), OutputJSON: models.JSONB(cleanJSON),
			LatencyMs: sql.NullInt32{Int32: int32(latency), Valid: true}, TokensIn: sql.NullInt32{Int32: tokensIn, Valid: true},
			TokensOut: sql.NullInt32{Int32: tokensOut, Valid: true}, TotalTokens: sql.NullInt32{Int32: tokensIn + tokensOut, Valid: true},
			Cost: sql.NullFloat64{Float64: cost, Valid: cost > 0}, Status: "success", CreatedAt: time.Now(),
		})
	}

	var parsed AIExtractionResult
	if err := json.Unmarshal([]byte(cleanJSON), &parsed); err != nil {
		return "", "", fmt.Errorf("invalid json from gemini: %w", err)
	}

	summary := fmt.Sprintf("Candidate has skills in %s. Experience: %s. Education: %s.",
		strings.Join(parsed.Skills, ", "), parsed.Experience, parsed.Education)
	return cleanJSON, summary, nil
}

func (s *AIService) ParseCV(ctx context.Context, companyID, candidateID string) error {
	start := time.Now()
	candidate, err := s.candidateRepo.GetByID(ctx, companyID, candidateID)
	if err != nil || candidate == nil {
		return fmt.Errorf("candidate not found")
	}

	if !candidate.CVFileID.Valid {
		return fmt.Errorf("candidate has no CV file")
	}

	fileRecord, err := s.fileRepo.GetByIDAndCompany(ctx, candidate.CVFileID.String, companyID)
	if err != nil || fileRecord == nil {
		return fmt.Errorf("CV file not found")
	}

	// 1. Read PDF text
	filePath := filepath.Join(s.uploadDir, fileRecord.StorageKey)
	text, err := ExtractTextFromPDF(filePath)
	if err != nil {
		return fmt.Errorf("failed to read PDF: %w", err)
	}

	// 2. Setup Gemini
	prompt := fmt.Sprintf(`Analyze the following resume text and extract the key information into a JSON format with exactly three fields:
1. "skills": A list of strings representing the technical and soft skills.
2. "experience": A short summary string of their work experience (e.g., "3 years Backend Developer").
3. "education": A short summary string of their education (e.g., "BS Computer Science").

Resume Text:
%s`, text)

	cleanJSON, tokensIn, tokensOut, err := s.generateGeminiJSON(ctx, prompt)
	if err != nil {
		return err
	}
	if s.logSvc != nil {
		latency := time.Since(start).Milliseconds()
		cost := float64(tokensIn)*0.30/1_000_000 + float64(tokensOut)*1.25/1_000_000
		s.logSvc.LogAsync(&models.AIRequestLog{
			CompanyID: companyID, CandidateID: sql.NullString{String: candidateID, Valid: true}, Provider: "gemini", Model: sql.NullString{String: DefaultTextModel, Valid: true},
			Operation: sql.NullString{String: "cv_parse", Valid: true}, InputJSON: models.JSONB(`{"template":"cv_parse"}`), OutputJSON: models.JSONB(cleanJSON),
			LatencyMs: sql.NullInt32{Int32: int32(latency), Valid: true}, TokensIn: sql.NullInt32{Int32: tokensIn, Valid: true}, TokensOut: sql.NullInt32{Int32: tokensOut, Valid: true},
			TotalTokens: sql.NullInt32{Int32: tokensIn + tokensOut, Valid: true}, Cost: sql.NullFloat64{Float64: cost, Valid: cost > 0}, Status: "success", CreatedAt: time.Now(),
		})
	}

	// Validate JSON
	var parsed AIExtractionResult
	if err := json.Unmarshal([]byte(cleanJSON), &parsed); err != nil {
		return fmt.Errorf("invalid json from gemini: %w", err)
	}

	// AI Summary (Short descriptive text)
	summary := fmt.Sprintf("Candidate has skills in %s. Experience: %s. Education: %s.",
		strings.Join(parsed.Skills, ", "), parsed.Experience, parsed.Education)

	// 3. Save back to candidate
	patch := map[string]any{
		"parsed_cv_json": json.RawMessage(cleanJSON),
		"ai_cv_summary":  summary,
	}
	_, err = s.candidateRepo.Update(ctx, companyID, candidateID, patch)
	if err != nil {
		log.Printf("failed to update candidate cv parsing result: %v", err)
		return err
	}

	return nil
}
