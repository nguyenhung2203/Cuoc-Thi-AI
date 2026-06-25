CREATE TABLE interview_transcripts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    interview_id UUID NOT NULL REFERENCES interviews(id),
    participant_id UUID REFERENCES interview_participants(id),
    speaker_type VARCHAR(50) NOT NULL CHECK(speaker_type IN ('recruiter', 'candidate', 'ai', 'system')),
    speaker_name VARCHAR(255),
    content TEXT NOT NULL,
    language VARCHAR(20),
    start_time_ms INT,
    end_time_ms INT,
    confidence NUMERIC(5,4),
    source VARCHAR(50) NOT NULL CHECK(source IN ('audio', 'chat', 'manual', 'ai')),
    is_final BOOLEAN NOT NULL DEFAULT false,
    edited_content TEXT,
    edited_by UUID REFERENCES users(id),
    edited_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX interview_transcripts_interview_created_idx ON interview_transcripts (interview_id, created_at);
CREATE INDEX interview_transcripts_interview_speaker_idx ON interview_transcripts (interview_id, speaker_type);

CREATE TABLE ai_suggestions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    interview_id UUID NOT NULL REFERENCES interviews(id),
    suggestion_type VARCHAR(100) NOT NULL,
    title VARCHAR(255),
    content TEXT NOT NULL,
    reason TEXT,
    target_skill VARCHAR(255),
    priority VARCHAR(50),
    confidence NUMERIC(5,4),
    context_json JSONB,
    accepted_by UUID REFERENCES users(id),
    accepted_at TIMESTAMPTZ,
    dismissed_by UUID REFERENCES users(id),
    dismissed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);\n