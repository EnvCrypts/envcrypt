-- name: GetUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    password_salt,
    user_public_key,
    encrypted_user_private_key,
    private_key_nonce,
    private_key_salt,
    recovery_encrypted_private_key,
    recovery_nonce,
    recovery_kdf_salt,
    argon_params
)
VALUES (
           $1, $2, $3, $4, $5, $6, $7, $8, $9,$10,$11
       )
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserCredentials :one
UPDATE users
SET password_hash = $2,
    password_salt = $3,
    argon_params = $4,
    encrypted_user_private_key = $5,
    private_key_nonce = $6,
    private_key_salt = $7
WHERE email = $1
RETURNING *;