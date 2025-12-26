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

-- name: GetProject :one

SELECT * from projects WHERE name = $1;


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
    created_by
)
VALUES (
           $1,
           $2,
           $3,
           $4,
           $5,
        $6
       )
RETURNING *;


-- name: GetEnv :one
SELECT * FROM env_versions WHERE project_id = $1 AND env_name = $2 AND version = $3;

