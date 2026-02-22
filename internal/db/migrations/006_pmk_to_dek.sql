-- +goose Up
ALTER TABLE project_wrapped_keys RENAME COLUMN wrapped_pmk TO wrapped_prk;
ALTER TABLE service_delegations RENAME COLUMN wrapped_pmk TO wrapped_prk;

ALTER TABLE env_versions ADD COLUMN wrapped_dek BYTEA;
ALTER TABLE env_versions ADD COLUMN dek_nonce BYTEA;
ALTER TABLE env_versions ADD COLUMN encryption_version INTEGER NOT NULL DEFAULT 2;

-- +goose Down
ALTER TABLE env_versions DROP COLUMN encryption_version;
ALTER TABLE env_versions DROP COLUMN dek_nonce;
ALTER TABLE env_versions DROP COLUMN wrapped_dek;

ALTER TABLE service_delegations RENAME COLUMN wrapped_prk TO wrapped_pmk;
ALTER TABLE project_wrapped_keys RENAME COLUMN wrapped_prk TO wrapped_pmk;
