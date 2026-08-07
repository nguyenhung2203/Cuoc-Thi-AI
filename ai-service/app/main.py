import logging

from fastapi import FastAPI

from app.api import generate_routes, health_routes, question_routes
from app.config import config

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("ai-service")

# Import-time validation is intentional: an incorrectly configured serving
# process must never accept requests and return fabricated AI output.
config.validate()

app = FastAPI(title="AI Interview Platform - AI Service", version="1.0.0")
app.include_router(health_routes.router, tags=["health"])
app.include_router(generate_routes.router, tags=["generate"])
app.include_router(question_routes.router, tags=["questions"])
