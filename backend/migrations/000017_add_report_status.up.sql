ALTER TABLE interviews ADD COLUMN report_status VARCHAR(50) DEFAULT 'pending';

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    link VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
