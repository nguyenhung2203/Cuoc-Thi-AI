import re

from app.models.question_models import QuestionGenerationResponse, QuestionItem
from app.utils.level_mapping import LEVEL_DIFFICULTIES

_ALLOWED_CATEGORIES = {"technical", "behavioral", "problem_solving", "communication"}
_ALLOWED_DIFFICULTIES = {"basic", "intermediate", "advanced"}
_SKILL_ALIASES = {"vue.js": "vue", "node.js": "node", "javascript": "js", "typescript": "ts"}


def _normalize_skill(value: str) -> str:
    value = re.sub(r"[^a-z0-9+#. ]+", " ", value.lower()).strip()
    return _SKILL_ALIASES.get(value, value)


def _scope_tokens(request):
    values = list(request.skill_tags) + list(request.requirements)
    phrases = {_normalize_skill(value) for value in values if value.strip()}
    tokens = {token for phrase in phrases for token in phrase.split() if len(token) > 1}
    return phrases, tokens


def _skill_in_scope(tag: str, phrases: set[str], tokens: set[str]) -> bool:
    normalized = _normalize_skill(tag)
    if normalized in phrases:
        return True
    return any(token in tokens for token in normalized.split())


def validate_questions(data, request) -> QuestionGenerationResponse:
    if not isinstance(data, dict) or not isinstance(data.get("questions"), list):
        raise ValueError("LLM response must contain questions array")

    warnings = list(data.get("warnings") or [])
    questions = []
    scope_phrases, scope_tokens = _scope_tokens(request)
    degraded = False
    seen_questions = set()
    previous_questions = {
        re.sub(r"\s+", " ", question.strip().lower())
        for question in request.previous_questions
        if question.strip()
    }
    require_skill = bool(scope_tokens)
    raw_questions = data["questions"]
    if len(raw_questions) > request.question_count:
        warnings.append("question_count_truncated")
        degraded = True

    for index, raw in enumerate(raw_questions[:request.question_count], 1):
        if not isinstance(raw, dict):
            warnings.append(f"question_{index}_invalid_shape")
            degraded = True
            continue
        try:
            item = QuestionItem.model_validate({"id": f"q-{index}", **raw})
        except Exception:
            warnings.append(f"question_{index}_invalid_schema")
            degraded = True
            continue
        normalized_question = re.sub(r"\s+", " ", item.question_text.strip().lower())
        if normalized_question in seen_questions or normalized_question in previous_questions:
            warnings.append(f"question_{index}_duplicate")
            degraded = True
            continue
        seen_questions.add(normalized_question)
        if item.difficulty not in LEVEL_DIFFICULTIES[request.level]:
            warnings.append(f"question_{index}_difficulty_mismatch")
            degraded = True
            continue
        if require_skill and not item.skill_tags:
            warnings.append(f"question_{index}_missing_skill_tags")
            degraded = True
            continue
        if scope_tokens and item.skill_tags and not any(
            _skill_in_scope(tag, scope_phrases, scope_tokens) for tag in item.skill_tags
        ):
            warnings.append(f"question_{index}_skill_out_of_scope")
            degraded = True
            continue
        if not item.expected_signals:
            warnings.append(f"question_{index}_missing_expected_signals")
            degraded = True
        questions.append(item)

    if not questions:
        raise ValueError("LLM returned no valid questions")
    if len(questions) < request.question_count:
        warnings.append("question_count_incomplete")
        degraded = True
    if request.level == "unknown":
        warnings.append("job_level_unknown")
        degraded = True
    if not request.job_description and not request.requirements:
        warnings.append("job_context_incomplete")
        degraded = True

    try:
        confidence = max(0.0, min(1.0, float(data.get("confidence", 0.5))))
    except (TypeError, ValueError):
        confidence = 0.3
    if degraded:
        confidence = min(confidence, 0.5)

    return QuestionGenerationResponse(
        status="degraded" if degraded else "complete",
        mode=request.mode,
        level=request.level,
        questions=questions,
        coverage=list(data.get("coverage") or []),
        warnings=warnings,
        confidence=confidence,
        requested_count=request.question_count,
        actual_count=len(questions),
    )
