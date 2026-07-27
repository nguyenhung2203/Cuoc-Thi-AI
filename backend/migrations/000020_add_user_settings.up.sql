-- Per-user settings blob (preferences: theme, language, notifications, ai
-- workspace toggles, privacy). Applies to both recruiter and candidate.
ALTER TABLE users ADD COLUMN IF NOT EXISTS settings JSONB NOT NULL DEFAULT '{}';
