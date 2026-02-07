-- +goose Up
CREATE TABLE refresh_tokens (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        expires_at TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '7 days'),

        UNIQUE (user_id, id)
);


CREATE INDEX idx_refresh_by_user
    ON refresh_tokens(user_id);

-- +goose Down
DROP TABLE refresh_tokens;
