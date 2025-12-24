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
SELECT * FROm project_members WHERE project_id = $1 AND user_id = $2;

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
