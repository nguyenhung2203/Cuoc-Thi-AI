import logging

from app.llm.base import BaseLLMClient
from app.llm.gemini_client import GeminiClient
from app.llm.response_parser import ResponseParser
from app.models.generate_models import GenerateRequest, GenerateResponse

logger = logging.getLogger("ai-service.generate")


class GenerateService:
    """Core gateway logic: render is done by Go; here we just call the LLM
    and shape the response into the standard envelope."""

    def __init__(self, client: BaseLLMClient = None):
        self._client = client or GeminiClient()
        self._parser = ResponseParser()

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
