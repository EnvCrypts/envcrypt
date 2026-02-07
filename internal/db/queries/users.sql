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
    argon_params
)
VALUES (
           $1, $2, $3, $4, $5, $6, $7, $8
       )
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;