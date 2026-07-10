import logging

from fastapi import FastAPI

from app.api import generate_routes, health_routes

logging.basicConfig(level=logging.INFO)

app = FastAPI(title="AI Interview Platform - AI Service", version="1.0.0")

# The Go backend renders prompts from DB templates and calls a single
# generation endpoint. Health reports live/mock mode.
app.include_router(health_routes.router, tags=["health"])
app.include_router(generate_routes.router, tags=["generate"])
