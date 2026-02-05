package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
)

type ServiceRoleServices struct {
	q *database.Queries
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
			return nil, errors.New("no roles found")
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

	serviceRole, err := s.q.CreateServiceRole(ctx, database.CreateServiceRoleParams{
		Name:                 requestBody.ServiceRoleName,
		ServiceRolePublicKey: requestBody.ServiceRolePublicKey,
		RepoPrincipal:        requestBody.RepoPrincipal,
		CreatedBy:            requestBody.CreatedBy,
	})
	if err != nil {
		if dberrors.IsUniqueViolation(err) == true {
			return nil, errors.New("service role already exists")
		}
		return nil, err
	}

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

	serviceRole, err := s.q.GetServiceRoleById(ctx, requestBody.ServiceRoleId)
	if err != nil {
		return err
	}
	if requestBody.CreatedBy != serviceRole.CreatedBy {
		return errors.New("not authorized to delete service role")
	}

	_, err = s.q.DeleteServiceRole(ctx, serviceRole.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceRoleServices) DelegateAccess(ctx context.Context, requestBody config.ServiceRoleDelegateRequest) error {

	projectRole, err := s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		UserID:    requestBody.DelegatedBy,
		ProjectID: requestBody.ProjectId,
	})
	if err != nil {
		return errors.New("user role not found")
	}

	if projectRole.Role != "admin" {
		return errors.New("user is not an admin")
	}
	if projectRole.IsRevoked == true {
		return errors.New("user access is revoked")
	}

	serviceRole, err := s.q.GetServiceRoleByPrincipal(ctx, requestBody.RepoPrincipal)
	if err != nil {
		return err
	}

	_, err = s.q.DelegateAccess(ctx, database.DelegateAccessParams{
		ServiceRoleID:    serviceRole.ID,
		ProjectID:        requestBody.ProjectId,
		Env:              requestBody.EnvName,
		WrappedPmk:       requestBody.WrappedPMK,
		WrapNonce:        requestBody.WrapNonce,
		WrapEphemeralPub: requestBody.EphemeralPublicKey,
		DelegatedBy:      requestBody.DelegatedBy,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceRoleServices) GetPerms(ctx context.Context, requestBody config.ServiceRolePermsRequest) (*config.ServiceRolePermsResponse, error) {

	serviceRole, err := s.q.GetServiceRoleByPrincipal(ctx, requestBody.RepoPrincipal)
	if err != nil {
		return nil, err
	}

	projectDelegated, err := s.q.GetDelegation(ctx, serviceRole.ID)
	if err != nil {
		return nil, err
	}

	return &config.ServiceRolePermsResponse{
		ProjectID:   projectDelegated.ProjectID,
		ProjectName: projectDelegated.ProjectName,
		Env:         projectDelegated.Env,
	}, nil
}
