ALTER TABLE question_bank
    ADD COLUMN IF NOT EXISTS follow_up_prompts JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS timebox_minutes INT NOT NULL DEFAULT 5,
    ADD COLUMN IF NOT EXISTS evidence_required BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS generation_mode VARCHAR(20),
    ADD COLUMN IF NOT EXISTS prompt_version VARCHAR(100),
    ADD COLUMN IF NOT EXISTS ai_metadata JSONB NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE question_bank
    ADD CONSTRAINT question_bank_generation_mode_check
    CHECK (generation_mode IS NULL OR generation_mode IN ('real', 'mock'));

ALTER TABLE question_bank
    ADD CONSTRAINT question_bank_timebox_check
    CHECK (timebox_minutes BETWEEN 1 AND 60);
