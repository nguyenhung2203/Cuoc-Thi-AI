-- Keep this data migration reversible for environments that still use the
-- previous verified text model. New installations use 3.1 Flash Lite defaults.
UPDATE system_settings
SET config_value = '"gemini-2.5-flash"'::jsonb, updated_at = NOW()
WHERE config_key = 'default_ai_model'
  AND config_value = '"gemini-3.1-flash-lite"'::jsonb;

UPDATE system_settings
SET config_value = '"gemini-2.0-flash-live-001"'::jsonb, updated_at = NOW()
WHERE config_key = 'default_ai_voice_model'
  AND config_value = '"gemini-2.5-flash-native-audio-latest"'::jsonb;

UPDATE ai_prompt_templates
SET model = 'gemini-2.5-flash'
WHERE model = 'gemini-3.1-flash-lite';
