package request

type CreateJobRequest struct {
	Title          string  `json:"title"           validate:"required,min=2,max=255"`
	Department     string  `json:"department"`
	Level          string  `json:"level"           validate:"omitempty,oneof=intern junior middle senior lead"`
	Location       string  `json:"location"`
	EmploymentType string  `json:"employment_type" validate:"omitempty,oneof=full_time part_time contract intern"`
	SalaryMin      float64 `json:"salary_min"      validate:"omitempty,min=0"`
	SalaryMax      float64 `json:"salary_max"      validate:"omitempty,min=0"`
	Currency       string  `json:"currency"`
	Description    string  `json:"description"     validate:"required,min=10"`
	Requirements   string  `json:"requirements"`
	Benefits       string  `json:"benefits"`
	Status         string  `json:"status"          validate:"required,oneof=draft open paused closed"`
}

type UpdateJobRequest struct {
	Title          *string  `json:"title"           validate:"omitempty,min=2,max=255"`
	Department     *string  `json:"department"`
	Level          *string  `json:"level"           validate:"omitempty,oneof=intern junior middle senior lead"`
	Location       *string  `json:"location"`
	EmploymentType *string  `json:"employment_type" validate:"omitempty,oneof=full_time part_time contract intern"`
	SalaryMin      *float64 `json:"salary_min"      validate:"omitempty,min=0"`
	SalaryMax      *float64 `json:"salary_max"      validate:"omitempty,min=0"`
	Currency       *string  `json:"currency"`
	Description    *string  `json:"description"     validate:"omitempty,min=10"`
	Requirements   *string  `json:"requirements"`
	Benefits       *string  `json:"benefits"`
	Status         *string  `json:"status"          validate:"omitempty,oneof=draft open paused closed"`
}

type ListJobsQuery struct {
	Status  string `json:"status"`
	Keyword string `json:"keyword"`
}
