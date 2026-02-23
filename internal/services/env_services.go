package services

import (
	"context"
	"encoding/json"

	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
)

type EnvServices struct {
	q     *database.Queries
	audit *AuditService
}

func NewEnvService(q *database.Queries) *EnvServices {
	return &EnvServices{
		q: q,
	}
}

func (s *EnvServices) GetEnv(ctx context.Context, requestBody config.GetEnvRequest) (*config.GetEnvResponse, error) {

	user, err := s.q.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPull, ActorType: config.ActorTypeUser, ActorID: "unknown", ActorEmail: requestBody.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusFailure, ErrMsg: helpers.Ptr("user not found")})
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("User", "Check the email address")
		}
		return nil, err
	}

	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID:    user.ID,
		ProjectID: requestBody.ProjectId,
		IsRevoked: false,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPull, ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: requestBody.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusFailure, ErrMsg: helpers.Ptr("permission denied")})
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrForbidden("You don't have permission to access this environment", "")
		}
		return nil, err
	}

	var env database.EnvVersion
	if requestBody.Version != nil {
		env, err = s.q.GetEnv(ctx, database.GetEnvParams{
			ProjectID: requestBody.ProjectId,
			EnvName:   requestBody.EnvName,
			Version:   *requestBody.Version,
		})
		if err != nil {
			if dberrors.IsNoRows(err) {
				return nil, helpers.ErrNotFound("Environment", "Check the environment name and version")
			}
			return nil, err
		}
	} else {
		env, err = s.q.GetLatestEnv(ctx, database.GetLatestEnvParams{
			ProjectID: requestBody.ProjectId,
			EnvName:   requestBody.EnvName,
		})
		if err != nil {
			s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPull, ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: requestBody.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusFailure, ErrMsg: helpers.Ptr("env not found")})
			if dberrors.IsNoRows(err) {
				return nil, helpers.ErrNotFound("Environment", "Check the environment name")
			}
			return nil, err
		}
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPull, ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: requestBody.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusSuccess})

	return &config.GetEnvResponse{
		CipherText:        env.Ciphertext,
		Nonce:             env.Nonce,
		WrappedDEK:        env.WrappedDek,
		DekNonce:          env.DekNonce,
		EncryptionVersion: env.EncryptionVersion,
	}, nil
}

func (s *EnvServices) GetEnvVersions(ctx context.Context, requestBody config.GetEnvVersionsRequest) (*config.GetEnvVersionsResponse, error) {
	user, err := s.q.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPull, ActorType: config.ActorTypeUser, ActorID: "unknown", ActorEmail: requestBody.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusFailure, ErrMsg: helpers.Ptr("user not found")})
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("User", "Check the email address")
		}
		return nil, err
	}

	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID:    user.ID,
		ProjectID: requestBody.ProjectId,
		IsRevoked: false,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPull, ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: requestBody.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusFailure, ErrMsg: helpers.Ptr("permission denied")})
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrForbidden("You don't have permission to access this environment", "")
		}
		return nil, err
	}

	envVersions, err := s.q.GetEnvVersions(ctx, database.GetEnvVersionsParams{
		ProjectID: requestBody.ProjectId,
		EnvName:   requestBody.EnvName,
	})
	if err != nil {
		return nil, err
	}

	var envResponses []config.EnvResponse

	for _, envVersion := range envVersions {
		var metadata config.Metadata
		err = json.Unmarshal(envVersion.Metadata, &metadata)
		if err != nil {
			continue
		}

		envResponses = append(envResponses, config.EnvResponse{
			CipherText:        envVersion.Ciphertext,
			Nonce:             envVersion.Nonce,
			WrappedDEK:        envVersion.WrappedDek,
			DekNonce:          envVersion.DekNonce,
			EncryptionVersion: envVersion.EncryptionVersion,
			Version:           envVersion.Version,
			Metadata:          metadata,
		})
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPull, ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: requestBody.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusSuccess})

	return &config.GetEnvVersionsResponse{EnvVersions: envResponses}, nil
}

func (s *EnvServices) AddEnv(ctx context.Context, requestBody config.AddEnvRequest) error {

	user, err := s.q.GetUserByID(ctx, requestBody.UserId)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPush, ActorType: config.ActorTypeUser, ActorID: requestBody.UserId.String(), ActorEmail: "unknown", ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusFailure, ErrMsg: helpers.Ptr("user not found")})
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("User", "")
		}
		return err
	}

	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID:    requestBody.UserId,
		ProjectID: requestBody.ProjectId,
		IsRevoked: false,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPush, ActorType: config.ActorTypeUser, ActorID: requestBody.UserId.String(), ActorEmail: user.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusFailure, ErrMsg: helpers.Ptr("user doesn't have permission to store env")})
		if dberrors.IsNoRows(err) {
			return helpers.ErrForbidden("You don't have permission to push to this environment", "")
		}
		return err
	}

	metadata, err := json.Marshal(requestBody.Metadata)
	if err != nil {
		return err
	}

	_, err = s.q.AddEnv(ctx, database.AddEnvParams{
		ProjectID:         requestBody.ProjectId,
		EnvName:           requestBody.EnvName,
		Ciphertext:        requestBody.CipherText,
		Nonce:             requestBody.Nonce,
		WrappedDek:        requestBody.WrappedDEK,
		DekNonce:          requestBody.DekNonce,
		EncryptionVersion: requestBody.EncryptionVersion,
		CreatedBy:         requestBody.UserId,
		Metadata:          metadata,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPush, ActorType: config.ActorTypeUser, ActorID: requestBody.UserId.String(), ActorEmail: user.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		if dberrors.IsUniqueViolation(err) {
			return helpers.ErrConflict("Environment with this version already exists", "")
		}
		return err
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPush, ActorType: config.ActorTypeUser, ActorID: requestBody.UserId.String(), ActorEmail: user.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusSuccess, Metadata: metadata})

	return nil
}

func (s *EnvServices) UpdateEnv(ctx context.Context, requestBody config.UpdateEnvRequest) error {
	user, err := s.q.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("User", "Check the email address")
		}
		return err
	}

	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID:    user.ID,
		ProjectID: requestBody.ProjectId,
		IsRevoked: false,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPush, ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: requestBody.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusFailure, ErrMsg: helpers.Ptr("permission denied")})
		if dberrors.IsNoRows(err) {
			return helpers.ErrForbidden("You don't have permission to update this environment", "")
		}
		return err
	}

	metadata, err := json.Marshal(requestBody.Metadata)
	if err != nil {
		return err
	}

	_, err = s.q.AddEnv(ctx, database.AddEnvParams{
		ProjectID:         requestBody.ProjectId,
		EnvName:           requestBody.EnvName,
		Ciphertext:        requestBody.CipherText,
		Nonce:             requestBody.Nonce,
		WrappedDek:        requestBody.WrappedDEK,
		DekNonce:          requestBody.DekNonce,
		EncryptionVersion: requestBody.EncryptionVersion,
		CreatedBy:         user.ID,
		Metadata:          metadata,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPush, ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: requestBody.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		return err
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPush, ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: requestBody.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusSuccess, Metadata: metadata})

	return nil
}

func (s *EnvServices) GetEnvForCI(ctx context.Context, requestBody config.GetEnvForCIRequest) (*config.GetEnvForCIResponse, error) {
	env, err := s.q.GetLatestEnv(ctx, database.GetLatestEnvParams{
		ProjectID: requestBody.ProjectId,
		EnvName:   requestBody.EnvName,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPull, ActorType: config.ActorTypeService, ActorID: "ci_session", ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusFailure, ErrMsg: helpers.Ptr("env not found")})
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("Environment", "Check the environment name")
		}
		return nil, err
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionEnvPull, ActorType: config.ActorTypeService, ActorID: "ci_session", ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, Status: config.StatusSuccess})

	return &config.GetEnvForCIResponse{
		CipherText:        env.Ciphertext,
		Nonce:             env.Nonce,
		WrappedDEK:        env.WrappedDek,
		DekNonce:          env.DekNonce,
		EncryptionVersion: env.EncryptionVersion,
	}, nil
}
