package request

type DecisionRequest struct {
	Decision string `json:"decision" validate:"required"`
	Comment  string `json:"comment"`
}
