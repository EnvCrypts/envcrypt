-- +goose Up
ALTER TABLE sessions DROP CONSTRAINT sessions_service_role_id_project_id_env_key;

-- +goose Down
ALTER TABLE sessions ADD CONSTRAINT sessions_service_role_id_project_id_env_key UNIQUE (service_role_id, project_id, env);
