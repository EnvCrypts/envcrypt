package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/errors"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (handler *Handler) CreateProject(w http.ResponseWriter, r *http.Request) error {

	var requestBody config.ProjectCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	if err := handler.Services.Projects.CreateProject(r.Context(), requestBody); err != nil {
		return err
	}

	var responseBody = config.ProjectCreateResponse{
		Message: "Project created successfully!",
	}
	helpers.WriteResponse(w, http.StatusCreated, responseBody)
	return nil
}

func (handler *Handler) ListProjects(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.ListProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	resp, err := handler.Services.Projects.ListProjects(r.Context(), requestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
	return nil
}

func (handler *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) error {

	var requestBody config.ProjectDeleteRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	if err := handler.Services.Projects.DeleteProject(r.Context(), requestBody); err != nil {
		return err
	}

	var responseBody = config.ProjectDeleteResponse{
		Message: "Project deleted successfully!",
	}
	helpers.WriteResponse(w, http.StatusOK, responseBody)
	return nil
}

func (handler *Handler) AddUserToProject(w http.ResponseWriter, r *http.Request) error {

	var requestBody config.AddUserToProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	if err := handler.Services.Projects.AddUserToProject(r.Context(), requestBody); err != nil {
		return err
	}

	var responseBody = config.ProjectCreateResponse{
		Message: "User added to project successfully!",
	}
	helpers.WriteResponse(w, http.StatusCreated, responseBody)
	return nil
}

func (handler *Handler) SetUserAccess(w http.ResponseWriter, r *http.Request) error {

	var requestBody config.SetAccessRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	if err := handler.Services.Projects.SetUserAccess(r.Context(), requestBody); err != nil {
		return err
	}

	var responseBody = config.ProjectCreateResponse{
		Message: "User access set successfully!",
	}
	helpers.WriteResponse(w, http.StatusOK, responseBody)
	return nil
}

func (handler *Handler) GetUserProjectKeys(w http.ResponseWriter, r *http.Request) error {

	var requestBody config.GetUserProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	resp, err := handler.Services.Projects.GetUserProject(r.Context(), requestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
	return nil
}

func (handler *Handler) GetMemberProject(w http.ResponseWriter, r *http.Request) error {

	var requestBody config.GetMemberProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	resp, err := handler.Services.Projects.GetMemberProject(r.Context(), requestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
	return nil
}

func (handler *Handler) RotateInit(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.RotateInitRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	resp, err := handler.Services.Projects.RotateInit(r.Context(), requestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
	return nil
}

func (handler *Handler) RotateCommit(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.RotateCommitRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	resp, err := handler.Services.Projects.RotateCommit(r.Context(), requestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
	return nil
}

func (handler *Handler) HandleProjectAuditLogs(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.ProjectAuditRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	if requestBody.ProjectID == uuid.Nil {
		return errors.Validation(map[string]string{"project_id": "project_id is required"})
	}

	resp, err := handler.Services.Audit.GetProjectAuditLogs(r.Context(), requestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
	return nil
}
