from fastapi import APIRouter

from app.config import config

router = APIRouter()


@router.get("/health")
def health_check():
    """Liveness/readiness for the configured live Gemini gateway."""
    config.validate()
    return {
        "status": "healthy",
        "mode": "live",
        "default_model": config.DEFAULT_MODEL,
    }
