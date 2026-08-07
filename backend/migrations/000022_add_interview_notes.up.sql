-- Internal recruiter notes for an interview (visible only to recruiter/admin).
-- Kept as free-form text; not shown to candidates.
ALTER TABLE interviews ADD COLUMN IF NOT EXISTS recruiter_notes TEXT NOT NULL DEFAULT '';
