package ai

// JDAnalysisResult matches the output JSON of analyze_jd template
type JDAnalysisResult struct {
	Summary              string              `json:"summary"`
	RequiredSkills       []string            `json:"required_skills"`
	NiceToHaveSkills     []string            `json:"nice_to_have_skills"`
	SeniorityAssessment  string              `json:"seniority_assessment"`
	MissingInformation   []string            `json:"missing_information"`
	InterviewFocusAreas  []string            `json:"interview_focus_areas"`
	SuggestedRubric      []SuggestedRubric   `json:"suggested_rubric"`
	SuggestedQuestions   []SuggestedQuestion `json:"suggested_questions"`
}

type SuggestedRubric struct {
	Name        string `json:"name"`
	Weight      int    `json:"weight"`
	Description string `json:"description"`
}

type SuggestedQuestion struct {
	QuestionText    string   `json:"question_text"`
	QuestionType    string   `json:"question_type"`
	TargetSkill     string   `json:"target_skill"`
	Difficulty      string   `json:"difficulty"`
	ExpectedSignals []string `json:"expected_signals"`
	WhyAsk          string   `json:"why_ask,omitempty"`
	RedFlags        []string `json:"red_flags,omitempty"`
}

// CVAnalysisResult matches the output JSON of analyze_cv template
type CVAnalysisResult struct {
	Summary                 string             `json:"summary"`
	Skills                  []CVSkill          `json:"skills"`
	ExperienceYearsEstimate int                `json:"experience_years_estimate"`
	WorkExperience          []CVWorkExperience `json:"work_experience"`
	Projects                []CVProject        `json:"projects"`
	Education               []string           `json:"education"`
	PotentialStrengths      []string           `json:"potential_strengths"`
	PotentialConcerns       []string           `json:"potential_concerns"`
	QuestionsToVerify       []string           `json:"questions_to_verify"`
	Confidence              float64            `json:"confidence"`
}

type CVSkill struct {
	Name      string `json:"name"`
	LevelHint string `json:"level_hint"`
	Evidence  string `json:"evidence"`
}

type CVWorkExperience struct {
	Company    string   `json:"company"`
	Role       string   `json:"role"`
	Duration   string   `json:"duration"`
	Highlights []string `json:"highlights"`
}

type CVProject struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	TechStack   []string `json:"tech_stack"`
}

// QuestionGenerationResult matches the output JSON of generate_questions template
type QuestionGenerationResult struct {
	Questions []SuggestedQuestion `json:"questions"`
}
