package services

import (
	"context"

	"github.com/vijayvenkatj/envcrypt/database"
)

type UserService struct {
	q *database.Queries
}

func NewUserService(q *database.Queries) *UserService {
	return &UserService{q: q}
}

func (s *UserService) GetByEmail(ctx context.Context, email string) ([]database.User, error) {
	return s.q.GetUsers(ctx)
}
