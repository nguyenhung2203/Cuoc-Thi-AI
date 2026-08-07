import json
import re
from typing import Any, Optional, Tuple

# Matches a ```json ... ``` or ``` ... ``` fenced block.
_FENCE_RE = re.compile(r"```(?:json)?\s*(.*?)\s*```", re.DOTALL | re.IGNORECASE)


def _strip_fences(text: str) -> str:
    m = _FENCE_RE.search(text)
    if m:
        return m.group(1).strip()
    return text.strip()


def _extract_balanced(text: str) -> Optional[str]:
    """Return the first balanced {...} or [...] block found in text.

    LLMs sometimes wrap JSON in prose; this pulls out the object/array by
    tracking brace/bracket depth while ignoring braces inside strings.
    """
    start = None
    opener = closer = ""
    for i, ch in enumerate(text):
        if ch in "{[":
            start = i
            opener = ch
            closer = "}" if ch == "{" else "]"
            break
    if start is None:
        return None

    depth = 0
    in_str = False
    esc = False
    for i in range(start, len(text)):
        ch = text[i]
        if in_str:
            if esc:
                esc = False
            elif ch == "\\":
                esc = True
            elif ch == '"':
                in_str = False
            continue
        if ch == '"':
            in_str = True
        elif ch == opener:
            depth += 1
        elif ch == closer:
            depth -= 1
            if depth == 0:
                return text[start:i + 1]
    return None


class ResponseParser:
    """Turns raw LLM text into a JSON object and derives the standard
    envelope fields (evidence, confidence, insufficient_data)."""

    @staticmethod
    def parse_json(text: str) -> Any:
        """Best-effort JSON extraction. Raises ValueError if nothing parses."""
        if text is None:
            raise ValueError("empty LLM response")

        candidate = _strip_fences(text)
        # Try direct parse first.
        try:
            return json.loads(candidate)
        except json.JSONDecodeError:
            pass

        # Fall back to extracting a balanced block.
        block = _extract_balanced(candidate) or _extract_balanced(text)
        if block is not None:
            return json.loads(block)

        raise ValueError("no valid JSON object found in LLM response")

    @staticmethod
    def build_envelope(data: Any) -> Tuple[Any, str, float, bool]:
        """Return (data, evidence, confidence, insufficient_data).

        The Go side reads these top-level fields. When the model embeds
        `evidence`/`confidence`/`insufficient_data` inside the object we lift
        them to the envelope; otherwise sensible defaults are used.
        """
        evidence = ""
        confidence = 0.0
        insufficient = False

        if isinstance(data, dict):
            ev = data.get("evidence")
            if isinstance(ev, str):
                evidence = ev
            conf = data.get("confidence")
            if isinstance(conf, (int, float)):
                confidence = float(conf)
            insf = data.get("insufficient_data") or data.get("insufficient_evidence")
            if isinstance(insf, bool):
                insufficient = insf
            # Some scoring prompts signal insufficiency via a status field.
            status = data.get("status")
            if isinstance(status, str) and status == "insufficient_evidence":
                insufficient = True

        return data, evidence, confidence, insufficient
