package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (handler *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {

	var requestBody config.ProjectCreateRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	err = handler.Services.Projects.CreateProject(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	var responseBody = config.ProjectCreateResponse{
		Message: "Project created successfully!",
	}
	helpers.WriteResponse(w, http.StatusCreated, responseBody)
}

func (handler *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	var requestBody config.ListProjectRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	resp, err := handler.Services.Projects.ListProjects(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
}

func (handler *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {

	var requestBody config.ProjectDeleteRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	err = handler.Services.Projects.DeleteProject(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	var responseBody = config.ProjectDeleteResponse{
		Message: "Project deleted successfully!",
	}
	helpers.WriteResponse(w, http.StatusOK, responseBody)
}

func (handler *Handler) AddUserToProject(w http.ResponseWriter, r *http.Request) {

	var requestBody config.AddUserToProjectRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	err = handler.Services.Projects.AddUserToProject(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	var responseBody = config.ProjectCreateResponse{
		Message: "User added to project successfully!",
	}
	helpers.WriteResponse(w, http.StatusCreated, responseBody)
}

func (handler *Handler) SetUserAccess(w http.ResponseWriter, r *http.Request) {

	var requestBody config.SetAccessRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	err = handler.Services.Projects.SetUserAccess(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	var responseBody = config.ProjectCreateResponse{
		Message: "User access set successfully!",
	}
	helpers.WriteResponse(w, http.StatusOK, responseBody)
}

func (handler *Handler) GetUserProjectKeys(w http.ResponseWriter, r *http.Request) {

	var requestBody config.GetUserProjectRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	resp, err := handler.Services.Projects.GetUserProject(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
}

func (handler *Handler) GetMemberProject(w http.ResponseWriter, r *http.Request) {

	var requestBody config.GetMemberProjectRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	resp, err := handler.Services.Projects.GetMemberProject(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
}

func (handler *Handler) RotateInit(w http.ResponseWriter, r *http.Request) {
	var requestBody config.RotateInitRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	resp, err := handler.Services.Projects.RotateInit(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
}

func (handler *Handler) RotateCommit(w http.ResponseWriter, r *http.Request) {
	var requestBody config.RotateCommitRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	resp, err := handler.Services.Projects.RotateCommit(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
}
