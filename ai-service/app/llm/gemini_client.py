import json
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
    """Google Gemini client via the current `google-genai` SDK.

    (Migrated off the deprecated `google-generativeai` package.) Configured
    lazily so the module imports cleanly even when the SDK or a key is absent
    (mock mode). The underlying client is thread-safe and reused across requests.
    """

    def __init__(self):
        self._client = None
        if not config.use_mock():
            self._configure()

    def _configure(self):
        try:
            from google import genai  # google-genai (new SDK)
        except ImportError as e:  # pragma: no cover - depends on env
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

        if config.use_mock():
            return self._mock_generate(prompt, model_name)

        if self._client is None:
            self._configure()

        try:
            resp = self._client.models.generate_content(
                model=model_name,
                contents=prompt,
                config={
                    "temperature": temperature,
                    "max_output_tokens": max_tokens,
                    "response_mime_type": "application/json",
                },
            )
        except Exception as e:  # SDK raises a variety of exception types
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
            text=text, tokens_in=tokens_in, tokens_out=tokens_out, model=model_name
        )

    @staticmethod
    def _extract_text(resp) -> str:
        # Newer SDKs expose .text; fall back to walking candidates/parts.
        text = getattr(resp, "text", None)
        if text:
            return text
        try:
            parts = resp.candidates[0].content.parts
            return "".join(getattr(p, "text", "") for p in parts)
        except Exception:
            raise LLMError("empty response from Gemini", retryable=True, status=502)

    @staticmethod
    def _mock_generate(prompt: str, model_name: str) -> LLMResult:
        """Deterministic, schema-agnostic JSON stub used when no key is set.

        The prompt already tells the model what JSON to emit, so for mock mode
        we emit a small generic object plus envelope hints. The Go layer
        unmarshals only the fields it needs, tolerating extras. We branch on a
        few keywords so common flows return plausible shapes.
        """
        low = prompt.lower()
        if "score" in low or "chấm điểm" in low or "criterion" in low:
            data = {
                "score": 3,
                "ai_comment": "[MOCK] Câu trả lời ở mức đạt yêu cầu.",
                "evidence": "[MOCK] Trích dẫn từ transcript.",
                "confidence": 0.5,
            }
        elif "report" in low or "báo cáo" in low:
            data = {
                "summary": "[MOCK] Tóm tắt buổi phỏng vấn.",
                "final_score": 3.0,
                "recommendation": "consider",
                "strengths": ["[MOCK] Điểm mạnh mẫu"],
                "weaknesses": ["[MOCK] Điểm yếu mẫu"],
                "risks": [],
                "reasoning": "[MOCK] Lý do đề xuất.",
            }
        elif "question" in low or "câu hỏi" in low:
            data = {
                "questions": [
                    {
                        "question_text": "[MOCK] Bạn hãy giới thiệu một dự án gần đây?",
                        "question_type": "experience",
                        "target_skill": "communication",
                        "difficulty": "middle",
                        "why_ask": "[MOCK]",
                        "expected_signals": ["rõ ràng", "có số liệu"],
                        "red_flags": [],
                    }
                ]
            }
        else:
            data = {
                "summary": "[MOCK] Kết quả phân tích mẫu.",
                "confidence": 0.5,
            }
        return LLMResult(text=json.dumps(data, ensure_ascii=False), model=model_name)
