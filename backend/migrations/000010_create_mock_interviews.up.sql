CREATE TABLE mock_interviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    target_role VARCHAR(255) NOT NULL,
    target_level VARCHAR(100),
    cv_file_id UUID REFERENCES files(id),
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    final_score NUMERIC(6,2),
    feedback_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mock_interview_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    mock_interview_id UUID NOT NULL REFERENCES mock_interviews(id),
    sender_type VARCHAR(50) NOT NULL CHECK(sender_type IN ('ai', 'candidate', 'system')),
    content TEXT NOT NULL,
    question_type VARCHAR(100),
    score_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);