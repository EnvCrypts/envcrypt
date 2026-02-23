package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
)

type ProjectService struct {
	q     *database.Queries
	db    *sql.DB
	audit *AuditService
}

func NewProjectService(q *database.Queries) *ProjectService {
	return &ProjectService{q: q}
}

func (s *ProjectService) CreateProject(ctx context.Context, createBody config.ProjectCreateRequest) error {

	creator, err := s.q.GetUserByID(ctx, createBody.UserId)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("User", "")
		}
		return err
	}

	project, err := s.q.CreateProject(ctx, database.CreateProjectParams{
		Name:      createBody.Name,
		CreatedBy: createBody.UserId,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionProjectCreate, ActorType: config.ActorTypeUser, ActorID: createBody.UserId.String(), ActorEmail: creator.Email, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		if dberrors.IsUniqueViolation(err) {
			return helpers.ErrConflict("Project with this name already exists", "Choose a different project name")
		}
		return err
	}

	_, err = s.q.AddUserToProject(ctx, database.AddUserToProjectParams{
		ProjectID: project.ID,
		UserID:    createBody.UserId,
		Role:      "admin",
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionProjectCreate, ActorType: config.ActorTypeUser, ActorID: createBody.UserId.String(), ActorEmail: creator.Email, ProjectID: &project.ID, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		return err
	}

	_, err = s.q.AddWrappedPRK(ctx, database.AddWrappedPRKParams{
		ProjectID:        project.ID,
		UserID:           createBody.UserId,
		WrappedPrk:       createBody.WrappedPRK,
		WrapNonce:        createBody.WrapNonce,
		WrapEphemeralPub: createBody.EphemeralPublicKey,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionProjectCreate, ActorType: config.ActorTypeUser, ActorID: createBody.UserId.String(), ActorEmail: creator.Email, ProjectID: &project.ID, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		return err
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionProjectCreate, ActorType: config.ActorTypeUser, ActorID: createBody.UserId.String(), ActorEmail: creator.Email, ProjectID: &project.ID, Status: config.StatusSuccess})
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

	actor, err := s.q.GetUserByID(ctx, requestBody.UserId)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("User", "")
		}
		return err
	}

	project, err := s.q.GetProject(ctx, database.GetProjectParams{
		Name:      requestBody.ProjectName,
		CreatedBy: requestBody.UserId,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("Project", "Check the project name or your permissions")
		}
		return err
	}

	projectRole, err := s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{UserID: project.CreatedBy, ProjectID: project.ID})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("Project role", "")
		}
		return err
	}
	if projectRole.IsRevoked == true {
		return helpers.ErrForbidden("Your access to this project has been revoked", "Contact the project admin")
	}

	if projectRole.Role != "admin" {
		return helpers.ErrForbidden("Only project admins can perform this action", "")
	}

	err = s.q.DeleteProject(ctx, project.ID)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionProjectDelete, ActorType: config.ActorTypeUser, ActorID: requestBody.UserId.String(), ActorEmail: actor.Email, ProjectID: &project.ID, Status: config.StatusFailure, ErrMsg: helpers.Ptr("unable to delete project")})
		return helpers.ErrInternal("Unable to delete project")
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionProjectDelete, ActorType: config.ActorTypeUser, ActorID: requestBody.UserId.String(), ActorEmail: actor.Email, ProjectID: &project.ID, Status: config.StatusSuccess})
	return nil
}

func (s *ProjectService) AddUserToProject(ctx context.Context, requestBody config.AddUserToProjectRequest) error {

	project, err := s.q.GetProject(ctx, database.GetProjectParams{
		Name:      requestBody.ProjectName,
		CreatedBy: requestBody.AdminId,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("Project", "Check the project name or your permissions")
		}
		return err
	}

	adminUser, err := s.q.GetUserByID(ctx, requestBody.AdminId)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("User", "")
		}
		return err
	}

	projectRole, err := s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{UserID: requestBody.AdminId, ProjectID: project.ID})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("Project role", "")
		}
		return err
	}

	if projectRole.Role != "admin" {
		return helpers.ErrForbidden("Only project admins can perform this action", "")
	}
	if projectRole.IsRevoked == true {
		return helpers.ErrForbidden("Your access to this project has been revoked", "Contact the project admin")
	}

	var role = "member"
	_, err = s.q.AddUserToProject(ctx, database.AddUserToProjectParams{
		ProjectID: project.ID,
		UserID:    requestBody.UserId,
		Role:      role,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionMembershipChange, ActorType: config.ActorTypeUser, ActorID: requestBody.AdminId.String(), ActorEmail: adminUser.Email, ProjectID: &project.ID, TargetID: helpers.Ptr(requestBody.UserId.String()), Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		if dberrors.IsUniqueViolation(err) {
			return helpers.ErrConflict("User is already a member of this project", "")
		}
		return helpers.ErrInternal("Unable to add user to project")
	}

	_, err = s.q.AddWrappedPRK(ctx, database.AddWrappedPRKParams{
		ProjectID:        project.ID,
		UserID:           requestBody.UserId,
		WrappedPrk:       requestBody.WrappedPRK,
		WrapNonce:        requestBody.WrapNonce,
		WrapEphemeralPub: requestBody.EphemeralPublicKey,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionMembershipChange, ActorType: config.ActorTypeUser, ActorID: requestBody.AdminId.String(), ActorEmail: adminUser.Email, ProjectID: &project.ID, TargetID: helpers.Ptr(requestBody.UserId.String()), Status: config.StatusFailure, ErrMsg: helpers.Ptr("unable to add wrapped prk")})
		return helpers.ErrInternal("Unable to add wrapped PRK")
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionMembershipChange, ActorType: config.ActorTypeUser, ActorID: requestBody.AdminId.String(), ActorEmail: adminUser.Email, ProjectID: &project.ID, TargetID: helpers.Ptr(requestBody.UserId.String()), Status: config.StatusSuccess})

	return nil
}

func (s *ProjectService) SetUserAccess(ctx context.Context, requestBody config.SetAccessRequest) error {

	project, err := s.q.GetProject(ctx, database.GetProjectParams{
		Name:      requestBody.ProjectName,
		CreatedBy: requestBody.AdminId,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("Project", "Check the project name or your permissions")
		}
		return err
	}

	adminUser, err := s.q.GetUserByID(ctx, requestBody.AdminId)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("User", "")
		}
		return err
	}

	projectRole, err := s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{UserID: requestBody.AdminId, ProjectID: project.ID})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("Project role", "")
		}
		return err
	}

	if projectRole.Role != "admin" {
		return helpers.ErrForbidden("Only project admins can perform this action", "")
	}
	if projectRole.IsRevoked == true {
		return helpers.ErrForbidden("Your access to this project has been revoked", "Contact the project admin")
	}

	user, err := s.q.GetUserByEmail(ctx, requestBody.UserEmail)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("User", "Check the email address")
		}
		return err
	}

	err = s.q.SetUserAccess(ctx, database.SetUserAccessParams{
		UserID:    user.ID,
		ProjectID: project.ID,
		IsRevoked: requestBody.IsRevoked,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionMembershipChange, ActorType: config.ActorTypeUser, ActorID: requestBody.AdminId.String(), ActorEmail: adminUser.Email, ProjectID: &project.ID, TargetID: helpers.Ptr(user.ID.String()), Status: config.StatusFailure, ErrMsg: helpers.Ptr("unable to revoke user access")})
		return helpers.ErrInternal("Unable to update user access")
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionMembershipChange, ActorType: config.ActorTypeUser, ActorID: requestBody.AdminId.String(), ActorEmail: adminUser.Email, ProjectID: &project.ID, TargetID: helpers.Ptr(user.ID.String()), Status: config.StatusSuccess})

	return nil
}

func (s *ProjectService) GetUserProject(ctx context.Context, requestBody config.GetUserProjectRequest) (*config.GetUserProjectResponse, error) {

	project, err := s.q.GetProject(ctx, database.GetProjectParams{
		Name:      requestBody.ProjectName,
		CreatedBy: requestBody.UserId,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("Project", "Check the project name or your permissions")
		}
		return nil, err
	}

	wrappedKey, err := s.q.GetProjectWrappedKey(ctx, database.GetProjectWrappedKeyParams{
		ProjectID: project.ID,
		UserID:    requestBody.UserId,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("Project key", "")
		}
		return nil, err
	}

	var response = &config.GetUserProjectResponse{
		ProjectId:          project.ID,
		WrappedPRK:         wrappedKey.WrappedPrk,
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
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("Project", "Check the project name or your permissions")
		}
		return nil, err
	}

	wrappedKey, err := s.q.GetProjectWrappedKey(ctx, database.GetProjectWrappedKeyParams{
		ProjectID: project,
		UserID:    requestBody.UserId,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("Project key", "")
		}
		return nil, err
	}

	var response = &config.GetMemberProjectResponse{
		ProjectId:          project,
		WrappedPRK:         wrappedKey.WrappedPrk,
		WrapNonce:          wrappedKey.WrapNonce,
		EphemeralPublicKey: wrappedKey.WrapEphemeralPub,
	}

	return response, nil
}

func (s *ProjectService) RotateInit(ctx context.Context, req config.RotateInitRequest) (*config.RotateInitResponse, error) {
	actor, err := s.q.GetUserByID(ctx, req.UserID)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("User", "")
		}
		return nil, err
	}

	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
		IsRevoked: false,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrForbidden("User doesn't have permission for this project", "")
		}
		return nil, err
	}

	project, err := s.q.GetProjectById(ctx, req.ProjectID)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("Project", "")
		}
		return nil, err
	}

	rotationData, err := s.q.GetRotationData(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}

	wrappedDEKs, err := s.q.GetProjectWrappedDEKs(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}

	resp := &config.RotateInitResponse{
		WrappedPRKs:      make([]config.WrappedKey, len(rotationData)),
		MemberPublicKeys: make([]config.MemberPublicKey, len(rotationData)),
		WrappedDEKs:      make([]config.WrappedDEK, len(wrappedDEKs)),
		PRKVersion:       project.PrkVersion,
	}

	for i, row := range rotationData {
		resp.WrappedPRKs[i] = config.WrappedKey{
			UserID:             row.UserID,
			WrappedPRK:         row.WrappedPrk,
			WrapNonce:          row.WrapNonce,
			EphemeralPublicKey: row.WrapEphemeralPub,
		}
		resp.MemberPublicKeys[i] = config.MemberPublicKey{
			UserID:    row.UserID,
			PublicKey: row.UserPublicKey,
		}
	}

	for i, dek := range wrappedDEKs {
		resp.WrappedDEKs[i] = config.WrappedDEK{
			EnvVersionID: dek.ID,
			WrappedDEK:   dek.WrappedDek,
			DekNonce:     dek.DekNonce,
		}
	}

	s.audit.Log(ctx, AuditEntry{
		Action:    config.ActionPRKRotate,
		ActorType: config.ActorTypeUser,
		ActorID:   req.UserID.String(),
		ActorEmail: actor.Email,
		ProjectID: &req.ProjectID,
		Status:    config.StatusSuccess,
		Metadata:  mustJSON(map[string]any{"phase": "init", "prk_version": project.PrkVersion}),
	})

	return resp, nil
}

func (s *ProjectService) RotateCommit(ctx context.Context, req config.RotateCommitRequest) (*config.RotateCommitResponse, error) {
	actor, err := s.q.GetUserByID(ctx, req.UserID)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("User", "")
		}
		return nil, err
	}

	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
		IsRevoked: false,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrForbidden("User doesn't have permission for this project", "")
		}
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	txQ := s.q.WithTx(tx)

	newVersion, err := txQ.IncrementPRKVersion(ctx, database.IncrementPRKVersionParams{
		ID:         req.ProjectID,
		PrkVersion: req.ExpectedPRKVersion,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrConflict("PRK version conflict: another rotation is in progress", "Retry the rotation")
		}
		return nil, err
	}

	for _, wrappedPRK := range req.NewWrappedPRKs {
		err = txQ.UpdateWrappedPRK(ctx, database.UpdateWrappedPRKParams{
			ProjectID:        req.ProjectID,
			UserID:           wrappedPRK.UserID,
			WrappedPrk:       wrappedPRK.WrappedPRK,
			WrapNonce:        wrappedPRK.WrapNonce,
			WrapEphemeralPub: wrappedPRK.EphemeralPublicKey,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update wrapped prk: %w", err)
		}
	}

	for _, dek := range req.NewWrappedDEKs {
		err = txQ.UpdateEnvVersionDEK(ctx, database.UpdateEnvVersionDEKParams{
			ID:         dek.EnvVersionID,
			WrappedDek: dek.NewWrappedDEK,
			DekNonce:   dek.NewDekNonce,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update wrapped dek: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.audit.Log(ctx, AuditEntry{
		Action:     config.ActionPRKRotate,
		ActorType:  config.ActorTypeUser,
		ActorID:    req.UserID.String(),
		ActorEmail: actor.Email,
		ProjectID:  &req.ProjectID,
		Status:     config.StatusSuccess,
		Metadata: mustJSON(map[string]any{
			"phase":              "commit",
			"old_prk_version":    req.ExpectedPRKVersion,
			"new_prk_version":    newVersion,
			"versions_rewrapped": len(req.NewWrappedDEKs),
		}),
	})

	return &config.RotateCommitResponse{NewPRKVersion: newVersion}, nil
}

func mustJSON(v any) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}
