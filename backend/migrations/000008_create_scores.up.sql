CREATE TABLE interview_scores (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    interview_id UUID NOT NULL REFERENCES interviews(id),
    rubric_criterion_id UUID REFERENCES rubric_criteria(id),
    criterion_name VARCHAR(255) NOT NULL,
    score NUMERIC(4,2),
    max_score NUMERIC(4,2) NOT NULL,
    weight NUMERIC(5,2) NOT NULL,
    weighted_score NUMERIC(6,2),
    evidence TEXT,
    ai_comment TEXT,
    confidence NUMERIC(5,4),
    status VARCHAR(50) NOT NULL,
    scored_by VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX interview_scores_unique_idx ON interview_scores (interview_id, criterion_name);