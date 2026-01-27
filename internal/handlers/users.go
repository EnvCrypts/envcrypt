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
	defer r.Body.Close()

	user, err := handler.Services.Users.Create(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response = config.CreateResponseBody{
		Message: "User created successfully",
		User:    *user,
	}

	helpers.WriteResponse(w, http.StatusCreated, response)
}

func (handler *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {

	var loginRequestBody config.LoginRequestBody
	err := json.NewDecoder(r.Body).Decode(&loginRequestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	user, err := handler.Services.Users.Login(r.Context(), loginRequestBody.Email, loginRequestBody.Password)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response = config.LoginResponseBody{
		Message: "Login successful",
		User:    *user,
	}

	helpers.WriteResponse(w, http.StatusOK, response)
}

func (handler *Handler) GetUserPublicKey(w http.ResponseWriter, r *http.Request) {

	var requestBody config.UserKeyRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	userId, publicKey, err := handler.Services.Users.GetUserPublicKey(r.Context(), requestBody.Email)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var response = config.UserKeyResponseBody{
		PublicKey: publicKey,
		UserId:    userId,
		Message:   "User public key successfully retrieved",
	}

	helpers.WriteResponse(w, http.StatusOK, response)
}
