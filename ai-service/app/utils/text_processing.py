import re

_WS_RE = re.compile(r"[ \t\f\v]+")
_MULTI_NL_RE = re.compile(r"\n{3,}")


def clean_text(text: str) -> str:
    """Normalize whitespace in free text without destroying paragraph breaks.

    Used to tidy CV/JD text before it is embedded in a prompt.
    """
    if not text:
        return ""
    text = text.replace("\r\n", "\n").replace("\r", "\n")
    # Collapse runs of spaces/tabs, trim each line, collapse 3+ blank lines.
    lines = [_WS_RE.sub(" ", ln).strip() for ln in text.split("\n")]
    joined = "\n".join(lines)
    return _MULTI_NL_RE.sub("\n\n", joined).strip()


def truncate(text: str, max_chars: int) -> str:
    """Hard cap text length to bound prompt/token size."""
    if text and len(text) > max_chars:
        return text[:max_chars]
    return text
