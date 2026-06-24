from fastapi import FastAPI

app = FastAPI(title="AI Interview Platform - AI Service")

@app.get("/health")
def health_check():
    return {"status": "healthy"}
