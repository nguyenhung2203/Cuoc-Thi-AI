from typing import Any, Optional
from pydantic import BaseModel, Field


class GenerateRequest(BaseModel):
    """Payload sent by the Go backend to POST /api/v1/generate."""
    prompt: str = Field(..., min_length=1)
    model: str = ""
    temperature: float = 0.2
    max_tokens: int = Field(default=2048, ge=1, le=32768)


class GenerateResponse(BaseModel):
    """Standard envelope returned to the Go backend."""
    data: Any = None
    evidence: str = ""
    confidence: float = 0.0
    insufficient_data: bool = False
    model: Optional[str] = None
    tokens_in: Optional[int] = None
    tokens_out: Optional[int] = None
