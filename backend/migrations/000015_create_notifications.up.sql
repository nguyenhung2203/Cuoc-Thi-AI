-- An earlier migration (000011) created a notifications table with a different
-- schema (content/read_at). This migration is the authoritative definition the
-- app repository relies on (message/is_read/link/company_id), so reconcile by
-- recreating it cleanly.
DROP TABLE IF EXISTS notifications CASCADE;
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID REFERENCES companies(id),
    user_id UUID REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    link VARCHAR(255),
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX notifications_user_id_idx ON notifications (user_id);
CREATE INDEX notifications_company_id_idx ON notifications (company_id);
