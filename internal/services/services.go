package services

import (
	"database/sql"

	"github.com/vijayvenkatj/envcrypt/database"
)

type Services struct {
	Users          *UserService
	Projects       *ProjectService
	Env            *EnvServices
	ServiceRoles   *ServiceRoleServices
	SessionService *SessionService
	Audit          *AuditService
	Snapshot       *SnapshotService
}

func NewServices(queries *database.Queries, auditService *AuditService, db *sql.DB) *Services {
	users := NewUserService(queries)
	users.audit = auditService

	projects := NewProjectService(queries)
	projects.audit = auditService
	projects.db = db

	env := NewEnvService(queries)
	env.audit = auditService

	serviceRoles := NewServiceRoleService(queries)
	serviceRoles.audit = auditService

	sessionService := NewSessionService(queries)
	sessionService.audit = auditService

	snapshot := NewSnapshotService(queries, db)
	snapshot.SetAuditService(auditService)

	return &Services{
		Users:          users,
		Projects:       projects,
		Env:            env,
		ServiceRoles:   serviceRoles,
		SessionService: sessionService,
		Audit:          auditService,
		Snapshot:       snapshot,
	}
}

