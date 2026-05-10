package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/errors"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (handler *Handler) ListServiceRoles(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.ServiceRoleListRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.ServiceRoles.List(r.Context(), requestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, responseBody)
	return nil
}

func (handler *Handler) GetServiceRole(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.ServiceRoleGetRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.ServiceRoles.Get(r.Context(), requestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, responseBody)
	return nil
}

func (handler *Handler) CreateServiceRole(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.ServiceRoleCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.ServiceRoles.Create(r.Context(), requestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusCreated, *responseBody)
	return nil
}

func (handler *Handler) DeleteServiceRole(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.ServiceRoleDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	if err := handler.Services.ServiceRoles.Delete(r.Context(), requestBody); err != nil {
		return err
	}

	responseBody := config.ServiceRoleDeleteResponse{
		Message: "Service role deleted!",
	}
	helpers.WriteResponse(w, http.StatusOK, responseBody)
	return nil
}

func (handler *Handler) DelegateAccess(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.ServiceRoleDelegateRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	if err := handler.Services.ServiceRoles.DelegateAccess(r.Context(), requestBody); err != nil {
		return err
	}

	responseBody := config.ServiceRoleDelegateResponse{
		Message: "Service role delegated access!",
	}
	helpers.WriteResponse(w, http.StatusOK, responseBody)
	return nil
}

func (handler *Handler) GetProjectKeys(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.ServiceRollProjectKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.SessionService.GetProjectKeys(r.Context(), requestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, responseBody)
	return nil
}

func (handler *Handler) GetPerms(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.ServiceRolePermsRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.ServiceRoles.GetPerms(r.Context(), requestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, responseBody)
	return nil
}
