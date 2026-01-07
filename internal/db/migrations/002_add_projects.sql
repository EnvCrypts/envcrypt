-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE projects (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

      name TEXT NOT NULL,
      created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

      created_at TIMESTAMP NOT NULL DEFAULT NOW(),

      UNIQUE (name,created_by)
);


CREATE TABLE project_members (
        project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

        role TEXT NOT NULL CHECK (role IN ('admin', 'member')),

        added_at TIMESTAMP NOT NULL DEFAULT NOW(),

        PRIMARY KEY (project_id, user_id)
);


CREATE TABLE project_wrapped_keys (
          project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
          user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

          wrapped_pmk BYTEA NOT NULL,
          wrap_nonce BYTEA NOT NULL,
          wrap_ephemeral_pub BYTEA NOT NULL,

          created_at TIMESTAMP NOT NULL DEFAULT NOW(),

          PRIMARY KEY (project_id, user_id)
);


CREATE TABLE env_versions (
          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

          project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,

          env_name TEXT NOT NULL,      -- e.g. production, staging
          version INTEGER NOT NULL,

          ciphertext BYTEA NOT NULL,
          nonce BYTEA NOT NULL,

          metadata JSONB NOT NULL DEFAULT '{"type":"env_created"}',

          created_by UUID NOT NULL REFERENCES users(id),
          created_at TIMESTAMP NOT NULL DEFAULT NOW(),

          UNIQUE (project_id, env_name, version)
);

-- +goose Down
DROP TABLE IF EXISTS env_versions;
DROP TABLE IF EXISTS project_wrapped_keys;
DROP TABLE IF EXISTS project_members;
DROP TABLE IF EXISTS projects;
