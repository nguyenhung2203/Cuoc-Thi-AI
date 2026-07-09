package request

type ScoreAnswerRequest struct {
	TranscriptIDs []string `json:"transcript_ids" validate:"required,min=1"`
	CriterionIDs  []string `json:"criterion_ids" validate:"required,min=1"`
}

type GenerateReportRequest struct {
	JobID string `json:"job_id" validate:"required"`
}

type SuggestFollowUpRequest struct {
	LastTranscriptID string `json:"last_transcript_id"`
	Focus            string `json:"focus"`
}
