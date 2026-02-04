package config

import (
	"time"

	"github.com/google/uuid"
)

type ServiceRole struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	ServiceRolePublicKey []byte `json:"service_role_public_key"`
	RepoPrincipal        string `json:"repo_principal"`

	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// ServiceRoleCreateRequest POST /service_role/create
type ServiceRoleCreateRequest struct {
	ServiceRoleName string `json:"service_role_name"`

	ServiceRolePublicKey []byte `json:"service_role_public_key"`

	RepoPrincipal string    `json:"repo_principal"`
	CreatedBy     uuid.UUID `json:"created_by"`
}
type ServiceRoleCreateResponse struct {
	Message     string      `json:"message"`
	ServiceRole ServiceRole `json:"service_role"`
}

// ServiceRoleGetRequest POST /service_role/get
type ServiceRoleGetRequest struct {
	RepoPrincipal string `json:"repo-principal"`
}
type ServiceRoleGetResponse struct {
	ServiceRole ServiceRole `json:"service_role"`
	Message     string      `json:"message"`
}

// ServiceRoleDeleteRequest POST /service_role/delete
type ServiceRoleDeleteRequest struct {
	ServiceRoleId uuid.UUID `json:"service_role_name"`
	CreatedBy     uuid.UUID `json:"created-by"`
}
type ServiceRoleDeleteResponse struct {
	Message string `json:"message"`
}
