import logging

from fastapi import FastAPI

from app.api import generate_routes, health_routes
from app.config import config

logging.basicConfig(level=logging.INFO)

logger = logging.getLogger("ai-service")

# Fail fast: in production a real Gemini key is mandatory. This refuses to boot
# in mock mode rather than silently serving [MOCK] responses to real users.
config.validate()

if config.use_mock():
    logger.warning(
        "AI service starting in MOCK mode (APP_ENV=%s) — responses are deterministic stubs, not real Gemini output",
        config.APP_ENV,
    )

app = FastAPI(title="AI Interview Platform - AI Service", version="1.0.0")

# The Go backend renders prompts from DB templates and calls a single
# generation endpoint. Health reports live/mock mode.
app.include_router(health_routes.router, tags=["health"])
app.include_router(generate_routes.router, tags=["generate"])
