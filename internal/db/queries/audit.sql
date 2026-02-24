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

-- name: GetProjectAuditLogsPaginated :many
SELECT *
FROM audit_logs
WHERE project_id = $1
  AND (sqlc.narg('actor_email')::text IS NULL OR actor_email = sqlc.narg('actor_email'))
  AND (sqlc.narg('action')::text IS NULL OR action = sqlc.narg('action'))
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('from_time')::timestamptz IS NULL OR timestamp >= sqlc.narg('from_time'))
  AND (sqlc.narg('to_time')::timestamptz IS NULL OR timestamp <= sqlc.narg('to_time'))
ORDER BY timestamp DESC
LIMIT sqlc.arg('limit_val') OFFSET sqlc.arg('offset_val');

-- name: CountProjectAuditLogs :one
SELECT COUNT(*)
FROM audit_logs
WHERE project_id = $1
  AND (sqlc.narg('actor_email')::text IS NULL OR actor_email = sqlc.narg('actor_email'))
  AND (sqlc.narg('action')::text IS NULL OR action = sqlc.narg('action'))
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('from_time')::timestamptz IS NULL OR timestamp >= sqlc.narg('from_time'))
  AND (sqlc.narg('to_time')::timestamptz IS NULL OR timestamp <= sqlc.narg('to_time'));
