import logging

from fastapi import APIRouter, HTTPException, Depends

from app.llm.base import LLMError
from app.models.generate_models import GenerateRequest, GenerateResponse
from app.services.generate_service import GenerateService, get_generate_service

logger = logging.getLogger("ai-service.api")

router = APIRouter()


@router.post("/api/v1/generate", response_model=GenerateResponse)
def generate(
    req: GenerateRequest,
    service: GenerateService = Depends(get_generate_service),
) -> GenerateResponse:
    """Single LLM generation endpoint used by the Go backend.

    Error mapping:
      - 400: invalid request / non-retryable input error
      - 429: rate limited / quota
      - 502/504: upstream LLM failure / timeout (retryable by Go)
    """
    try:
        return service.generate(req)
    except LLMError as e:
        logger.warning("LLM error (status=%s retryable=%s): %s",
                       e.status, e.retryable, e.message)
        raise HTTPException(status_code=e.status, detail=e.message)
    except Exception as e:  # pragma: no cover - defensive
        logger.exception("unexpected error in /generate")
        raise HTTPException(status_code=500, detail=str(e))
