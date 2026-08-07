package request

type CreateQuestionBankRequest struct {
	JobID           *string  `json:"job_id"`
	QuestionText    string   `json:"question_text" validate:"required"`
	QuestionType    string   `json:"question_type" validate:"required"`
	SkillTags       []string `json:"skill_tags"`
	Level           string   `json:"level"`
	ExpectedSignals []string `json:"expected_signals"`
}

type UpdateQuestionBankRequest struct {
	QuestionText    *string   `json:"question_text"`
	QuestionType    *string   `json:"question_type"`
	SkillTags       *[]string `json:"skill_tags"`
	Level           *string   `json:"level"`
	ExpectedSignals *[]string `json:"expected_signals"`
}

type GenerateQuestionsRequest struct {
	QuestionTypes []string `json:"question_types"`
	Level         string   `json:"level"`
	Mode          string   `json:"mode"`
	Count         int      `json:"count"`
}
