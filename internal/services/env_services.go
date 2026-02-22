package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
)

type EnvServices struct {
	q *database.Queries
}

func NewEnvService(q *database.Queries) *EnvServices {
	return &EnvServices{
		q: q,
	}
}

func (s *EnvServices) GetEnv(ctx context.Context, requestBody config.GetEnvRequest) (*config.GetEnvResponse, error) {

	user, err := s.q.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, errors.New("user not found")
		}
		log.Println("Error getting user")
		return nil, err
	}

	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID:    user.ID,
		ProjectID: requestBody.ProjectId,
		IsRevoked: false,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, errors.New("user doesn't have permission to get env")
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
				return nil, errors.New("env not found")
			}
			log.Print(err.Error())
			return nil, err
		}
	} else {
		env, err = s.q.GetLatestEnv(ctx, database.GetLatestEnvParams{
			ProjectID: requestBody.ProjectId,
			EnvName:   requestBody.EnvName,
		})
		if err != nil {
			if dberrors.IsNoRows(err) {
				return nil, errors.New("env not found")
			}
			log.Print(err.Error())
			return nil, err
		}
	}

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
		if dberrors.IsNoRows(err) {
			return nil, errors.New("user not found")
		}
		log.Println("Error getting user")
		return nil, err
	}

	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID:    user.ID,
		ProjectID: requestBody.ProjectId,
		IsRevoked: false,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, errors.New("user doesn't have permission to get env")
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
			log.Print(err.Error())
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

	return &config.GetEnvVersionsResponse{EnvVersions: envResponses}, nil
}

func (s *EnvServices) AddEnv(ctx context.Context, requestBody config.AddEnvRequest) error {

	_, err := s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID:    requestBody.UserId,
		ProjectID: requestBody.ProjectId,
		IsRevoked: false,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return errors.New("user doesn't have permission to store env")
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
		if dberrors.IsUniqueViolation(err) {
			return errors.New("env with this version already exists")
		}
		return err
	}

	return nil
}

func (s *EnvServices) UpdateEnv(ctx context.Context, requestBody config.UpdateEnvRequest) error {
	user, err := s.q.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return errors.New("user not found")
		}
		return err
	}

	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID:    user.ID,
		ProjectID: requestBody.ProjectId,
		IsRevoked: false,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return errors.New("user doesn't have permission to update env")
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
		return err
	}

	return nil
}

func (s *EnvServices) GetEnvForCI(ctx context.Context, requestBody config.GetEnvForCIRequest) (*config.GetEnvForCIResponse, error) {
	env, err := s.q.GetLatestEnv(ctx, database.GetLatestEnvParams{
		ProjectID: requestBody.ProjectId,
		EnvName:   requestBody.EnvName,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, errors.New("env not found")
		}
		return nil, err
	}

	return &config.GetEnvForCIResponse{
		CipherText:        env.Ciphertext,
		Nonce:             env.Nonce,
		WrappedDEK:        env.WrappedDek,
		DekNonce:          env.DekNonce,
		EncryptionVersion: env.EncryptionVersion,
	}, nil
}
