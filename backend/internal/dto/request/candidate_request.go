package request

type CreateCandidateRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Email    string `json:"email"     validate:"required,email"`
	Phone    string `json:"phone"`
	Source   string `json:"source"    validate:"omitempty,oneof=linkedin referral import manual"`
	JobID    string `json:"job_id"`   // optional — if provided, create job_candidates record
}

type UpdateCandidateRequest struct {
	FullName *string  `json:"full_name" validate:"omitempty,min=2,max=255"`
	Email    *string  `json:"email"     validate:"omitempty,email"`
	Phone    *string  `json:"phone"`
	Source   *string  `json:"source"    validate:"omitempty,oneof=linkedin referral import manual"`
	Status   *string  `json:"status"    validate:"omitempty,oneof=new screening invited interviewing completed passed rejected talent_pool"`
	Tags     []string `json:"tags"`
}

type ListCandidatesQuery struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Keyword string `json:"keyword"`
}

type AssignCandidateRequest struct {
	PipelineStatus string `json:"pipeline_status" validate:"required,oneof=new screening invited interviewing completed passed rejected talent_pool"`
}
