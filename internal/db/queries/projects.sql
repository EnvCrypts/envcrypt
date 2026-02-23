-- name: GetUserProjects :many
SELECT * FROM projects WHERE created_by = $1 ORDER BY created_at desc;

-- name: CreateProject :one
INSERT INTO projects (
    name,
    created_by
)
VALUES (
           $1,
           $2
       )
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;

-- name: GetProject :one
SELECT * from projects WHERE name = $1 AND created_by = $2;

-- name: GetProjectById :one
SELECT * from projects WHERE id = $1;

-- name: ListProjectsWithRole :many
SELECT
    p.id,
    p.name,
    p.created_by,
    p.created_at,
    pm.role,
    pm.is_revoked
FROM projects p JOIN project_members pm ON pm.project_id = p.id
WHERE pm.user_id = $1
ORDER BY p.created_at DESC;

-- name: GetMemberProject :one
SELECT p.id
FROM projects p
         JOIN project_members pm ON pm.project_id = p.id
WHERE p.name = $1
  AND pm.user_id = $2;


-- name: AddUserToProject :one
INSERT INTO project_members (
    project_id,
    user_id,
    role
)
VALUES (
           $1,
           $2,
            $3
       )
RETURNING *;

-- name: GetUserProjectRole :one
SELECT * FROM project_members WHERE project_id = $1 AND user_id = $2 and is_revoked = $3;

-- name: SetUserAccess :exec
UPDATE project_members
SET is_revoked = $3
WHERE user_id = $1 AND project_id = $2;


-- name: AddWrappedPRK :one
INSERT INTO project_wrapped_keys (
    project_id,
    user_id,

    wrapped_prk,
    wrap_nonce,
    wrap_ephemeral_pub
)
VALUES (
           $1,
           $2,
           $3,
        $4,
        $5
       )
RETURNING *;

-- name: GetProjectWrappedKey :one
SELECT * FROM project_wrapped_keys WHERE project_id = $1 AND user_id = $2;


-- name: AddEnv :one
INSERT INTO env_versions (
    project_id,
    env_name,
    version,
    nonce,
    ciphertext,
    wrapped_dek,
    dek_nonce,
    encryption_version,
    created_by,
    metadata
)
SELECT
    $1,
    $2,
    COALESCE(MAX(version), 0) + 1,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
FROM env_versions
WHERE project_id = $1 AND env_name = $2
RETURNING *;




-- name: GetEnv :one
SELECT * FROM env_versions WHERE project_id = $1 AND env_name = $2 AND version = $3;

-- name: GetLatestEnv :one
SELECT * FROM env_versions WHERE project_id = $1 AND env_name = $2 ORDER BY version DESC LIMIT 1;

-- name: GetEnvVersions :many
SELECT * FROM env_versions WHERE project_id = $1 AND env_name = $2 ORDER BY version DESC;

-- name: GetRotationData :many
SELECT
    pwk.user_id,
    pwk.wrapped_prk,
    pwk.wrap_nonce,
    pwk.wrap_ephemeral_pub,
    u.user_public_key
FROM project_wrapped_keys pwk
JOIN users u ON u.id = pwk.user_id
WHERE pwk.project_id = $1;

-- name: GetProjectWrappedDEKs :many
SELECT id, wrapped_dek, dek_nonce
FROM env_versions
WHERE project_id = $1 AND wrapped_dek IS NOT NULL;

-- name: UpdateWrappedPRK :exec
UPDATE project_wrapped_keys
SET wrapped_prk = $3, wrap_nonce = $4, wrap_ephemeral_pub = $5
WHERE project_id = $1 AND user_id = $2;

-- name: UpdateEnvVersionDEK :exec
UPDATE env_versions
SET wrapped_dek = $2, dek_nonce = $3
WHERE id = $1;

-- name: IncrementPRKVersion :one
UPDATE projects
SET prk_version = prk_version + 1
WHERE id = $1 AND prk_version = $2
RETURNING prk_version;
