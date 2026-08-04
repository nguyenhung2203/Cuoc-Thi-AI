-- Align existing AI configuration with the verified quota-friendly defaults.
-- Newer Gemini 3.5/3 Flash Live identifiers are intentionally not guessed here.
UPDATE system_settings
SET config_value = '"gemini-3.1-flash-lite"'::jsonb, updated_at = NOW()
WHERE config_key = 'default_ai_model'
  AND config_value = '"gemini-2.5-flash"'::jsonb;

UPDATE system_settings
SET config_value = '"gemini-2.5-flash-native-audio-latest"'::jsonb, updated_at = NOW()
WHERE config_key = 'default_ai_voice_model'
  AND config_value IN ('"gemini-2.0-flash-live-001"'::jsonb, '"gemini-2.5-flash"'::jsonb);

UPDATE ai_prompt_templates
SET model = 'gemini-3.1-flash-lite'
WHERE model IN ('gemini-1.5-flash', 'gemini-1.5-pro', 'gemini-2.5-pro', 'gemini-2.5-flash');
