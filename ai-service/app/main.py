from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import os
import requests
import json

app = FastAPI(title="AI Interview Platform - AI Service")

class AIPayload(BaseModel):
    model: str = "gemini-1.5-flash"
    prompt: str
    temperature: float = 0.7
    max_tokens: int = 1000

@app.get("/health")
def health_check():
    return {"status": "healthy"}

@app.post("/api/v1/generate")
def generate(payload: AIPayload):
    api_key = os.environ.get("GEMINI_API_KEY")
    if not api_key:
        raise HTTPException(status_code=500, detail="GEMINI_API_KEY is not set")
    
    url = f"https://generativelanguage.googleapis.com/v1beta/models/{payload.model}:generateContent?key={api_key}"
    
    req_body = {
        "contents": [{"parts": [{"text": payload.prompt}]}],
        "generationConfig": {
            "temperature": payload.temperature,
            "maxOutputTokens": payload.max_tokens,
        }
    }
    
    resp = requests.post(url, json=req_body)
    if resp.status_code != 200:
        raise HTTPException(status_code=500, detail=resp.text)
        
    data = resp.json()
    try:
        text = data["candidates"][0]["content"]["parts"][0]["text"]
        # Try to parse the text as JSON, usually Gemini returns a markdown code block like ```json ... ```
        text = text.replace("```json", "").replace("```", "").strip()
        parsed_data = json.loads(text)
    except Exception as e:
        # If parsing fails, just return it as a raw string or dictionary
        parsed_data = {"raw_text": text} if 'text' in locals() else {}

    return {
        "data": parsed_data,
        "evidence": "",
        "confidence": 0.9,
        "insufficient_data": False
    }
