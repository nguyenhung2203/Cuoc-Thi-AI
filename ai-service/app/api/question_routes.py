from fastapi import APIRouter, Depends, HTTPException

from app.llm.base import LLMError
from app.models.question_models import QuestionGenerationRequest, QuestionGenerationResponse
from app.services.generate_service import GenerateService, get_generate_service

router = APIRouter()


@router.post("/api/v1/questions/generate", response_model=QuestionGenerationResponse)
def generate_questions(
    req: QuestionGenerationRequest,
    service: GenerateService = Depends(get_generate_service),
) -> QuestionGenerationResponse:
    try:
        return service.generate_questions(req)
    except LLMError as e:
        raise HTTPException(status_code=e.status, detail=e.message) from e
    except ValueError as e:
        raise HTTPException(status_code=502, detail=str(e)) from e
    except Exception as e:  # pragma: no cover
        raise HTTPException(status_code=500, detail="question generation failed") from e


@router.post("/api/v1/questions/follow-up", response_model=QuestionGenerationResponse)
def generate_follow_up(
    req: QuestionGenerationRequest,
    service: GenerateService = Depends(get_generate_service),
) -> QuestionGenerationResponse:
    """Generate one evidence-based follow-up without exposing a second LLM contract."""
    try:
        return service.generate_follow_up(req)
    except LLMError as e:
        raise HTTPException(status_code=e.status, detail=e.message) from e
    except ValueError as e:
        raise HTTPException(status_code=502, detail=str(e)) from e
    except Exception as e:  # pragma: no cover
        raise HTTPException(status_code=500, detail="follow-up generation failed") from e
