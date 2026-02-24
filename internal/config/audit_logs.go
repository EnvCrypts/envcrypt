package config

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)


// Action types
const (
	ActionEnvPull             = "env.pull"
	ActionEnvPush             = "env.push"
	ActionPRKRotate           = "prk.rotate"
	ActionLogin               = "login"
	ActionLogout              = "logout"
	ActionRegister            = "register"
	ActionMembershipChange    = "membership.change"
	ActionProjectCreate       = "project.create"
	ActionProjectDelete       = "project.delete"
	ActionServiceRoleCreate   = "service_role.create"
	ActionServiceRoleDelete   = "service_role.delete"
	ActionServiceRoleDelegate = "service_role.delegate"
)

// Actor types
const (
	ActorTypeUser    = "user"
	ActorTypeService = "service"
	ActorTypeSystem  = "system"
)

// Status types
const (
	StatusSuccess = "success"
	StatusFailure = "failure"
)


type AuditLog struct {
	ID           uuid.UUID       `json:"id"`
	Timestamp    time.Time       `json:"timestamp"`
	RequestID    string          `json:"request_id"`
	ActorType    string          `json:"actor_type"`
	ActorID      string          `json:"actor_id"`
	ActorEmail   string          `json:"actor_email"`
	Action       string          `json:"action"`
	ProjectID    *uuid.UUID      `json:"project_id"`
	Environment  *string         `json:"environment"`
	TargetID     *string         `json:"target_id"`
	IPAddress    *string         `json:"ip_address"`
	UserAgent    *string         `json:"user_agent"`
	Status       string          `json:"status"`
	ErrorMessage *string         `json:"error_message"`
	Metadata     json.RawMessage `json:"metadata"`
}

type ProjectAuditRequest struct {
	ProjectID  uuid.UUID  `json:"project_id"`
	Limit      int32      `json:"limit"`
	Offset     int32      `json:"offset"`
	ActorEmail *string    `json:"actor_email"`
	Action     *string    `json:"action"`
	Status     *string    `json:"status"`
	From       *time.Time `json:"from"`
	To         *time.Time `json:"to"`
}

type ProjectAuditResponse struct {
	Logs       []AuditLog `json:"logs"`
	Pagination struct {
		Limit  int32 `json:"limit"`
		Offset int32 `json:"offset"`
		Total  int64 `json:"total"`
	} `json:"pagination"`
}
