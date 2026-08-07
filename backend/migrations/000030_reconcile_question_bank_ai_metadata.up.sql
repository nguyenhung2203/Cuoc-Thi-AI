-- Reconcile AI question metadata columns for databases where migration 000029
-- was recorded but the ALTER statements were not present in the live schema.
ALTER TABLE question_bank
    ADD COLUMN IF NOT EXISTS follow_up_prompts JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS timebox_minutes INT NOT NULL DEFAULT 5,
    ADD COLUMN IF NOT EXISTS evidence_required BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS generation_mode VARCHAR(20),
    ADD COLUMN IF NOT EXISTS prompt_version VARCHAR(100),
    ADD COLUMN IF NOT EXISTS ai_metadata JSONB NOT NULL DEFAULT '{}'::jsonb;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'question_bank_generation_mode_check'
    ) THEN
        ALTER TABLE question_bank ADD CONSTRAINT question_bank_generation_mode_check
            CHECK (generation_mode IS NULL OR generation_mode IN ('real', 'mock'));
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'question_bank_timebox_check'
    ) THEN
        ALTER TABLE question_bank ADD CONSTRAINT question_bank_timebox_check
            CHECK (timebox_minutes BETWEEN 1 AND 60);
    END IF;
END $$;
