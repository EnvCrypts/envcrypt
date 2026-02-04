package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/envcrypt/database"
)

type SessionService struct {
	q *database.Queries
}

func NewSessionService(q *database.Queries) *SessionService {
	return &SessionService{q: q}
}

func (s *SessionService) Create(ctx context.Context, projectId, serviceRoleId uuid.UUID, envName, repoPrincipal string) (*uuid.UUID, error) {

	serviceRole, err := s.q.GetServiceRoleByPrincipal(ctx, repoPrincipal)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("service role for %s not found", repoPrincipal))
	}

	sid, err := s.q.HasAccess(ctx, database.HasAccessParams{
		ProjectID:     projectId,
		ServiceRoleID: serviceRole.ID,
		Env:           envName,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("access denied for %s", envName))
	}
	if sid != serviceRoleId {
		return nil, errors.New(fmt.Sprintf("access denied for %s", envName))
	}

	session, err := s.q.CreateSession(ctx, database.CreateSessionParams{
		ProjectID:     projectId,
		Env:           envName,
		ServiceRoleID: serviceRole.ID,
		GithubRepo:    repoPrincipal,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create session for %s", envName))
	}

	return &session.ID, nil
}
