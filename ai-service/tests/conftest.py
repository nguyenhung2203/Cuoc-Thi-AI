import os
import sys

# Serving configuration remains strict. Tests provide a syntactically valid
# non-production credential and inject a provider double at the API boundary.
os.environ.setdefault("GEMINI_API_KEY", "test-only-injected-provider-key")

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
if ROOT not in sys.path:
    sys.path.insert(0, ROOT)
