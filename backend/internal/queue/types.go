package queue

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
)

// GenerateReportPayload is the JSON payload for a report generation job.
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
