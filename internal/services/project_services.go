package services

import (
	"context"
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
