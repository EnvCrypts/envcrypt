
-- name: CreateCISession :one
INSERT INTO sessions (
    identity_type,
    service_role_id,
    project_id,
    env,
    github_repo
)
VALUES ('ci',$1, $2, $3, $4)
RETURNING id, created_at, expires_at;

-- name: CreateUserSession :one
INSERT INTO sessions (
    identity_type,
    user_id
)
VALUES ('user',$1)
RETURNING id, created_at, expires_at;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 AND expires_at > NOW();


-- name: RefreshToken :one
SELECT id, created_at, expires_at
FROM refresh_tokens
WHERE user_id = $1 AND expires_at > NOW();

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id)
VALUES ($1)
    RETURNING id, created_at, expires_at;

-- name: DeleteRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

-- name: DeleteUserAccessTokens :exec
DELETE FROM sessions
WHERE identity_type = 'user'
  AND user_id = $1;
