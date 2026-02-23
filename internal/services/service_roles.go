package services

import (
	"context"
	"fmt"

	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
)

type ServiceRoleServices struct {
	q     *database.Queries
	audit *AuditService
}

func NewServiceRoleService(q *database.Queries) *ServiceRoleServices {
	return &ServiceRoleServices{
		q: q,
	}
}

func (s *ServiceRoleServices) List(ctx context.Context, requestBody config.ServiceRoleListRequest) (*config.ServiceRoleListResponse, error) {

	serviceRolesDB, err := s.q.GetServiceRolesByAdmin(ctx, requestBody.CreatedBy)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("Service role", "")
		}
		return nil, err
	}

	serviceRoles := make([]config.ServiceRole, len(serviceRolesDB))
	for i := range serviceRolesDB {
		serviceRoles[i] = config.ServiceRole{
			ID:                   serviceRolesDB[i].ID,
			Name:                 serviceRolesDB[i].Name,
			ServiceRolePublicKey: serviceRolesDB[i].ServiceRolePublicKey,
			RepoPrincipal:        serviceRolesDB[i].RepoPrincipal,
			CreatedAt:            serviceRolesDB[i].CreatedAt,
			CreatedBy:            serviceRolesDB[i].CreatedBy,
		}
	}

	return &config.ServiceRoleListResponse{ServiceRoles: serviceRoles}, nil
}

func (s *ServiceRoleServices) Create(ctx context.Context, requestBody config.ServiceRoleCreateRequest) (*config.ServiceRoleCreateResponse, error) {

	creator, err := s.q.GetUserByID(ctx, requestBody.CreatedBy)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("User", "")
		}
		return nil, err
	}

	serviceRole, err := s.q.CreateServiceRole(ctx, database.CreateServiceRoleParams{
		Name:                 requestBody.ServiceRoleName,
		ServiceRolePublicKey: requestBody.ServiceRolePublicKey,
		RepoPrincipal:        requestBody.RepoPrincipal,
		CreatedBy:            requestBody.CreatedBy,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionServiceRoleCreate, ActorType: config.ActorTypeUser, ActorID: requestBody.CreatedBy.String(), ActorEmail: creator.Email, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		if dberrors.IsUniqueViolation(err) == true {
			return nil, helpers.ErrConflict("Service role already exists", "Choose a different name or principal")
		}
		return nil, err
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionServiceRoleCreate, ActorType: config.ActorTypeUser, ActorID: requestBody.CreatedBy.String(), ActorEmail: creator.Email, TargetID: helpers.Ptr(serviceRole.ID.String()), Status: config.StatusSuccess})

	return &config.ServiceRoleCreateResponse{
		Message: fmt.Sprintf("service role '%s' created", serviceRole.Name),
		ServiceRole: config.ServiceRole{
			ID:                   serviceRole.ID,
			Name:                 serviceRole.Name,
			ServiceRolePublicKey: requestBody.ServiceRolePublicKey,
			RepoPrincipal:        requestBody.RepoPrincipal,
			CreatedBy:            requestBody.CreatedBy,
			CreatedAt:            serviceRole.CreatedAt,
		},
	}, nil
}

func (s *ServiceRoleServices) Get(ctx context.Context, requestBody config.ServiceRoleGetRequest) (*config.ServiceRoleGetResponse, error) {

	serviceRole, err := s.q.GetServiceRoleByPrincipal(ctx, requestBody.RepoPrincipal)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("Service role", "")
		}
		return nil, err
	}

	return &config.ServiceRoleGetResponse{
		ServiceRole: config.ServiceRole{
			ID:                   serviceRole.ID,
			Name:                 serviceRole.Name,
			ServiceRolePublicKey: serviceRole.ServiceRolePublicKey,
			RepoPrincipal:        serviceRole.RepoPrincipal,
			CreatedBy:            serviceRole.CreatedBy,
			CreatedAt:            serviceRole.CreatedAt,
		},
		Message: fmt.Sprintf("service role '%s retrieved'", serviceRole.Name),
	}, nil
}

func (s *ServiceRoleServices) Delete(ctx context.Context, requestBody config.ServiceRoleDeleteRequest) error {

	actor, err := s.q.GetUserByID(ctx, requestBody.CreatedBy)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("User", "")
		}
		return err
	}

	serviceRole, err := s.q.GetServiceRoleById(ctx, requestBody.ServiceRoleId)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("Service role", "")
		}
		return err
	}
	if requestBody.CreatedBy != serviceRole.CreatedBy {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionServiceRoleDelete, ActorType: config.ActorTypeUser, ActorID: requestBody.CreatedBy.String(), ActorEmail: actor.Email, TargetID: helpers.Ptr(serviceRole.ID.String()), Status: config.StatusFailure, ErrMsg: helpers.Ptr("not authorized to delete service role")})
		return helpers.ErrForbidden("Not authorized to delete this service role", "Only the creator can delete it")
	}

	_, err = s.q.DeleteServiceRole(ctx, serviceRole.ID)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionServiceRoleDelete, ActorType: config.ActorTypeUser, ActorID: requestBody.CreatedBy.String(), ActorEmail: actor.Email, TargetID: helpers.Ptr(serviceRole.ID.String()), Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		return err
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionServiceRoleDelete, ActorType: config.ActorTypeUser, ActorID: requestBody.CreatedBy.String(), ActorEmail: actor.Email, TargetID: helpers.Ptr(serviceRole.ID.String()), Status: config.StatusSuccess})
	return nil
}

func (s *ServiceRoleServices) DelegateAccess(ctx context.Context, requestBody config.ServiceRoleDelegateRequest) error {

	actor, err := s.q.GetUserByID(ctx, requestBody.DelegatedBy)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("User", "")
		}
		return err
	}

	projectRole, err := s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID:    requestBody.DelegatedBy,
		ProjectID: requestBody.ProjectId,
	})
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

	serviceRole, err := s.q.GetServiceRoleByPrincipal(ctx, requestBody.RepoPrincipal)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return helpers.ErrNotFound("Service role", "")
		}
		return err
	}

	_, err = s.q.DelegateAccess(ctx, database.DelegateAccessParams{
		ServiceRoleID:    serviceRole.ID,
		ProjectID:        requestBody.ProjectId,
		Env:              requestBody.EnvName,
		WrappedPrk:       requestBody.WrappedPRK,
		WrapNonce:        requestBody.WrapNonce,
		WrapEphemeralPub: requestBody.EphemeralPublicKey,
		DelegatedBy:      requestBody.DelegatedBy,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionServiceRoleDelegate, ActorType: config.ActorTypeUser, ActorID: requestBody.DelegatedBy.String(), ActorEmail: actor.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, TargetID: helpers.Ptr(serviceRole.ID.String()), Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		if dberrors.IsUniqueViolation(err) {
			return helpers.ErrConflict("Service role is already delegated to a project", "")
		}
		return err
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionServiceRoleDelegate, ActorType: config.ActorTypeUser, ActorID: requestBody.DelegatedBy.String(), ActorEmail: actor.Email, ProjectID: &requestBody.ProjectId, Environment: &requestBody.EnvName, TargetID: helpers.Ptr(serviceRole.ID.String()), Status: config.StatusSuccess})
	return nil
}

func (s *ServiceRoleServices) GetPerms(ctx context.Context, requestBody config.ServiceRolePermsRequest) (*config.ServiceRolePermsResponse, error) {

	serviceRole, err := s.q.GetServiceRoleByPrincipal(ctx, requestBody.RepoPrincipal)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("Service role", "")
		}
		return nil, err
	}

	projectDelegated, err := s.q.GetDelegation(ctx, serviceRole.ID)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, helpers.ErrNotFound("Delegation", "Ensure the service role is delegated to a project")
		}
		return nil, err
	}

	return &config.ServiceRolePermsResponse{
		ProjectID:   projectDelegated.ProjectID,
		ProjectName: projectDelegated.ProjectName,
		Env:         projectDelegated.Env,
	}, nil
}
