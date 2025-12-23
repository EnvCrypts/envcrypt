package handlers

import (
	"encoding/json"
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

	helpers.WriteResponse(w, http.StatusCreated, "Project created successfully")
}
