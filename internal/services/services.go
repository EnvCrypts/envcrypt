package services

import "github.com/vijayvenkatj/envcrypt/database"

type Services struct {
	Users          *UserService
	Projects       *ProjectService
	Env            *EnvServices
	ServiceRoles   *ServiceRoleServices
	SessionService *SessionService
}

func NewServices(queries *database.Queries) *Services {
	return &Services{
		Users:          NewUserService(queries),
		Projects:       NewProjectService(queries),
		Env:            NewEnvService(queries),
		ServiceRoles:   NewServiceRoleService(queries),
		SessionService: NewSessionService(queries),
	}
}
