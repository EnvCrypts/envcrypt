package config

import "github.com/google/uuid"

type Metadata struct {
	Type string `json:"type"`
}

type GetEnvRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	Email     string    `json:"user_email"`

	EnvName string `json:"env_name"`
	Version *int32 `json:"version"`
}

type GetEnvResponse struct {
	CipherText        []byte `json:"cipher_text"`
	Nonce             []byte `json:"nonce"`
	WrappedDEK        []byte `json:"wrapped_dek,omitempty"`
	DekNonce          []byte `json:"dek_nonce,omitempty"`
	EncryptionVersion int32  `json:"encryption_version"`
}

type GetEnvVersionsRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	Email     string    `json:"user_email"`

	EnvName string `json:"env_name"`
}

type EnvResponse struct {
	CipherText        []byte   `json:"cipher_text"`
	Nonce             []byte   `json:"nonce"`
	WrappedDEK        []byte   `json:"wrapped_dek,omitempty"`
	DekNonce          []byte   `json:"dek_nonce,omitempty"`
	EncryptionVersion int32    `json:"encryption_version"`
	Version           int32    `json:"version"`
	Metadata          Metadata `json:"metadata"`
}
type GetEnvVersionsResponse struct {
	EnvVersions []EnvResponse `json:"env_versions"`
}
type AddEnvRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	UserId    uuid.UUID `json:"user_id"`

	EnvName           string `json:"env_name"`
	CipherText        []byte `json:"cipher_text"`
	Nonce             []byte `json:"nonce"`
	WrappedDEK        []byte `json:"wrapped_dek"`
	DekNonce          []byte `json:"dek_nonce"`
	EncryptionVersion int32  `json:"encryption_version"`

	Metadata Metadata `json:"metadata"`
}

type AddEnvResponse struct {
	Message string `json:"message"`
}

type UpdateEnvRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	Email     string    `json:"user_email"`

	EnvName           string `json:"env_name"`
	CipherText        []byte `json:"cipher_text"`
	Nonce             []byte `json:"nonce"`
	WrappedDEK        []byte `json:"wrapped_dek"`
	DekNonce          []byte `json:"dek_nonce"`
	EncryptionVersion int32  `json:"encryption_version"`

	Metadata Metadata `json:"metadata"`
}

type UpdateEnvResponse struct {
	Message string `json:"message"`
}

type GetEnvForCIRequest struct {
	ProjectId uuid.UUID `json:"project_id"`
	EnvName   string    `json:"env_name"`
}
type GetEnvForCIResponse struct {
	CipherText        []byte `json:"cipher_text"`
	Nonce             []byte `json:"nonce"`
	WrappedDEK        []byte `json:"wrapped_dek"`
	DekNonce          []byte `json:"dek_nonce"`
	EncryptionVersion int32  `json:"encryption_version"`
}
