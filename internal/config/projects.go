package config

import "github.com/google/uuid"

type ProjectCreateRequest struct {
	Name               string    `json:"name"`
	UserId             uuid.UUID `json:"user_id"`
	WrappedPMK         []byte    `json:"wrapped_pmk"`
	WrapNonce          []byte    `json:"wrap_nonce"`
	EphemeralPublicKey []byte    `json:"ephemeral_public_key"`
}

type ProjectCreateResponse struct {
	Message string `json:"message"`
}

type AddUserToProjectRequest struct {
	ProjectId          uuid.UUID `json:"project_id"`
	AdminId            uuid.UUID `json:"admin_id"`
	UserId             uuid.UUID `json:"user_id"`
	WrappedPMK         []byte    `json:"wrapped_pmk"`
	WrapNonce          []byte    `json:"wrap_nonce"`
	EphemeralPublicKey []byte    `json:"ephemeral_public_key"`
}
type AddUserToProjectResponse struct {
	Message string `json:"message"`
}
