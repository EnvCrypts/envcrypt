package config

import "github.com/google/uuid"

type GetEnvRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	Email     string    `json:"user_email"`

	EnvName string `json:"env_name"`
	Version *int32 `json:"version"`
}

type GetEnvResponse struct {
	CipherText []byte `json:"cipher_text"`
	Nonce      []byte `json:"nonce"`
}

type GetEnvVersionsRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	Email     string    `json:"user_email"`

	EnvName string `json:"env_name"`
}

type EnvResponse struct {
	CipherText []byte `json:"cipher_text"`
	Nonce      []byte `json:"nonce"`
	Version    int32  `json:"version"`
}
type GetEnvVersionsResponse struct {
	EnvVersions []EnvResponse `json:"env_versions"`
}
type AddEnvRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	Email     string    `json:"user_email"`

	EnvName    string `json:"env_name"`
	CipherText []byte `json:"cipher_text"`
	Nonce      []byte `json:"nonce"`
}

type AddEnvResponse struct {
	Message string `json:"message"`
}

type UpdateEnvRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	Email     string    `json:"user_email"`

	EnvName    string `json:"env_name"`
	CipherText []byte `json:"cipher_text"`
	Nonce      []byte `json:"nonce"`
}

type UpdateEnvResponse struct {
	Message string `json:"message"`
}
