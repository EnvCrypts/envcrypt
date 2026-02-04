-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE service_roles (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       name TEXT NOT NULL,

       service_role_public_key BYTEA NOT NULL,
       repo_principal TEXT NOT NULL,
       is_revoked BOOL NOT NULL DEFAULT false,

       created_by UUID NOT NULL REFERENCES users(id),
       created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_service_roles_repo
    ON service_roles(repo_principal);

CREATE TABLE service_delegations (
         service_role_id UUID NOT NULL
             REFERENCES service_roles(id) ON DELETE CASCADE,

         project_id UUID NOT NULL
             REFERENCES projects(id) ON DELETE CASCADE,

         env TEXT NOT NULL DEFAULT 'dev',

         wrapped_pmk BYTEA NOT NULL,
         wrap_nonce BYTEA NOT NULL,
         wrap_ephemeral_pub BYTEA NOT NULL,

         created_at TIMESTAMP NOT NULL DEFAULT NOW(),
         delegated_by UUID NOT NULL REFERENCES users(id),

         PRIMARY KEY (service_role_id, project_id, env)
);


CREATE INDEX idx_delegations_by_project
    ON service_delegations(project_id, env);

-- +goose Down
DROP TABLE service_delegations;
DROP TABLE service_roles;
