package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (handler *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {

	var requestBody config.ProjectCreateRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	err = handler.Services.Projects.CreateProject(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var responseBody = config.ProjectCreateResponse{
		Message: "Project created successfully!",
	}
	helpers.WriteResponse(w, http.StatusCreated, responseBody)
}

func (handler *Handler) AddUserToProject(w http.ResponseWriter, r *http.Request) {

	var requestBody config.AddUserToProjectRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	err = handler.Services.Projects.AddUserToProject(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var responseBody = config.ProjectCreateResponse{
		Message: "User added to project successfully!",
	}
	helpers.WriteResponse(w, http.StatusCreated, responseBody)
}

func (handler *Handler) GetUserProjectKeys(w http.ResponseWriter, r *http.Request) {

	var requestBody config.GetUserProjectRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	resp, err := handler.Services.Projects.GetUserProject(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		log.Print(err.Error())
		return
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
}
