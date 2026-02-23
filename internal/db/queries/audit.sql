-- name: CreateAuditLog :exec
INSERT INTO audit_logs (
    id,
    timestamp,
    request_id,
    actor_type,
    actor_id,
    actor_email,
    action,
    project_id,
    environment,
    target_id,
    ip_address,
    user_agent,
    status,
    error_message,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
);
