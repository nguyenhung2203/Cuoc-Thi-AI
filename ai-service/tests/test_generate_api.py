from fastapi.testclient import TestClient

from app.llm.base import BaseLLMClient, LLMResult
from app.main import app
from app.services.generate_service import GenerateService, get_generate_service


class StubLLMClient(BaseLLMClient):
    def generate(self, prompt, model="", temperature=0.2, max_tokens=2048):
        if "criterion" in prompt.lower():
            text = '{"score":3,"evidence":"Transcript evidence","confidence":0.5}'
        else:
            text = '{"questions":[{"question_text":"Describe a recent project"}]}'
        return LLMResult(
            text=text,
            tokens_in=17,
            tokens_out=9,
            model=model or "gemini-2.5-flash",
        )


app.dependency_overrides[get_generate_service] = lambda: GenerateService(StubLLMClient())
client = TestClient(app)


def test_health_reports_live_mode():
    response = client.get("/health")
    assert response.status_code == 200
    body = response.json()
    assert body["status"] == "healthy"
    assert body["mode"] == "live"


def test_generate_questions_shape_and_usage():
    response = client.post("/api/v1/generate", json={
        "prompt": "Sinh danh sách câu hỏi phỏng vấn (questions) cho vị trí backend.",
        "model": "gemini-2.5-flash",
    })
    assert response.status_code == 200
    body = response.json()
    assert isinstance(body["data"]["questions"], list)
    assert body["insufficient_data"] is False
    assert body["tokens_in"] == 17
    assert body["tokens_out"] == 9
    assert body["model"] == "gemini-2.5-flash"


def test_generate_score_lifts_envelope():
    response = client.post("/api/v1/generate", json={
        "prompt": "Chấm điểm câu trả lời theo criterion Technical Knowledge.",
    })
    assert response.status_code == 200
    body = response.json()
    assert body["data"]["score"] == 3
    assert body["evidence"] == "Transcript evidence"
    assert body["confidence"] == 0.5


def test_generate_rejects_empty_prompt():
    response = client.post("/api/v1/generate", json={"prompt": ""})
    assert response.status_code == 422
