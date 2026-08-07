ALTER TABLE question_bank
    DROP CONSTRAINT IF EXISTS question_bank_generation_mode_check,
    DROP CONSTRAINT IF EXISTS question_bank_timebox_check,
    DROP COLUMN IF EXISTS follow_up_prompts,
    DROP COLUMN IF EXISTS timebox_minutes,
    DROP COLUMN IF EXISTS evidence_required,
    DROP COLUMN IF EXISTS generation_mode,
    DROP COLUMN IF EXISTS prompt_version,
    DROP COLUMN IF EXISTS ai_metadata;
