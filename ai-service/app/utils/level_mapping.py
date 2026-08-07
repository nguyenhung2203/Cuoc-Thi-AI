from typing import Literal

LEVEL_ALIASES = {"mid": "middle", "medior": "middle"}
LEVELS = ("intern", "fresher", "junior", "middle", "senior", "lead", "unknown")
MODES = ("real", "mock")


def normalize_level(value: str | None) -> str:
    value = (value or "unknown").strip().lower()
    value = LEVEL_ALIASES.get(value, value)
    return value if value in LEVELS else "unknown"


def normalize_mode(value: str | None) -> str:
    value = (value or "real").strip().lower()
    if value not in MODES:
        raise ValueError("mode must be real or mock")
    return value


LEVEL_DIFFICULTIES = {
    "intern": {"basic", "intermediate"},
    "fresher": {"basic", "intermediate"},
    "junior": {"basic", "intermediate"},
    "middle": {"intermediate", "advanced"},
    "senior": {"intermediate", "advanced"},
    "lead": {"intermediate", "advanced"},
    "unknown": {"basic", "intermediate"},
}

def level_guidance(level: str) -> str:
    return {
        "intern": "fundamentals, learning ability, coursework and basic behavior",
        "fresher": "fundamentals, project learning, implementation basics and teamwork",
        "junior": "implementation, debugging, testing and teamwork with guidance",
        "middle": "design, trade-offs, delivery, edge cases and operational ownership",
        "senior": "architecture, scale, reliability, risk, mentoring and ownership",
        "lead": "strategy, architecture, leadership, stakeholder alignment and outcomes",
        "unknown": "foundations and role-relevant questions; avoid assuming seniority",
    }[level]
