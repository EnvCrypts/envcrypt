package services

import "github.com/vijayvenkatj/envcrypt/database"

type Services struct {
	Users *UserService
}

func NewServices(queries *database.Queries) *Services {
	return &Services{
		Users: NewUserService(queries),
	}
}
