from typing import Any, Literal

from pydantic import BaseModel, Field, field_validator

from app.utils.level_mapping import normalize_level, normalize_mode


class QuestionGenerationRequest(BaseModel):
    job_title: str = Field(..., min_length=1, max_length=200)
    job_description: str = Field(default="", max_length=12000)
    requirements: list[str] = Field(default_factory=list, max_length=50)
    level: str = "unknown"
    mode: str = "real"
    language: str = Field(default="vi", min_length=2, max_length=20)
    skill_tags: list[str] = Field(default_factory=list, max_length=50)
    rubric: dict[str, Any] = Field(default_factory=dict)
    question_count: int = Field(default=5, ge=1, le=20)
    previous_questions: list[str] = Field(default_factory=list, max_length=50)
    recent_answer: str = Field(default="", max_length=12000)
    model: str = Field(default="", max_length=100)
    temperature: float = Field(default=0.2, ge=0, le=2)
    max_tokens: int = Field(default=4096, ge=1, le=32768)

    @field_validator("level")
    @classmethod
    def valid_level(cls, value: str) -> str:
        return normalize_level(value)

    @field_validator("mode")
    @classmethod
    def valid_mode(cls, value: str) -> str:
        return normalize_mode(value)


class QuestionItem(BaseModel):
    id: str = ""
    question_text: str = Field(..., min_length=1)
    category: Literal["technical", "behavioral", "problem_solving", "communication"] = "technical"
    difficulty: Literal["basic", "intermediate", "advanced"] = "intermediate"
    skill_tags: list[str] = Field(default_factory=list)
    expected_signals: list[str] = Field(default_factory=list)
    follow_up_prompts: list[str] = Field(default_factory=list)
    timebox_minutes: int = Field(default=5, ge=1, le=60)
    evidence_required: bool = True


class QuestionGenerationResponse(BaseModel):
    status: Literal["complete", "degraded"] = "complete"
    mode: str
    level: str
    questions: list[QuestionItem]
    coverage: list[str] = Field(default_factory=list)
    warnings: list[str] = Field(default_factory=list)
    confidence: float = Field(default=0, ge=0, le=1)
    prompt_version: str = "question_generation_v2"
    requested_count: int = 0
    actual_count: int = 0
