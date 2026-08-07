import json


def validate_json(data: str) -> bool:
    """Return True if `data` is a parseable JSON document."""
    if not isinstance(data, str) or not data.strip():
        return False
    try:
        json.loads(data)
        return True
    except (json.JSONDecodeError, ValueError):
        return False
