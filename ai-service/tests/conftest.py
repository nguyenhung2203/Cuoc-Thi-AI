import os
import sys

# Force mock mode for all tests so no real Gemini key is needed.
os.environ["AI_MOCK"] = "true"
os.environ.setdefault("GEMINI_API_KEY", "")

# Ensure the ai-service root is importable (app package).
ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
if ROOT not in sys.path:
    sys.path.insert(0, ROOT)
