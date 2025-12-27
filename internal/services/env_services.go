package services

import (
	"context"
	"errors"
	"log"

	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
)

type EnvServices struct {
	q *database.Queries
}

func NewEnvServices(q *database.Queries) *EnvServices {
	return &EnvServices{
		q: q,
	}
}

func (s *EnvServices) GetEnv(ctx context.Context, requestBody config.GetEnvRequest) (*config.GetEnvResponse, error) {

	user, err := s.q.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		log.Println("Error getting user")
		return nil, err
	}

	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID:    user.ID,
		ProjectID: requestBody.ProjectId,
	})
	if err != nil {
		return nil, errors.New("user doesn't have permission to get env")
	}

	env, err := s.q.GetEnv(ctx, database.GetEnvParams{
		ProjectID: requestBody.ProjectId,
		EnvName:   requestBody.EnvName,
		Version:   requestBody.Version,
	})
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return &config.GetEnvResponse{
		CipherText: env.Ciphertext,
		Nonce:      env.Nonce,
	}, nil
}

func (s *EnvServices) AddEnv(ctx context.Context, requestBody config.AddEnvRequest) error {
	user, err := s.q.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		return err
	}

	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID:    user.ID,
		ProjectID: requestBody.ProjectId,
	})
	if err != nil {
		return errors.New("user doesn't have permission to store env")
	}

	_, err = s.q.AddEnv(ctx, database.AddEnvParams{
		ProjectID:  requestBody.ProjectId,
		EnvName:    requestBody.EnvName,
		Version:    requestBody.Version,
		Ciphertext: requestBody.CipherText,
		Nonce:      requestBody.Nonce,
		CreatedBy:  user.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *EnvServices) UpdateEnv(ctx context.Context, requestBody config.UpdateEnvRequest) error {
	user, err := s.q.GetUserByEmail(ctx, requestBody.Email)
	if err != nil {
		return err
	}

	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID: user.ID,
	})
	if err != nil {
		return errors.New("user doesn't have permission to update env")
	}

	_, err = s.q.UpdateEnv(ctx, database.UpdateEnvParams{
		ProjectID:  requestBody.ProjectId,
		EnvName:    requestBody.EnvName,
		Version:    requestBody.Version,
		Ciphertext: requestBody.CipherText,
		Nonce:      requestBody.Nonce,
	})
	if err != nil {
		return err
	}

	return nil
}
