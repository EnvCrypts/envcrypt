package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (handler *Handler) ListServiceRoles(w http.ResponseWriter, r *http.Request) {
	var requestBody config.ServiceRoleListRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.ServiceRoles.List(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteResponse(w, http.StatusCreated, responseBody)
}

func (handler *Handler) GetServiceRole(w http.ResponseWriter, r *http.Request) {
	var requestBody config.ServiceRoleGetRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.ServiceRoles.Get(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteResponse(w, http.StatusCreated, responseBody)
}

func (handler *Handler) CreateServiceRole(w http.ResponseWriter, r *http.Request) {
	var requestBody config.ServiceRoleCreateRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Print(err)
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.ServiceRoles.Create(r.Context(), requestBody)
	log.Print(err)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteResponse(w, http.StatusCreated, *responseBody)
}

func (handler *Handler) DeleteServiceRole(w http.ResponseWriter, r *http.Request) {
	var requestBody config.ServiceRoleDeleteRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	err = handler.Services.ServiceRoles.Delete(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseBody := config.ServiceRoleDeleteResponse{
		Message: "Service role deleted!",
	}
	helpers.WriteResponse(w, http.StatusOK, responseBody)
}

func (handler *Handler) DelegateAccess(w http.ResponseWriter, r *http.Request) {
	var requestBody config.ServiceRoleDelegateRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	err = handler.Services.ServiceRoles.DelegateAccess(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseBody := config.ServiceRoleDelegateResponse{
		Message: "Service role delegated access!",
	}
	helpers.WriteResponse(w, http.StatusOK, responseBody)
}

func (handler *Handler) GetProjectKeys(w http.ResponseWriter, r *http.Request) {
	var requestBody config.ServiceRollProjectKeyRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.SessionService.GetProjectKeys(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteResponse(w, http.StatusOK, responseBody)
}

func (handler *Handler) GetPerms(w http.ResponseWriter, r *http.Request) {
	var requestBody config.ServiceRolePermsRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.ServiceRoles.GetPerms(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteResponse(w, http.StatusOK, responseBody)
}
