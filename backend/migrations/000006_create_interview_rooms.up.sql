CREATE TABLE interview_rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    interview_id UUID NOT NULL UNIQUE REFERENCES interviews(id),
    room_code VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'waiting',
    provider VARCHAR(100),
    connection_config JSONB,
    opened_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE interview_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    interview_id UUID NOT NULL REFERENCES interviews(id),
    user_id UUID REFERENCES users(id),
    participant_type VARCHAR(50) NOT NULL CHECK(participant_type IN ('recruiter', 'candidate', 'ai', 'guest')),
    display_name VARCHAR(255) NOT NULL,
    joined_at TIMESTAMPTZ,
    left_at TIMESTAMPTZ,
    last_seen_at TIMESTAMPTZ,
    connection_state VARCHAR(50) NOT NULL DEFAULT 'offline',
    media_status JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);\n