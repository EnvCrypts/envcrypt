

-- name: GetServiceRoleById :one
SELECT * FROM service_roles WHERE id = $1;

-- name: GetServiceRoleByPrincipal :one
SELECT * FROM service_roles WHERE repo_principal = $1;

-- name: GetServiceRolesByAdmin :many
SELECT * FROM service_roles WHERE created_by = $1;

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


-- name: DelegateAccess :one
INSERT INTO service_delegations (
    service_role_id,
    project_id,
    env,
    wrapped_pmk,
    wrap_nonce,
    wrap_ephemeral_pub,
    delegated_by
)
VALUES (
           $1,  -- service_role_id
           $2,  -- project_id
           $3,  -- env
           $4,  -- wrapped_pmk
           $5,  -- wrap_nonce
           $6,  -- wrap_ephemeral_pub (admin_eph_pub)
           $7   -- delegated_by (admin user id)
       )
RETURNING
service_role_id,
project_id,
env,
created_at;


-- name: HasAccess :one
SELECT service_role_id FROM service_delegations
WHERE service_role_id = $1
  AND project_id = $2
  AND env = $3;

-- name: GetDelegatedKeys :one
SELECT * FROM service_delegations WHERE service_role_id = $1 AND project_id = $2 AND env = $3;

-- name: GetDelegation :one
SELECT
    d.service_role_id,
    d.project_id,
    d.env,
    d.created_at,
    p.name        AS project_name
FROM service_delegations d
         JOIN projects p
              ON d.project_id = p.id
WHERE d.service_role_id = $1
    LIMIT 1;