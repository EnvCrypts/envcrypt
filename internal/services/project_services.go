package services

import (
	"context"
	"errors"
	"log"

	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
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

	wrappedPmk, err := s.q.AddWrappedPMK(ctx, database.AddWrappedPMKParams{
		ProjectID:        project.ID,
		UserID:           createBody.UserId,
		WrappedPmk:       createBody.WrappedPMK,
		WrapNonce:        createBody.WrapNonce,
		WrapEphemeralPub: createBody.EphemeralPublicKey,
	})
	if err != nil {
		return err
	}

	log.Print(wrappedPmk)

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

	wrappedPmk, err := s.q.AddWrappedPMK(ctx, database.AddWrappedPMKParams{
		ProjectID:        requestBody.ProjectId,
		UserID:           requestBody.UserId,
		WrappedPmk:       requestBody.WrappedPMK,
		WrapNonce:        requestBody.WrapNonce,
		WrapEphemeralPub: requestBody.EphemeralPublicKey,
	})
	if err != nil {
		return err
	}

	log.Print(wrappedPmk)

	return nil
}
