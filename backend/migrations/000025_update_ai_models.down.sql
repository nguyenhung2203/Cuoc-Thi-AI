-- No safe rollback: the previous model names (gemini-1.5-*) are retired by Google
-- and would break generateContent. Intentionally a no-op.
SELECT 1;
