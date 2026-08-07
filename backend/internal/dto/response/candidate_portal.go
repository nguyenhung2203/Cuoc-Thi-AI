package response

import "time"

type CandidatePortalDashboardStats struct {
	UpcomingInterviews  int     `json:"upcoming_interviews"`
	CompletedMockTests  int     `json:"completed_mock_tests"`
	AverageMockScore    float64 `json:"average_mock_score"`
	ProfileCompleteness int     `json:"profile_completeness"` // percentage 0-100
}

type CandidatePortalInterview struct {
	ID              string     `json:"id" db:"id"`
	Title           string     `json:"title" db:"title"`
	CompanyID       string     `json:"company_id" db:"company_id"`
	CompanyName     string     `json:"company_name" db:"company_name"`
	JobID           string     `json:"job_id" db:"job_id"`
	JobTitle        string     `json:"job_title" db:"job_title"`
	Mode            string     `json:"mode" db:"mode"`
	Status          string     `json:"status" db:"status"`
	ScheduledAt     *time.Time `json:"scheduled_at" db:"scheduled_at"`
	InviteTokenHash string     `json:"-" db:"invite_token_hash"`
	JoinLink        string     `json:"join_link,omitempty" db:"-"`
}

type CandidatePortalProfile struct {
	UserID         string                   `json:"user_id"`
	FullName       string                   `json:"full_name"`
	Email          string                   `json:"email"`
	Phone          string                   `json:"phone"`
	AvatarURL      string                   `json:"avatar_url"`
	CVFileID       string                   `json:"cv_file_id"`
	CVUrl          string                   `json:"cv_url"`
	CVName         string                   `json:"cv_name"`
	CVParseStatus  string                   `json:"cv_parse_status,omitempty"`
	ParsedData     interface{}              `json:"parsed_data,omitempty"`
	InterviewScore *CandidateInterviewScore `json:"interview_score,omitempty"`
}

type CandidateInterviewScore struct {
	FinalScore float64  `json:"final_score"`
	Summary    string   `json:"summary"`
	Strengths  []string `json:"strengths,omitempty"`
	Weaknesses []string `json:"weaknesses,omitempty"`
	Advice     []string `json:"advice,omitempty"`
}

type CandidatePortalCVUpload struct {
	Message     string      `json:"message"`
	FileName    string      `json:"file_name"`
	CVUrl       string      `json:"cv_url"`
	CVFileID    string      `json:"cv_file_id"`
	ParseStatus string      `json:"parse_status"`
	ParsedData  interface{} `json:"parsed_data,omitempty"`
}

type CandidatePortalApplication struct {
	ID           string     `json:"id" db:"id"`
	JobID        string     `json:"job_id" db:"job_id"`
	JobTitle     string     `json:"job_title" db:"job_title"`
	CompanyID    string     `json:"company_id" db:"company_id"`
	CompanyName  string     `json:"company_name" db:"company_name"`
	Status       string     `json:"status" db:"status"`
	AppliedAt    *time.Time `json:"applied_at" db:"applied_at"`
	CVName       string     `json:"cv_name" db:"cv_name"`
	CVStorageKey string     `json:"-" db:"cv_storage_key"`
	CVUrl        string     `json:"cv_url,omitempty" db:"-"`
}

// CandidatePortalCVReview is AI feedback for improving a candidate CV (portal tool).
type CandidatePortalCVReview struct {
	Summary         string   `json:"summary"`
	Issues          []string `json:"issues"`
	Suggestions     []string `json:"suggestions"`
	MissingSections []string `json:"missing_sections"`
	Strengths       []string `json:"strengths"`
}
