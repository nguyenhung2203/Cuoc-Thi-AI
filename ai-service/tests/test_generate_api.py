from fastapi.testclient import TestClient

from app.main import app

client = TestClient(app)


def test_health_reports_mock_mode():
    r = client.get("/health")
    assert r.status_code == 200
    body = r.json()
    assert body["status"] == "healthy"
    assert body["mode"] == "mock"


def test_generate_questions_shape():
    r = client.post("/api/v1/generate", json={
        "prompt": "Sinh danh sách câu hỏi phỏng vấn (questions) cho vị trí backend.",
        "model": "gemini-1.5-flash",
    })
    assert r.status_code == 200
    body = r.json()
    assert "data" in body
    assert "questions" in body["data"]
    assert isinstance(body["data"]["questions"], list)
    assert body["insufficient_data"] is False


def test_generate_score_lifts_envelope():
    r = client.post("/api/v1/generate", json={
        "prompt": "Chấm điểm câu trả lời theo criterion Technical Knowledge.",
    })
    assert r.status_code == 200
    body = r.json()
    assert body["data"]["score"] == 3
    # evidence/confidence lifted from the object into the envelope
    assert body["evidence"] != ""
    assert body["confidence"] == 0.5


def test_generate_rejects_empty_prompt():
    r = client.post("/api/v1/generate", json={"prompt": ""})
    # Pydantic min_length=1 -> 422 validation error
    assert r.status_code == 422
