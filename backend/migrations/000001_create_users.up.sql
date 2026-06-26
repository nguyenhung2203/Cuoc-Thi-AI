CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role AS ENUM ('admin', 'recruiter', 'candidate');
CREATE TYPE user_status AS ENUM ('pending', 'active', 'inactive', 'blocked');

CREATE TABLE users (
    id                UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    email             VARCHAR(255) NOT NULL,
    password_hash     TEXT,
    full_name         VARCHAR(255) NOT NULL,
    avatar_url        TEXT,
    role              user_role    NOT NULL DEFAULT 'recruiter',
    status            user_status  NOT NULL DEFAULT 'pending',
    last_login_at     TIMESTAMPTZ,
    email_verified_at TIMESTAMPTZ,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE UNIQUE INDEX users_email_unique ON users (email) WHERE deleted_at IS NULL;
CREATE INDEX users_role_idx            ON users (role);
CREATE INDEX users_status_idx          ON users (status);
