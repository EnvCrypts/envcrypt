package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
)

type SessionService struct {
	q *database.Queries
}

func NewSessionService(q *database.Queries) *SessionService {
	return &SessionService{q: q}
}

func (s *SessionService) Create(ctx context.Context, repoPrincipal string) (*uuid.UUID, *uuid.UUID, error) {

	serviceRole, err := s.q.GetServiceRoleByPrincipal(ctx, repoPrincipal)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, nil, errors.New(fmt.Sprintf("service role for %s not found", repoPrincipal))
		}
		return nil, nil, err
	}

	projectDelegation, err := s.q.GetDelegation(ctx, serviceRole.ID)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, nil, errors.New(fmt.Sprintf("service delegation for %s not found", repoPrincipal))
		}
		return nil, nil, err
	}

	session, err := s.q.CreateSession(ctx, database.CreateSessionParams{
		ProjectID:     projectDelegation.ProjectID,
		Env:           projectDelegation.Env,
		ServiceRoleID: serviceRole.ID,
		GithubRepo:    repoPrincipal,
	})
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New(fmt.Sprintf("failed to create session for %s/%s", projectDelegation.ProjectName, projectDelegation.Env))
	}

	return &session.ID, &projectDelegation.ProjectID, nil
}

func (s *SessionService) GetProjectKeys(ctx context.Context, requestBody config.ServiceRollProjectKeyRequest) (*config.ServiceRollProjectKeyResponse, error) {

	session, err := s.q.GetSession(ctx, requestBody.SessionID)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, errors.New(fmt.Sprintf("failed to get session for %s", requestBody.SessionID))
		}
		return nil, err
	}

	if (session.ProjectID != requestBody.ProjectID) || (session.Env != requestBody.Env) {
		return nil, errors.New(fmt.Sprintf("project id and env are not the same"))
	}

	projectKeys, err := s.q.GetDelegatedKeys(ctx, database.GetDelegatedKeysParams{
		ProjectID:     session.ProjectID,
		Env:           session.Env,
		ServiceRoleID: session.ServiceRoleID,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, errors.New(fmt.Sprintf("failed to get project keys for %s", requestBody.SessionID))
		}
		return nil, err
	}

	return &config.ServiceRollProjectKeyResponse{
		ProjectId:          session.ProjectID,
		WrappedPMK:         projectKeys.WrappedPmk,
		WrapNonce:          projectKeys.WrapNonce,
		EphemeralPublicKey: projectKeys.WrapEphemeralPub,
	}, nil
}
