from fastapi import APIRouter

from app.config import config

router = APIRouter()


@router.get("/health")
def health_check():
    """Liveness + basic readiness. Reports whether the service will call a
    real LLM or fall back to mock mode (no/placeholder key)."""
    return {
        "status": "healthy",
        "mode": "mock" if config.use_mock() else "live",
        "default_model": config.DEFAULT_MODEL,
    }
