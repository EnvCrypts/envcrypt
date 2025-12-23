package config

import "github.com/google/uuid"

type ProjectCreateRequest struct {
	Name               string    `json:"name"`
	UserId             uuid.UUID `json:"user_id"`
	WrappedPMK         []byte    `json:"wrapped_pmk"`
	WrapNonce          []byte    `json:"wrap_nonce"`
	EphemeralPublicKey []byte    `json:"ephemeral_public_key"`
}
