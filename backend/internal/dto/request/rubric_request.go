package request

type CreateRubricRequest struct {
	JobID       *string                        `json:"job_id"`
	Name        string                         `json:"name" validate:"required,max=255"`
	Description *string                        `json:"description"`
	Criteria    []CreateRubricCriterionRequest `json:"criteria" validate:"required,min=1,dive"`
}

type CreateRubricCriterionRequest struct {
	Name         string  `json:"name" validate:"required,max=255"`
	Description  *string `json:"description"`
	Weight       float64 `json:"weight" validate:"required,min=0"`
	MinScore     int     `json:"min_score" validate:"required"`
	MaxScore     int     `json:"max_score" validate:"required"`
	ScoringGuide string  `json:"scoring_guide" validate:"required"`
	OrderIndex   int     `json:"order_index"`
}
