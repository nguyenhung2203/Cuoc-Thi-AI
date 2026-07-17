-- Migrate deprecated Gemini model names to currently-available ones.
-- Google retired gemini-1.5-flash and gemini-1.5-pro on the v1beta API.
-- The free-tier key has NO quota for any Pro model (2.5-pro / 3.1-pro show 0/0),
-- so every template is moved to gemini-2.5-flash (best available in the usable tier).
-- Idempotent: safe to run repeatedly.
UPDATE ai_prompt_templates
SET model = 'gemini-2.5-flash'
WHERE model IN ('gemini-1.5-flash', 'gemini-1.5-pro', 'gemini-2.5-pro');
