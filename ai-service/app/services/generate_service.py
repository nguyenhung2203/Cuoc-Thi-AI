import logging

from app.llm.base import BaseLLMClient
from app.llm.gemini_client import GeminiClient
from app.llm.response_parser import ResponseParser
from app.models.generate_models import GenerateRequest, GenerateResponse
from app.models.question_models import QuestionGenerationRequest, QuestionGenerationResponse
from app.prompts.question_generation_v2 import build_prompt
from app.utils.question_validator import validate_questions

logger = logging.getLogger("ai-service.generate")


class GenerateService:
    """Core gateway logic: render is done by Go; here we just call the LLM
    and shape the response into the standard envelope."""

    def __init__(self, client: BaseLLMClient = None):
        self._client = client or GeminiClient()
        self._parser = ResponseParser()

    def generate_questions(self, req: QuestionGenerationRequest) -> QuestionGenerationResponse:
        result = self._client.generate(
            prompt=build_prompt(req),
            model=req.model,
            temperature=req.temperature,
            max_tokens=req.max_tokens,
        )
        try:
            data = self._parser.parse_json(result.text)
            return validate_questions(data, req)
        except ValueError as e:
            logger.warning("Question generation validation failed: %s", e)
            raise ValueError(str(e)) from e

    def generate_follow_up(self, req: QuestionGenerationRequest) -> QuestionGenerationResponse:
        follow_up_req = req.model_copy(update={"question_count": 1})
        if not follow_up_req.recent_answer.strip():
            follow_up_req = follow_up_req.model_copy(update={
                "recent_answer": "Không có câu trả lời; hãy hỏi một câu clarification ngắn."
            })
        return self.generate_questions(follow_up_req)

    def generate(self, req: GenerateRequest) -> GenerateResponse:
        result = self._client.generate(
            prompt=req.prompt,
            model=req.model,
            temperature=req.temperature,
            max_tokens=req.max_tokens,
        )

        try:
            data = self._parser.parse_json(result.text)
        except ValueError as e:
            logger.warning("LLM returned non-JSON output: %s", e)
            raise ValueError("Gemini returned malformed JSON") from e

        data, evidence, confidence, insufficient = self._parser.build_envelope(data)
        return GenerateResponse(
            data=data,
            evidence=evidence,
            confidence=confidence,
            insufficient_data=insufficient,
            model=result.model,
            tokens_in=result.tokens_in,
            tokens_out=result.tokens_out,
        )


# Shared singleton (FastAPI reuses across requests).
_service: GenerateService = None


def get_generate_service() -> GenerateService:
    global _service
    if _service is None:
        _service = GenerateService()
    return _service
