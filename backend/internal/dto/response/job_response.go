package response

import "time"

type JobListItem struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Department     string    `json:"department,omitempty"`
	Level          string    `json:"level,omitempty"`
	Status         string    `json:"status"`
	CandidateCount int       `json:"candidate_count"`
	InterviewCount int       `json:"interview_count"`
	AvgFitScore    float64   `json:"avg_fit_score"`
	CreatedAt      time.Time `json:"created_at"`
}

type JobStats struct {
	CandidateCount   int `json:"candidate_count"`
	InterviewCount   int `json:"interview_count"`
	ReportReadyCount int `json:"report_ready_count"`
}

type JobDetail struct {
	ID             string      `json:"id"`
	CompanyID      string      `json:"company_id"`
	Title          string      `json:"title"`
	Department     string      `json:"department,omitempty"`
	Level          string      `json:"level,omitempty"`
	Location       string      `json:"location,omitempty"`
	EmploymentType string      `json:"employment_type,omitempty"`
	SalaryMin      float64     `json:"salary_min,omitempty"`
	SalaryMax      float64     `json:"salary_max,omitempty"`
	Currency       string      `json:"currency,omitempty"`
	Description    string      `json:"description"`
	Requirements   string      `json:"requirements,omitempty"`
	Benefits       string      `json:"benefits,omitempty"`
	Status         string      `json:"status"`
	AISummary      string      `json:"ai_summary,omitempty"`
	AIAnalysisJSON interface{} `json:"ai_analysis_json,omitempty"`
	Stats          JobStats    `json:"stats"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
