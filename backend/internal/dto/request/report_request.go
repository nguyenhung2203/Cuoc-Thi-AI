package request

type UpdateReportDecisionRequest struct {
	Decision string `json:"decision" validate:"required,oneof=strong_hire hire consider next_round reject insufficient_data"`
	Comment  string `json:"comment"`
}
