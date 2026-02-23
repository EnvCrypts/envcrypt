-- +goose Up
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),
    request_id TEXT NOT NULL,
    actor_type TEXT NOT NULL, -- user | service | system
    actor_id TEXT NOT NULL,
    actor_email TEXT NOT NULL,
    action TEXT NOT NULL,
    project_id UUID,
    environment TEXT,
    target_id TEXT,
    ip_address INET,
    user_agent TEXT,
    status TEXT NOT NULL, -- success | failure
    error_message TEXT,
    metadata JSONB
);

CREATE INDEX idx_audit_logs_project_id ON audit_logs(project_id);
CREATE INDEX idx_audit_logs_actor_id ON audit_logs(actor_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp DESC);

-- +goose Down
DROP TABLE audit_logs;
