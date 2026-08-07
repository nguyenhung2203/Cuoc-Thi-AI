package ai

// JDAnalysisResult matches the output JSON of analyze_jd template
type JDAnalysisResult struct {
	Summary             string              `json:"summary"`
	RequiredSkills      []string            `json:"required_skills"`
	NiceToHaveSkills    []string            `json:"nice_to_have_skills"`
	SeniorityAssessment string              `json:"seniority_assessment"`
	MissingInformation  []string            `json:"missing_information"`
	InterviewFocusAreas []string            `json:"interview_focus_areas"`
	SuggestedRubric     []SuggestedRubric   `json:"suggested_rubric"`
	SuggestedQuestions  []SuggestedQuestion `json:"suggested_questions"`
}

type SuggestedRubric struct {
	Name        string `json:"name"`
	Weight      int    `json:"weight"`
	Description string `json:"description"`
}

type SuggestedQuestion struct {
	ID               string   `json:"id,omitempty"`
	QuestionText     string   `json:"question_text"`
	QuestionType     string   `json:"question_type,omitempty"`
	Category         string   `json:"category,omitempty"`
	TargetSkill      string   `json:"target_skill,omitempty"`
	SkillTags        []string `json:"skill_tags,omitempty"`
	Difficulty       string   `json:"difficulty"`
	ExpectedSignals  []string `json:"expected_signals"`
	FollowUpPrompts  []string `json:"follow_up_prompts,omitempty"`
	TimeboxMinutes   int      `json:"timebox_minutes,omitempty"`
	EvidenceRequired bool     `json:"evidence_required,omitempty"`
	WhyAsk           string   `json:"why_ask,omitempty"`
	RedFlags         []string `json:"red_flags,omitempty"`
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

type QuestionGenerationResult struct {
	Status         string              `json:"status,omitempty"`
	Mode           string              `json:"mode,omitempty"`
	Level          string              `json:"level,omitempty"`
	Questions      []SuggestedQuestion `json:"questions"`
	Coverage       []string            `json:"coverage,omitempty"`
	Warnings       []string            `json:"warnings,omitempty"`
	Confidence     float64             `json:"confidence,omitempty"`
	PromptVersion  string              `json:"prompt_version,omitempty"`
	RequestedCount int                 `json:"requested_count,omitempty"`
	ActualCount    int                 `json:"actual_count,omitempty"`
}

// ReportGenerationResult matches the output JSON of generate_report template
type ReportGenerationResult struct {
	Summary            string   `json:"summary"`
	FinalScore         float64  `json:"final_score"`
	Recommendation     string   `json:"recommendation"` // strong_hire, hire, consider, next_round, reject, insufficient_data
	Strengths          []string `json:"strengths"`
	Weaknesses         []string `json:"weaknesses"`
	Risks              []string `json:"risks"`
	EvidenceJSON       any      `json:"evidence_json"` // Can be a map mapping criterion to evidence
	AIReasoningSummary string   `json:"ai_reasoning_summary"`
	ImprovementAdvice  []string `json:"improvement_advice,omitempty"`
	CommunicationScore float64  `json:"communication_score,omitempty"`
	ToneScore          float64  `json:"tone_score,omitempty"`
	PersonalityScore   float64  `json:"personality_score,omitempty"`
}
