package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/errors"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (handler *Handler) GetEnv(w http.ResponseWriter, r *http.Request) error {

	var RequestBody config.GetEnvRequest

	if err := json.NewDecoder(r.Body).Decode(&RequestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	resp, err := handler.Services.Env.GetEnv(r.Context(), RequestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
	return nil
}

func (handler *Handler) GetEnvVersions(w http.ResponseWriter, r *http.Request) error {

	var RequestBody config.GetEnvVersionsRequest

	if err := json.NewDecoder(r.Body).Decode(&RequestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	resp, err := handler.Services.Env.GetEnvVersions(r.Context(), RequestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
	return nil
}

func (handler *Handler) AddEnv(w http.ResponseWriter, r *http.Request) error {
	var RequestBody config.AddEnvRequest

	if err := json.NewDecoder(r.Body).Decode(&RequestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	if err := handler.Services.Env.AddEnv(r.Context(), RequestBody); err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "env added successfully",
	})
	return nil
}

func (handler *Handler) UpdateEnv(w http.ResponseWriter, r *http.Request) error {
	var RequestBody config.UpdateEnvRequest
	if err := json.NewDecoder(r.Body).Decode(&RequestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	if err := handler.Services.Env.UpdateEnv(r.Context(), RequestBody); err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "env updated successfully",
	})
	return nil
}

func (handler *Handler) GetCIEnv(w http.ResponseWriter, r *http.Request) error {

	var requestBody config.GetEnvForCIRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.Env.GetEnvForCI(r.Context(), requestBody)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, responseBody)
	return nil
}
