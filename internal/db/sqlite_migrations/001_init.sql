-- +goose Up
PRAGMA foreign_keys = ON;

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,

    password_hash TEXT NOT NULL,
    password_salt BLOB NOT NULL,

    user_public_key BLOB NOT NULL,

    encrypted_user_private_key BLOB NOT NULL,
    private_key_nonce BLOB NOT NULL,
    private_key_salt BLOB NOT NULL,

    argon_params TEXT NOT NULL DEFAULT '{"time":3,"memory":65536,"parallelism":1}',

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    recovery_encrypted_private_key BLOB,
    recovery_nonce BLOB,
    recovery_kdf_salt BLOB
);

CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_by TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    prk_version INTEGER NOT NULL DEFAULT 1,
    UNIQUE (name, created_by)
);

CREATE TABLE project_members (
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('admin', 'member')),
    is_revoked INTEGER NOT NULL DEFAULT 0,
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, user_id)
);

CREATE TABLE project_wrapped_keys (
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wrapped_prk BLOB NOT NULL,
    wrap_nonce BLOB NOT NULL,
    wrap_ephemeral_pub BLOB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (project_id, user_id)
);

CREATE TABLE env_versions (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    env_name TEXT NOT NULL,
    version INTEGER NOT NULL,
    ciphertext BLOB NOT NULL,
    nonce BLOB NOT NULL,
    metadata TEXT NOT NULL DEFAULT '{"type":"env_created"}',
    created_by TEXT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    wrapped_dek BLOB,
    dek_nonce BLOB,
    encryption_version INTEGER NOT NULL DEFAULT 2,
    UNIQUE (project_id, env_name, version)
);

CREATE TABLE service_roles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    service_role_public_key BLOB NOT NULL,
    repo_principal TEXT NOT NULL,
    is_revoked INTEGER NOT NULL DEFAULT 0,
    created_by TEXT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_service_roles_repo
    ON service_roles(repo_principal);

CREATE TABLE service_delegations (
    service_role_id TEXT NOT NULL REFERENCES service_roles(id) ON DELETE CASCADE,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    env TEXT NOT NULL DEFAULT 'dev',
    wrapped_prk BLOB NOT NULL,
    wrap_nonce BLOB NOT NULL,
    wrap_ephemeral_pub BLOB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    delegated_by TEXT NOT NULL REFERENCES users(id),
    PRIMARY KEY (service_role_id, project_id, env)
);

CREATE INDEX idx_delegations_by_project
    ON service_delegations(project_id, env);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    identity_type TEXT NOT NULL CHECK (identity_type IN ('user', 'ci')),
    user_id TEXT NULL REFERENCES users(id) ON DELETE CASCADE,
    service_role_id TEXT NULL REFERENCES service_roles(id) ON DELETE CASCADE,
    project_id TEXT NULL REFERENCES projects(id) ON DELETE CASCADE,
    env TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, '+10 minutes')),
    github_repo TEXT NULL
);

CREATE INDEX idx_sessions_ci_lookup
    ON sessions(service_role_id, project_id, env);

CREATE INDEX idx_sessions_user_lookup
    ON sessions(user_id);

CREATE TABLE refresh_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, '+7 days')),
    UNIQUE (user_id, id)
);

CREATE INDEX idx_refresh_by_user
    ON refresh_tokens(user_id);

CREATE TABLE audit_logs (
    id TEXT PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    request_id TEXT NOT NULL,
    actor_type TEXT NOT NULL,
    actor_id TEXT NOT NULL,
    actor_email TEXT NOT NULL,
    action TEXT NOT NULL,
    project_id TEXT,
    environment TEXT,
    target_id TEXT,
    ip_address TEXT,
    user_agent TEXT,
    status TEXT NOT NULL,
    error_message TEXT,
    metadata TEXT
);

CREATE INDEX idx_audit_logs_project_id ON audit_logs(project_id);
CREATE INDEX idx_audit_logs_actor_id ON audit_logs(actor_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp DESC);

-- +goose Down
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS service_delegations;
DROP TABLE IF EXISTS service_roles;
DROP TABLE IF EXISTS env_versions;
DROP TABLE IF EXISTS project_wrapped_keys;
DROP TABLE IF EXISTS project_members;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;
