package services

import (
	"github.com/vijayvenkatj/envcrypt/database"
)

type Services struct {
	Users          *UserService
	Projects       *ProjectService
	Env            *EnvServices
	ServiceRoles   *ServiceRoleServices
	SessionService *SessionService
	Audit          *AuditService
}

func NewServices(queries *database.Queries, auditService *AuditService) *Services {
	users := NewUserService(queries)
	users.audit = auditService

	projects := NewProjectService(queries)
	projects.audit = auditService

	env := NewEnvService(queries)
	env.audit = auditService

	serviceRoles := NewServiceRoleService(queries)
	serviceRoles.audit = auditService

	sessionService := NewSessionService(queries)
	sessionService.audit = auditService

	return &Services{
		Users:          users,
		Projects:       projects,
		Env:            env,
		ServiceRoles:   serviceRoles,
		SessionService: sessionService,
		Audit:          auditService,
	}
}
