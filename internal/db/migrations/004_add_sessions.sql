-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE sessions (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

      service_role_id UUID NOT NULL REFERENCES service_roles(id) ON DELETE CASCADE,

      project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
      env TEXT NOT NULL DEFAULT 'dev',

      created_at TIMESTAMP NOT NULL DEFAULT NOW(),
      expires_at TIMESTAMP NOT NULL,

      github_repo TEXT NOT NULL,   -- e.g. "github:org/repo:main"

      UNIQUE (service_role_id, project_id, env)
);

CREATE INDEX idx_sessions_lookup ON sessions(service_role_id, project_id, env);

-- +goose Down
DROP TABLE sessions;
