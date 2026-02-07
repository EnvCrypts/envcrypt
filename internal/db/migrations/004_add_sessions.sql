-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE sessions (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

      identity_type TEXT NOT NULL CHECK (identity_type IN ('user', 'ci')),

      user_id UUID NULL REFERENCES users(id) ON DELETE CASCADE,

      service_role_id UUID NULL REFERENCES service_roles(id) ON DELETE CASCADE,
      project_id UUID NULL REFERENCES projects(id) ON DELETE CASCADE,
      env TEXT NULL,

      created_at TIMESTAMP NOT NULL DEFAULT NOW(),
      expires_at TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '10 minutes'),

      github_repo TEXT NULL   -- only for CI
);


CREATE INDEX idx_sessions_ci_lookup
    ON sessions(service_role_id, project_id, env);

CREATE INDEX idx_sessions_user_lookup
    ON sessions(user_id);

-- +goose Down
DROP TABLE sessions;
