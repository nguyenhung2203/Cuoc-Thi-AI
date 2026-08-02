import os


class Config:
    """Configuration for the live Gemini gateway.

    A serving process always requires a real provider credential. Test doubles
    must be injected into GenerateService by tests; runtime mock mode is not
    supported.
    """

    APP_ENV: str = os.getenv("APP_ENV", "development").strip().lower()
    GEMINI_API_KEY: str = os.getenv("GEMINI_API_KEY", "").strip()
    DEFAULT_MODEL: str = os.getenv("AI_DEFAULT_MODEL", "gemini-3.1-flash-lite")
    DEFAULT_TEMPERATURE: float = float(os.getenv("AI_DEFAULT_TEMPERATURE", "0.2"))
    DEFAULT_MAX_TOKENS: int = int(os.getenv("AI_DEFAULT_MAX_TOKENS", "2048"))
    LLM_TIMEOUT_SECONDS: float = float(os.getenv("AI_LLM_TIMEOUT_SECONDS", "60"))

    @classmethod
    def key_placeholder(cls) -> bool:
        return cls.GEMINI_API_KEY.lower() in {"", "your-gemini-api-key", "changeme", "placeholder"}

    @classmethod
    def validate(cls) -> None:
        if cls.key_placeholder():
            raise RuntimeError(
                "GEMINI_API_KEY is missing or a placeholder; the AI service requires a live Gemini credential"
            )


config = Config()
config.validate()
