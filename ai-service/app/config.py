import os


def _get_bool(name: str, default: bool = False) -> bool:
    val = os.getenv(name)
    if val is None:
        return default
    return val.strip().lower() in ("1", "true", "yes", "on")


class Config:
    """Central configuration for the AI service (thin Gemini gateway).

    The Go backend renders prompts from DB templates and calls
    POST /api/v1/generate with {model, prompt, temperature, max_tokens}.
    This service only talks to the LLM and returns the standard envelope
    {data, evidence, confidence, insufficient_data}.
    """

    GEMINI_API_KEY: str = os.getenv("GEMINI_API_KEY", "").strip()

    # Default model used when the caller does not specify one.
    DEFAULT_MODEL: str = os.getenv("AI_DEFAULT_MODEL", "gemini-1.5-flash")

    # Generation defaults (overridable per-request).
    DEFAULT_TEMPERATURE: float = float(os.getenv("AI_DEFAULT_TEMPERATURE", "0.2"))
    DEFAULT_MAX_TOKENS: int = int(os.getenv("AI_DEFAULT_MAX_TOKENS", "2048"))

    # Request timeout to the upstream LLM, seconds.
    LLM_TIMEOUT_SECONDS: float = float(os.getenv("AI_LLM_TIMEOUT_SECONDS", "60"))

    # MOCK mode: when true (or when no API key is present) the service returns
    # a deterministic, schema-valid stub instead of calling Gemini. This keeps
    # the full pipeline testable in CI without a real key. With a real key and
    # AI_MOCK unset, live Gemini is used.
    MOCK_MODE: bool = _get_bool("AI_MOCK", default=False)

    @classmethod
    def use_mock(cls) -> bool:
        return cls.MOCK_MODE or cls.key_placeholder()

    @classmethod
    def key_placeholder(cls) -> bool:
        """True when the configured key is empty or an obvious placeholder."""
        return cls.GEMINI_API_KEY in ("", "your-gemini-api-key", "changeme")


config = Config()
