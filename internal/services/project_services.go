package services

import (
	"context"
	"errors"

	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
)

type ProjectService struct {
	q *database.Queries
}

func NewProjectService(q *database.Queries) *ProjectService {
	return &ProjectService{q: q}
}

func (s *ProjectService) CreateProject(ctx context.Context, createBody config.ProjectCreateRequest) error {

	project, err := s.q.CreateProject(ctx, database.CreateProjectParams{
		Name:      createBody.Name,
		CreatedBy: createBody.UserId,
	})
	if err != nil {
		if dberrors.IsUniqueViolation(err) {
			return errors.New("project already exists")
		}
		return err
	}

	_, err = s.q.AddUserToProject(ctx, database.AddUserToProjectParams{
		ProjectID: project.ID,
		UserID:    createBody.UserId,
		Role:      "admin",
	})
	if err != nil {
		return err
	}

	_, err = s.q.AddWrappedPMK(ctx, database.AddWrappedPMKParams{
		ProjectID:        project.ID,
		UserID:           createBody.UserId,
		WrappedPmk:       createBody.WrappedPMK,
		WrapNonce:        createBody.WrapNonce,
		WrapEphemeralPub: createBody.EphemeralPublicKey,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) AddUserToProject(ctx context.Context, requestBody config.AddUserToProjectRequest) error {

	projectRole, err := s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{UserID: requestBody.AdminId, ProjectID: requestBody.ProjectId})
	if err != nil {
		return err
	}

	if projectRole.Role != "admin" {
		return errors.New("user is not an admin")
	}

	_, err = s.q.AddUserToProject(ctx, database.AddUserToProjectParams{
		ProjectID: requestBody.ProjectId,
		UserID:    requestBody.UserId,
		Role:      "member",
	})
	if err != nil {
		return err
	}

	_, err = s.q.AddWrappedPMK(ctx, database.AddWrappedPMKParams{
		ProjectID:        requestBody.ProjectId,
		UserID:           requestBody.UserId,
		WrappedPmk:       requestBody.WrappedPMK,
		WrapNonce:        requestBody.WrapNonce,
		WrapEphemeralPub: requestBody.EphemeralPublicKey,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) GetUserProject(ctx context.Context, requestBody config.GetUserProjectRequest) (*config.GetUserProjectResponse, error) {

	project, err := s.q.GetProject(ctx, database.GetProjectParams{
		Name:      requestBody.ProjectName,
		CreatedBy: requestBody.UserId,
	})
	if err != nil {
		return nil, err
	}

	wrappedKey, err := s.q.GetProjectWrappedKey(ctx, database.GetProjectWrappedKeyParams{
		ProjectID: project.ID,
		UserID:    requestBody.UserId,
	})
	if err != nil {
		return nil, err
	}

	var response = &config.GetUserProjectResponse{
		ProjectId:          project.ID,
		WrappedPMK:         wrappedKey.WrappedPmk,
		WrapNonce:          wrappedKey.WrapNonce,
		EphemeralPublicKey: wrappedKey.WrapEphemeralPub,
	}

	return response, nil
}
