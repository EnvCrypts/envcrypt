package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
)

type SessionService struct {
	q     *database.Queries
	audit *AuditService
}

func NewSessionService(q *database.Queries) *SessionService {
	return &SessionService{q: q}
}

func (s *SessionService) Create(ctx context.Context, repoPrincipal string) (*uuid.UUID, *uuid.UUID, error) {

	serviceRole, err := s.q.GetServiceRoleByPrincipal(ctx, repoPrincipal)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionLogin, ActorType: config.ActorTypeService, ActorID: "unknown", ActorEmail: repoPrincipal, Status: config.StatusFailure, ErrMsg: helpers.Ptr("service role not found")})
		if dberrors.IsNoRows(err) {
			return nil, nil, helpers.ErrNotFound("Service role", fmt.Sprintf("No service role found for principal '%s'", repoPrincipal))
		}
		return nil, nil, err
	}

	projectDelegation, err := s.q.GetDelegation(ctx, serviceRole.ID)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionLogin, ActorType: config.ActorTypeService, ActorID: serviceRole.ID.String(), ActorEmail: repoPrincipal, Status: config.StatusFailure, ErrMsg: helpers.Ptr("service delegation not found")})
		if dberrors.IsNoRows(err) {
			return nil, nil, helpers.ErrNotFound("Delegation", "Ensure the service role is delegated to a project")
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
		s.audit.Log(ctx, AuditEntry{Action: config.ActionLogin, ActorType: config.ActorTypeService, ActorID: serviceRole.ID.String(), ActorEmail: repoPrincipal, ProjectID: &projectDelegation.ProjectID, Environment: &projectDelegation.Env, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		return nil, nil, helpers.ErrInternal(fmt.Sprintf("Failed to create session for %s/%s", projectDelegation.ProjectName, projectDelegation.Env))
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionLogin, ActorType: config.ActorTypeService, ActorID: serviceRole.ID.String(), ActorEmail: repoPrincipal, ProjectID: &projectDelegation.ProjectID, Environment: &projectDelegation.Env, TargetID: helpers.Ptr(session.ID.String()), Status: config.StatusSuccess})

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
					return nil, nil, helpers.ErrConflict(fmt.Sprintf("Refresh token for user %s already exists", userID), "")
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
			return nil, nil, helpers.ErrConflict(fmt.Sprintf("User session for %s already exists", userID), "")
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
			return nil, helpers.ErrUnauthorized("SESSION_EXPIRED", fmt.Sprintf("Session %s is invalid or expired", requestBody.SessionID), "Please log in again")
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
		return nil, helpers.ErrUnauthorized("SESSION_INVALID", fmt.Sprintf("Session %s has incomplete data", requestBody.SessionID), "")
	}

	if (sessionProjectID != requestBody.ProjectID) || (sessionEnv != requestBody.Env) {
		return nil, helpers.ErrForbidden("Project ID and environment do not match session", "")
	}

	projectKeys, err := s.q.GetDelegatedKeys(ctx, database.GetDelegatedKeysParams{
		ProjectID:     sessionProjectID,
		Env:           sessionEnv,
		ServiceRoleID: sessionServiceRoleID,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("Project keys", "Ensure delegation is configured for this service role")
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
			return helpers.ErrUnauthorized("SESSION_EXPIRED", fmt.Sprintf("Session %s is invalid or expired", sessionID), "Please log in again")
		}
		return err
	}

	return nil
}
