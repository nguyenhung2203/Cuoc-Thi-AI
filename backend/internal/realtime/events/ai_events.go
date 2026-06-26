package events

// ── Client → Server ──────────────────────────────────────────────────────────

// AIRequestSuggestionPayload is sent by the recruiter to ask for a follow-up.
type AIRequestSuggestionPayload struct {
	Focus            string `json:"focus,omitempty"` // e.g. "technical_depth"
	LastTranscriptID string `json:"last_transcript_id,omitempty"`
}

// AIRequestScoreUpdatePayload asks the AI to refresh rubric scores.
type AIRequestScoreUpdatePayload struct {
	CriterionIDs []string `json:"criterion_ids,omitempty"`
	Scope        string   `json:"scope,omitempty"` // "latest_answer" | "full"
}

// ── Server → Client (Recruiter only, unless mock interview) ──────────────────

// AIThinkingPayload is sent immediately when the server starts an AI task,
// before the result is ready.
type AIThinkingPayload struct {
	Task               AITaskType `json:"task"`
	RequestID          string     `json:"request_id"`
	Message            string     `json:"message"`
	ExpectedDurationMs int        `json:"expected_duration_ms"`
}

// AISuggestionPayload carries a single AI-generated question/hint for the recruiter.
type AISuggestionPayload struct {
	SuggestionID   string           `json:"suggestion_id"`
	SuggestionType AISuggestionType `json:"suggestion_type"`
	Content        string           `json:"content"`
	Reason         string           `json:"reason"`
	TargetSkill    string           `json:"target_skill,omitempty"`
	Priority       string           `json:"priority"` // "low" | "medium" | "high"
	Confidence     float64          `json:"confidence"`
}

// RubricScore is a single criterion result inside an AI score update.
type RubricScore struct {
	CriterionName string      `json:"criterion_name"`
	Score         *float64    `json:"score"` // nil when insufficient evidence
	MaxScore      float64     `json:"max_score"`
	Evidence      string      `json:"evidence,omitempty"`
	AIComment     string      `json:"ai_comment,omitempty"`
	Confidence    float64     `json:"confidence"`
	Status        ScoreStatus `json:"status"`
}

// AIScoreUpdatePayload carries updated rubric scores for the recruiter.
type AIScoreUpdatePayload struct {
	Scores []RubricScore `json:"scores"`
}

// AIWarningPayload signals a data quality issue without stopping the interview.
type AIWarningPayload struct {
	WarningType     string `json:"warning_type"`
	Message         string `json:"message"`
	SuggestedAction string `json:"suggested_action,omitempty"`
}

// AIErrorPayload describes an AI failure and whether recovery is possible.
type AIErrorPayload struct {
	ErrorType         AIErrorType `json:"error_type"`
	RequestID         string      `json:"request_id,omitempty"`
	AffectedFeatures  []string    `json:"affected_features"`
	Severity          AISeverity  `json:"severity"`
	Message           string      `json:"message"`
	Recoverable       bool        `json:"recoverable"`
	RetryAfterSeconds int         `json:"retry_after_seconds,omitempty"`
}

// ReportReadyPayload is sent to the recruiter when the post-interview report is done.
type ReportReadyPayload struct {
	ReportID       string  `json:"report_id"`
	FinalScore     float64 `json:"final_score"`
	Recommendation string  `json:"recommendation"`
}

// ErrorPayload is the generic error envelope returned to the originating client.
type ErrorPayload struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Recoverable bool   `json:"recoverable"`
}
