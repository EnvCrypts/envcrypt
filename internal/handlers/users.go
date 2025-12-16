package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (handler *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := handler.Services.Users.GetAllUsers(r.Context())
	if err != nil {
		return
	}

	helpers.WriteResponse(w, http.StatusOK, users)
	return
}

func (handler *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody config.CreateRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = handler.Services.Users.Create(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteResponse(w, http.StatusCreated, "User created successfully")
}
