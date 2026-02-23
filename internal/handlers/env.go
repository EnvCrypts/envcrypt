package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (handler *Handler) GetEnv(w http.ResponseWriter, r *http.Request) {

	var RequestBody config.GetEnvRequest

	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	resp, err := handler.Services.Env.GetEnv(r.Context(), RequestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
}

func (handler *Handler) GetEnvVersions(w http.ResponseWriter, r *http.Request) {

	var RequestBody config.GetEnvVersionsRequest

	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	resp, err := handler.Services.Env.GetEnvVersions(r.Context(), RequestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
}

func (handler *Handler) AddEnv(w http.ResponseWriter, r *http.Request) {
	var RequestBody config.AddEnvRequest

	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	err = handler.Services.Env.AddEnv(r.Context(), RequestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	helpers.WriteResponse(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "env added successfully",
	})
}

func (handler *Handler) UpdateEnv(w http.ResponseWriter, r *http.Request) {
	var RequestBody config.UpdateEnvRequest
	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	err = handler.Services.Env.UpdateEnv(r.Context(), RequestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	helpers.WriteResponse(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "env updated successfully",
	})
}

func (handler *Handler) GetCIEnv(w http.ResponseWriter, r *http.Request) {

	var requestBody config.GetEnvForCIRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.Env.GetEnvForCI(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	helpers.WriteResponse(w, http.StatusOK, responseBody)
}
