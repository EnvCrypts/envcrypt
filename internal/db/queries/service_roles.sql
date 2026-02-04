

-- name: GetServiceRoleById :one
SELECT * FROM service_roles WHERE id = $1;

-- name: GetServiceRoleByPrincipal :one
SELECT * FROM service_roles WHERE repo_principal = $1;

-- name: CreateServiceRole :one
INSERT INTO service_roles (
    name,
    service_role_public_key,
    repo_principal,
    created_by
)
VALUES (
           $1,   -- name
           $2,   -- service_role_public_key (BYTEA)
           $3,   -- repo_principal
           $4    -- created_by (UUID)
       )
RETURNING *;

-- name: DeleteServiceRole :one
DELETE FROM service_roles
WHERE id = $1
    RETURNING id;
