package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (handler *Handler) GetEnv(w http.ResponseWriter, r *http.Request) {

	var RequestBody config.GetEnvRequest

	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {

		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	resp, err := handler.Services.Env.GetEnv(r.Context(), RequestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
}

func (handler *Handler) AddEnv(w http.ResponseWriter, r *http.Request) {
	var RequestBody config.AddEnvRequest

	err := json.NewDecoder(r.Body).Decode(&RequestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	err = handler.Services.Env.AddEnv(r.Context(), RequestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
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
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	err = handler.Services.Env.UpdateEnv(r.Context(), RequestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteResponse(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "env updated successfully",
	})
}
