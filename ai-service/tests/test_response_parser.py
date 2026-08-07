from app.llm.response_parser import ResponseParser

p = ResponseParser()


def test_parses_plain_json():
    assert p.parse_json('{"a": 1}') == {"a": 1}


def test_parses_fenced_json():
    text = "Here you go:\n```json\n{\"a\": 2}\n```\nthanks"
    assert p.parse_json(text) == {"a": 2}


def test_parses_json_with_surrounding_prose():
    text = 'Sure! {"score": 4, "ok": true} done.'
    assert p.parse_json(text) == {"score": 4, "ok": True}


def test_parses_nested_braces_in_strings():
    text = '{"note": "use {curly} braces", "n": 1}'
    assert p.parse_json(text) == {"note": "use {curly} braces", "n": 1}


def test_raises_on_no_json():
    try:
        p.parse_json("no json here")
        assert False, "expected ValueError"
    except ValueError:
        pass


def test_build_envelope_lifts_fields():
    data, ev, conf, insf = p.build_envelope({
        "x": 1, "evidence": "quote", "confidence": 0.9,
    })
    assert ev == "quote"
    assert conf == 0.9
    assert insf is False


def test_build_envelope_status_insufficient():
    _, _, _, insf = p.build_envelope({"status": "insufficient_evidence"})
    assert insf is True
