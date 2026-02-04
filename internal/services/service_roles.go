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
