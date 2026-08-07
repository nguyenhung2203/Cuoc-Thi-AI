package queue

import "time"

// JobType identifies an async task kind (asynq task type name).
type JobType string

const (
	// TypeGenerateReport generates an interview report after it ends.
	TypeGenerateReport = "report:generate"
	// TypeAnalyzeCV parses/analyzes a candidate CV.
	TypeAnalyzeCV = "cv:analyze"
	// TypeBatchTranscript processes a batch of transcript segments.
	TypeBatchTranscript = "transcript:batch"
	// TypeRecomputeMatches recomputes all fit scores for a user after a CV change.
	TypeRecomputeMatches = "match:recompute"
	// TypeSendOTPEmail sends a registration or password-reset OTP email.
	TypeSendOTPEmail = "email:send-otp"
)

type SendOTPEmailPayload struct {
	To           string    `json:"to"`
	Purpose      string    `json:"purpose"`
	GenerationID string    `json:"generation_id"`
	TemplateType string    `json:"template_type"`
	CreatedAt    time.Time `json:"created_at"`
}

type GenerateReportPayload struct {
	CompanyID   string `json:"company_id"`
	InterviewID string `json:"interview_id"`
	JobID       string `json:"job_id"`
	GeneratedBy string `json:"generated_by"`
	RecruiterID string `json:"recruiter_id"`
}

// AnalyzeCVPayload is the JSON payload for a CV analysis job.
type AnalyzeCVPayload struct {
	CompanyID   string `json:"company_id"`
	CandidateID string `json:"candidate_id"`
}

// BatchTranscriptPayload is the JSON payload for a transcript batch job.
type BatchTranscriptPayload struct {
	CompanyID   string `json:"company_id"`
	InterviewID string `json:"interview_id"`
}

// RecomputeMatchesPayload is the JSON payload for recomputing a user's fit scores.
type RecomputeMatchesPayload struct {
	UserID string `json:"user_id"`
}
