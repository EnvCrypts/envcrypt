-- +goose Up
ALTER TABLE users ADD COLUMN recovery_encrypted_private_key BYTEA;
ALTER TABLE users ADD COLUMN recovery_nonce BYTEA;
ALTER TABLE users ADD COLUMN recovery_kdf_salt BYTEA;

-- +goose Down
ALTER TABLE users DROP COLUMN recovery_encrypted_private_key;
ALTER TABLE users DROP COLUMN recovery_nonce;
ALTER TABLE users DROP COLUMN recovery_kdf_salt;