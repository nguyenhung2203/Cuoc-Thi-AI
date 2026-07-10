package response

import "time"

// CriterionScoreResponse is one rubric criterion's AI score with evidence.
type CriterionScoreResponse struct {
	CriterionName string  `json:"criterion_name"`
	Score         float64 `json:"score"`
	MaxScore      float64 `json:"max_score"`
	Weight        float64 `json:"weight"`
	WeightedScore float64 `json:"weighted_score"`
	Evidence      string  `json:"evidence"`
	AIComment     string  `json:"ai_comment"`
	Confidence    float64 `json:"confidence"`
	Status        string  `json:"status"`
}

// ReportResponse is the API-facing shape of an interview report.
type ReportResponse struct {
	ID               string                    `json:"id"`
	InterviewID      string                    `json:"interview_id"`
	Summary          string                    `json:"summary"`
	FinalScore       float64                   `json:"final_score"`
	Recommendation   string                    `json:"recommendation"`
	Strengths        []string                  `json:"strengths"`
	Weaknesses       []string                  `json:"weaknesses"`
	Risks            []string                  `json:"risks"`
	Scores           []CriterionScoreResponse  `json:"scores"`
	RecruiterDecision string                   `json:"recruiter_decision,omitempty"`
	RecruiterComment  string                   `json:"recruiter_comment,omitempty"`
	ReportStatus     string                    `json:"report_status"`
	CreatedAt        time.Time                 `json:"created_at"`
}
