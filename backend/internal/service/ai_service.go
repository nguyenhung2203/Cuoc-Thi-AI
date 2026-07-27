package service

import (
	"bytes"
	"context"
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

	"backend/internal/repository"
)

type AIService struct {
	fileRepo      *repository.FileRepository
	candidateRepo *repository.CandidateRepository
	geminiKeys    []string
	uploadDir     string
}

func NewAIService(
	fileRepo *repository.FileRepository,
	candidateRepo *repository.CandidateRepository,
	geminiKey string,
	uploadDir string,
) *AIService {
	rand.Seed(time.Now().UnixNano())
	
	// Read GEMINI_API_KEYS from env if present, else fallback to geminiKey
	keysStr := os.Getenv("GEMINI_API_KEYS")
	if keysStr == "" {
		keysStr = geminiKey
	}
	
	var keys []string
	for _, k := range strings.Split(keysStr, ",") {
		k = strings.TrimSpace(k)
		if k != "" {
			keys = append(keys, k)
		}
	}
	
	if len(keys) == 0 {
		keys = []string{""} // Fallback to empty string if missing
	}

	if uploadDir == "" {
		uploadDir = "uploads"
	}

	return &AIService{
		fileRepo:      fileRepo,
		candidateRepo: candidateRepo,
		geminiKeys:    keys,
		uploadDir:     uploadDir,
	}
}

// generateGeminiJSON sends a prompt to Gemini (new google.golang.org/genai SDK)
// requesting a JSON response, rotating through the configured API keys with
// retry on failure. Returns the cleaned JSON text (markdown fences stripped).
func (s *AIService) generateGeminiJSON(ctx context.Context, prompt string) (string, error) {
	if len(s.geminiKeys) == 0 || s.geminiKeys[0] == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not configured")
	}

	var lastErr error
	startIndex := rand.Intn(len(s.geminiKeys))
	for i := 0; i < len(s.geminiKeys); i++ {
		idx := (startIndex + i) % len(s.geminiKeys)
		key := s.geminiKeys[idx]

		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  key,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			lastErr = fmt.Errorf("failed to create gemini client: %w", err)
			continue
		}

		resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash",
			genai.Text(prompt),
			&genai.GenerateContentConfig{ResponseMIMEType: "application/json"},
		)
		if err != nil {
			lastErr = fmt.Errorf("gemini generation failed: %w", err)
			log.Printf("[WARN] AI Key %d failed, retrying with next key... error: %v", idx, err)
			continue
		}

		text := resp.Text()
		if text == "" {
			lastErr = fmt.Errorf("empty response from gemini")
			continue
		}

		clean := strings.TrimPrefix(text, "```json\n")
		clean = strings.TrimPrefix(clean, "```json")
		clean = strings.TrimSuffix(clean, "\n```")
		clean = strings.TrimSuffix(clean, "```")
		return strings.TrimSpace(clean), nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("all api keys failed")
	}
	return "", lastErr
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

	cleanJSON, err := s.generateGeminiJSON(ctx, prompt)
	if err != nil {
		return "", "", err
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

	cleanJSON, err := s.generateGeminiJSON(ctx, prompt)
	if err != nil {
		return err
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

