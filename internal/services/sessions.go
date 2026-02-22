package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

	session, err := s.q.CreateCISession(ctx, database.CreateCISessionParams{
		ProjectID:     uuid.NullUUID{UUID: projectDelegation.ProjectID, Valid: true},
		Env:           sql.NullString{String: projectDelegation.Env, Valid: true},
		ServiceRoleID: uuid.NullUUID{UUID: serviceRole.ID, Valid: true},
		GithubRepo:    sql.NullString{String: repoPrincipal, Valid: true},
	})

	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("failed to create session for %s/%s", projectDelegation.ProjectName, projectDelegation.Env))
	}

	return &session.ID, &projectDelegation.ProjectID, nil
}

func (s *SessionService) Refresh(ctx context.Context, userID uuid.UUID) (*uuid.UUID, *uuid.UUID, error) {

	var accessToken, refreshToken uuid.UUID

	refreshTokenDB, err := s.q.RefreshToken(ctx, userID)
	if err != nil {
		if dberrors.IsNoRows(err) {
			refreshTokenDB, err := s.q.CreateRefreshToken(ctx, userID)
			if err != nil {
				if dberrors.IsUniqueViolation(err) {
					return nil, nil, errors.New(fmt.Sprintf("refresh token for user %s already exists", userID))
				}
				return nil, nil, err
			}
			refreshToken = refreshTokenDB.ID
		} else {
			return nil, nil, err
		}
	}
	refreshToken = refreshTokenDB.ID

	accessTokenDB, err := s.q.CreateUserSession(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		if dberrors.IsUniqueViolation(err) {
			return nil, nil, errors.New(fmt.Sprintf("user session for %s already exists", userID))
		}
		return nil, nil, err
	}
	accessToken = accessTokenDB.ID

	return &accessToken, &refreshToken, nil
}

func (s *SessionService) GetProjectKeys(ctx context.Context, requestBody config.ServiceRollProjectKeyRequest) (*config.ServiceRollProjectKeyResponse, error) {

	session, err := s.q.GetSession(ctx, requestBody.SessionID)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, errors.New(fmt.Sprintf("failed to get session for %s", requestBody.SessionID))
		}
		return nil, err
	}

	var sessionProjectID, sessionServiceRoleID uuid.UUID
	var sessionEnv string

	if session.ProjectID.Valid && session.ServiceRoleID.Valid && session.Env.Valid {
		sessionProjectID = session.ProjectID.UUID
		sessionServiceRoleID = session.ServiceRoleID.UUID
		sessionEnv = session.Env.String
	} else {
		return nil, errors.New(fmt.Sprintf("failed to get session for %s", requestBody.SessionID))
	}

	if (sessionProjectID != requestBody.ProjectID) || (sessionEnv != requestBody.Env) {
		return nil, errors.New(fmt.Sprintf("project id and env are not the same"))
	}

	projectKeys, err := s.q.GetDelegatedKeys(ctx, database.GetDelegatedKeysParams{
		ProjectID:     sessionProjectID,
		Env:           sessionEnv,
		ServiceRoleID: sessionServiceRoleID,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, errors.New(fmt.Sprintf("failed to get project keys for %s", requestBody.SessionID))
		}
		return nil, err
	}

	return &config.ServiceRollProjectKeyResponse{
		ProjectId:          sessionProjectID,
		WrappedPRK:         projectKeys.WrappedPrk,
		WrapNonce:          projectKeys.WrapNonce,
		EphemeralPublicKey: projectKeys.WrapEphemeralPub,
	}, nil
}

func (s *SessionService) GetSession(ctx context.Context, sessionID uuid.UUID) error {
	_, err := s.q.GetSession(ctx, sessionID)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return errors.New(fmt.Sprintf("failed to get session for %s", sessionID))
		}
		return err
	}

	return nil
}
