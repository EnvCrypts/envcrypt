-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email TEXT NOT NULL UNIQUE,

    password_hash TEXT NOT NULL,
    password_salt BYTEA NOT NULL,

    user_public_key BYTEA NOT NULL,

    encrypted_user_private_key BYTEA NOT NULL,
    private_key_nonce BYTEA NOT NULL,
    private_key_salt BYTEA NOT NULL,

    argon_params JSONB NOT NULL DEFAULT '{"time":3,"memory":65536,"parallelism":1}',

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;
