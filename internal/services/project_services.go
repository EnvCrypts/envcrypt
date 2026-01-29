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

func (s *ProjectService) ListProjects(ctx context.Context, requestBody config.ListProjectRequest) (*config.ListProjectResponse, error) {

	projects, err := s.q.ListProjectsWithRole(ctx, requestBody.UserId)
	if err != nil {
		return nil, err
	}

	resp := &config.ListProjectResponse{
		Projects: make([]config.Project, len(projects)),
	}

	for i, project := range projects {
		resp.Projects[i] = config.Project{
			Id:        project.ID,
			Name:      project.Name,
			Role:      project.Role,
			IsRevoked: project.IsRevoked,
		}
	}

	return resp, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, requestBody config.ProjectDeleteRequest) error {

	project, err := s.q.GetProject(ctx, database.GetProjectParams{
		Name:      requestBody.ProjectName,
		CreatedBy: requestBody.UserId,
	})
	if err != nil {
		return errors.New("project not found")
	}

	projectRole, err := s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{UserID: project.CreatedBy, ProjectID: project.ID})
	if err != nil {
		return errors.New("project role not found")
	}
	if projectRole.IsRevoked == true {
		return errors.New("user access is revoked")
	}

	if projectRole.Role != "admin" {
		return errors.New("user is not an admin")
	}

	err = s.q.DeleteProject(ctx, project.ID)
	if err != nil {
		return errors.New("unable to delete project")
	}

	return nil
}

func (s *ProjectService) AddUserToProject(ctx context.Context, requestBody config.AddUserToProjectRequest) error {

	project, err := s.q.GetProject(ctx, database.GetProjectParams{
		Name:      requestBody.ProjectName,
		CreatedBy: requestBody.AdminId,
	})
	if err != nil {
		return errors.New("project not found")
	}

	projectRole, err := s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{UserID: requestBody.AdminId, ProjectID: project.ID})
	if err != nil {
		return errors.New("user role not found")
	}

	if projectRole.Role != "admin" {
		return errors.New("user is not an admin")
	}
	if projectRole.IsRevoked == true {
		return errors.New("user access is revoked")
	}

	var role = "member"
	if requestBody.Role == "admin" {
		role = "admin"
	}
	_, err = s.q.AddUserToProject(ctx, database.AddUserToProjectParams{
		ProjectID: project.ID,
		UserID:    requestBody.UserId,
		Role:      role,
	})
	if err != nil {
		if dberrors.IsUniqueViolation(err) {
			return errors.New("project already has user")
		}
		return errors.New("unable to add user to project")
	}

	_, err = s.q.AddWrappedPMK(ctx, database.AddWrappedPMKParams{
		ProjectID:        project.ID,
		UserID:           requestBody.UserId,
		WrappedPmk:       requestBody.WrappedPMK,
		WrapNonce:        requestBody.WrapNonce,
		WrapEphemeralPub: requestBody.EphemeralPublicKey,
	})
	if err != nil {
		return errors.New("unable to add wrapped pmk")
	}

	return nil
}

func (s *ProjectService) SetUserAccess(ctx context.Context, requestBody config.SetAccessRequest) error {

	project, err := s.q.GetProject(ctx, database.GetProjectParams{
		Name:      requestBody.ProjectName,
		CreatedBy: requestBody.AdminId,
	})
	if err != nil {
		return errors.New("project not found")
	}

	projectRole, err := s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{UserID: requestBody.AdminId, ProjectID: project.ID})
	if err != nil {
		return errors.New("user role not found")
	}

	if projectRole.Role != "admin" {
		return errors.New("user is not an admin")
	}
	if projectRole.IsRevoked == true {
		return errors.New("user access is revoked")
	}

	user, err := s.q.GetUserByEmail(ctx, requestBody.UserEmail)
	if err != nil {
		return errors.New("user not found")
	}

	err = s.q.SetUserAccess(ctx, database.SetUserAccessParams{
		UserID:    user.ID,
		ProjectID: project.ID,
		IsRevoked: requestBody.IsRevoked,
	})
	if err != nil {
		return errors.New("unable to revoke user access")
	}

	return nil
}

func (s *ProjectService) GetUserProject(ctx context.Context, requestBody config.GetUserProjectRequest) (*config.GetUserProjectResponse, error) {

	project, err := s.q.GetProject(ctx, database.GetProjectParams{
		Name:      requestBody.ProjectName,
		CreatedBy: requestBody.UserId,
	})
	if err != nil {
		return nil, errors.New("project not found")
	}

	wrappedKey, err := s.q.GetProjectWrappedKey(ctx, database.GetProjectWrappedKeyParams{
		ProjectID: project.ID,
		UserID:    requestBody.UserId,
	})
	if err != nil {
		return nil, errors.New("project wrapped key not found")
	}

	var response = &config.GetUserProjectResponse{
		ProjectId:          project.ID,
		WrappedPMK:         wrappedKey.WrappedPmk,
		WrapNonce:          wrappedKey.WrapNonce,
		EphemeralPublicKey: wrappedKey.WrapEphemeralPub,
	}

	return response, nil
}

func (s *ProjectService) GetMemberProject(ctx context.Context, requestBody config.GetMemberProjectRequest) (*config.GetMemberProjectResponse, error) {
	project, err := s.q.GetMemberProject(ctx, database.GetMemberProjectParams{
		Name:   requestBody.ProjectName,
		UserID: requestBody.UserId,
	})
	if err != nil {
		return nil, errors.New("project not found")
	}

	return &config.GetMemberProjectResponse{
		ProjectId: project,
	}, nil
}
