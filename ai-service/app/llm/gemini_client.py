import logging

from app.config import config
from app.llm.base import (
    BaseLLMClient,
    LLMResult,
    LLMError,
    LLMTimeoutError,
    LLMInvalidRequestError,
)

logger = logging.getLogger("ai-service.gemini")


class GeminiClient(BaseLLMClient):
    """Google Gemini client via the current ``google-genai`` SDK."""

    def __init__(self):
        self._client = None
        self._configure()

    def _configure(self):
        try:
            from google import genai
        except ImportError as e:  # pragma: no cover - depends on environment
            raise LLMError(
                f"google-genai not installed: {e}", retryable=False, status=500
            )
        self._client = genai.Client(api_key=config.GEMINI_API_KEY)

    def generate(
        self,
        prompt: str,
        model: str = "",
        temperature: float = 0.2,
        max_tokens: int = 2048,
    ) -> LLMResult:
        if not prompt or not prompt.strip():
            raise LLMInvalidRequestError("prompt must not be empty")

        model_name = model or config.DEFAULT_MODEL

        try:
            resp = self._client.models.generate_content(
                model=model_name,
                contents=prompt,
                config={
                    "temperature": temperature,
                    "max_output_tokens": max_tokens,
                    "response_mime_type": "application/json",
                    "thinking_config": {"thinking_budget": 0},
                },
            )
        except Exception as e:
            msg = str(e)
            low = msg.lower()
            if "timeout" in low or "deadline" in low:
                raise LLMTimeoutError(msg)
            if "api key" in low or "invalid" in low or "permission" in low:
                raise LLMInvalidRequestError(msg)
            if "quota" in low or "rate" in low or "429" in low or "resource_exhausted" in low:
                raise LLMError(msg, retryable=True, status=429)
            raise LLMError(msg, retryable=True, status=502)

        text = self._extract_text(resp)
        usage = getattr(resp, "usage_metadata", None)
        tokens_in = getattr(usage, "prompt_token_count", 0) if usage else 0
        tokens_out = getattr(usage, "candidates_token_count", 0) if usage else 0
        return LLMResult(
            text=text,
            tokens_in=tokens_in or 0,
            tokens_out=tokens_out or 0,
            model=model_name,
        )

    @staticmethod
    def _extract_text(resp) -> str:
        text = getattr(resp, "text", None)
        if text:
            return text
        try:
            parts = resp.candidates[0].content.parts
            text = "".join(getattr(part, "text", "") for part in parts)
            if text:
                return text
        except Exception:
            pass
        raise LLMError("empty response from Gemini", retryable=True, status=502)
