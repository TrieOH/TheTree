-- +goose Up
-- Created at 2025-10-11T21:24:52-03:00

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_id
    ON users(id);

-- +goose Down
DROP INDEX IF EXISTS idx_users_id;
DROP EXTENSION IF EXISTS "pgcrypto";
DROP TABLE IF EXISTS users;
     