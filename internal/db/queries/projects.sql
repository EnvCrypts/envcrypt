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

-- name: ListProjectsWithRole :many
SELECT
    p.id,
    p.name,
    p.created_by,
    p.created_at,
    pm.role
FROM projects p JOIN project_members pm ON pm.project_id = p.id
WHERE pm.user_id = $1
ORDER BY p.created_at DESC;

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
SELECT * FROM project_members WHERE project_id = $1 AND user_id = $2;

-- name: AddWrappedPMK :one

INSERT INTO project_wrapped_keys (
    project_id,
    user_id,

    wrapped_pmk,
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
    $6
FROM env_versions
WHERE project_id = $1 AND env_name = $2
RETURNING *;




-- name: GetEnv :one
SELECT * FROM env_versions WHERE project_id = $1 AND env_name = $2 AND version = $3;

-- name: GetLatestEnv :one
SELECT * FROM env_versions WHERE project_id = $1 AND env_name = $2 ORDER BY version DESC LIMIT 1;

-- name: GetEnvVersions :many
SELECT * FROM env_versions WHERE project_id = $1 AND env_name = $2 ORDER BY version DESC;
