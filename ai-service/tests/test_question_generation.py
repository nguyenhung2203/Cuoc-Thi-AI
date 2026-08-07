from fastapi.testclient import TestClient

from app.llm.base import BaseLLMClient, LLMResult
from app.main import app
from app.services.generate_service import GenerateService, get_generate_service


class QuestionStub(BaseLLMClient):
    response = '{"questions": [{"question_text": "Explain Vue", "category": "technical", "difficulty": "intermediate", "skill_tags": ["Vue"], "expected_signals": ["concept"]}], "confidence": 0.9}'

    def generate(self, prompt, model="", temperature=0.2, max_tokens=2048):
        return LLMResult(self.response, 1, 1, model or "stub")


def client_for(stub=None):
    app.dependency_overrides[get_generate_service] = lambda: GenerateService(stub or QuestionStub())
    return TestClient(app)


def payload(**kwargs):
    return {"job_title": "Frontend Developer", "job_description": "Build Vue web applications", "level": "fresher", "mode": "real", "skill_tags": ["Vue"], "question_count": 1, **kwargs}


def test_question_generation_contract():
    response = client_for().post("/api/v1/questions/generate", json=payload())
    assert response.status_code == 200
    body = response.json()
    assert body["status"] == "complete"
    assert body["level"] == "fresher"
    assert body["actual_count"] == 1


def test_mid_normalizes_to_middle():
    response = client_for().post("/api/v1/questions/generate", json=payload(level="mid"))
    assert response.json()["level"] == "middle"


def test_follow_up_forces_one_question():
    response = client_for().post("/api/v1/questions/follow-up", json=payload(question_count=5, recent_answer="I used Vue components."))
    assert response.status_code == 200
    body = response.json()
    assert body["requested_count"] == 1
    assert body["actual_count"] == 1


    original = QuestionStub.response
    QuestionStub.response = '{"questions": [null, "bad"]}'
    response = client_for().post("/api/v1/questions/generate", json=payload())
    assert response.status_code == 502
    QuestionStub.response = original
