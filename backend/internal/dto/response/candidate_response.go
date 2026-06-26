package response

import "time"

type CVFile struct {
	ID           string `json:"id"`
	OriginalName string `json:"original_name"`
	DownloadURL  string `json:"download_url"` // signed URL
}

type CandidateListItem struct {
	ID                string     `json:"id"`
	FullName          string     `json:"full_name"`
	Email             string     `json:"email"`
	Phone             string     `json:"phone,omitempty"`
	Status            string     `json:"status"`
	Source            string     `json:"source,omitempty"`
	FitScore          float64    `json:"fit_score,omitempty"`
	LatestInterviewAt *time.Time `json:"latest_interview_at,omitempty"`
}

type CandidateDetail struct {
	ID          string     `json:"id"`
	CompanyID   string     `json:"company_id"`
	FullName    string     `json:"full_name"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone,omitempty"`
	AvatarURL   string     `json:"avatar_url,omitempty"`
	Status      string     `json:"status"`
	Source      string     `json:"source,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	CVFile      *CVFile    `json:"cv_file,omitempty"`
	AICVSummary string     `json:"ai_cv_summary,omitempty"`
	Skills      []string   `json:"skills,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type JobCandidateItem struct {
	ID             string     `json:"id"`
	CandidateID    string     `json:"candidate_id"`
	FullName       string     `json:"full_name"`
	Email          string     `json:"email"`
	PipelineStatus string     `json:"pipeline_status"`
	FitScore       float64    `json:"fit_score,omitempty"`
	AppliedAt      *time.Time `json:"applied_at,omitempty"`
}
