#!/usr/bin/env bash
# Run the full automated test suite across all three layers.
# Usage: bash scripts/test-all.sh
# Exits non-zero if any suite fails.

set -uo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
FAIL=0

section() { echo ""; echo "==================== $1 ===================="; }

section "Backend (Go) — build + vet + test"
(
  cd "$ROOT/backend" || exit 1
  go build ./... && go vet ./... && go test ./...
) || FAIL=1

section "AI service (Python) — pytest (mock mode)"
(
  cd "$ROOT/ai-service" || exit 1
  if [ -x ".venv/Scripts/python.exe" ]; then
    PY=".venv/Scripts/python.exe"
  elif [ -x ".venv/bin/python" ]; then
    PY=".venv/bin/python"
  else
    PY="python"
  fi
  AI_MOCK=true "$PY" -m pytest -q
) || FAIL=1

section "Frontend (Vue) — build"
(
  cd "$ROOT/frontend" || exit 1
  npm run build
) || FAIL=1

echo ""
if [ "$FAIL" -eq 0 ]; then
  echo "✅ ALL TEST SUITES PASSED"
else
  echo "❌ ONE OR MORE SUITES FAILED"
fi
exit $FAIL
