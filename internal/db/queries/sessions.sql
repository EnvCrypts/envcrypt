
-- name: CreateSession :one
INSERT INTO sessions (
    service_role_id,
    project_id,
    env,
    github_repo
)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, expires_at;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1;