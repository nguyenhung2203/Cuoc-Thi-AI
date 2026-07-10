from abc import ABC, abstractmethod
from dataclasses import dataclass


@dataclass
class LLMResult:
    """Raw text returned by an LLM plus optional usage metadata."""
    text: str
    tokens_in: int = 0
    tokens_out: int = 0
    model: str = ""


class BaseLLMClient(ABC):
    """Abstract LLM client. Implementations must be safe to reuse across
    requests (FastAPI shares one instance)."""

    @abstractmethod
    def generate(
        self,
        prompt: str,
        model: str = "",
        temperature: float = 0.2,
        max_tokens: int = 2048,
    ) -> LLMResult:
        """Send a single prompt and return the raw completion text.

        Must raise LLMError (or subclass) on any upstream failure so the
        API layer can map it to the right HTTP status.
        """
        raise NotImplementedError


class LLMError(Exception):
    """Base class for LLM failures."""

    def __init__(self, message: str, *, retryable: bool = True, status: int = 502):
        super().__init__(message)
        self.message = message
        self.retryable = retryable
        self.status = status


class LLMTimeoutError(LLMError):
    def __init__(self, message: str = "LLM request timed out"):
        super().__init__(message, retryable=True, status=504)


class LLMInvalidRequestError(LLMError):
    def __init__(self, message: str = "invalid LLM request"):
        super().__init__(message, retryable=False, status=400)
