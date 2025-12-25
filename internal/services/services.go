package services

import "github.com/vijayvenkatj/envcrypt/database"

type Services struct {
	Users    *UserService
	Projects *ProjectService
	Env      *EnvServices
}

func NewServices(queries *database.Queries) *Services {
	return &Services{
		Users:    NewUserService(queries),
		Projects: NewProjectService(queries),
		Env:      NewEnvServices(queries),
	}
}
