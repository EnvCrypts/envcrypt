package handlers

import (
	"encoding/json"
	"errors"
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
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	user, err := handler.Services.Users.Create(r.Context(), requestBody)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	accessToken, refreshToken, err := handler.Services.SessionService.Refresh(r.Context(), user.Id)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	var response = config.CreateResponseBody{
		Message: "User created successfully",
		User:    *user,
		Session: config.SessionBody{
			AccessToken:  *accessToken,
			RefreshToken: *refreshToken,
			ExpiresIn:    600,
		},
	}

	helpers.WriteResponse(w, http.StatusCreated, response)
}

func (handler *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {

	var loginRequestBody config.LoginRequestBody
	err := json.NewDecoder(r.Body).Decode(&loginRequestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	user, err := handler.Services.Users.Login(r.Context(), loginRequestBody.Email, loginRequestBody.Password)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	accessToken, refreshToken, err := handler.Services.SessionService.Refresh(r.Context(), user.Id)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	var response = config.LoginResponseBody{
		Message: "Login successful",
		User:    *user,
		Session: config.SessionBody{
			AccessToken:  *accessToken,
			RefreshToken: *refreshToken,
			ExpiresIn:    600,
		},
	}

	helpers.WriteResponse(w, http.StatusOK, response)
}

func (handler *Handler) GetUserPublicKey(w http.ResponseWriter, r *http.Request) {

	var requestBody config.UserKeyRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	userId, publicKey, err := handler.Services.Users.GetUserPublicKey(r.Context(), requestBody.Email)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	var response = config.UserKeyResponseBody{
		PublicKey: publicKey,
		UserId:    userId,
		Message:   "User public key successfully retrieved",
	}

	helpers.WriteResponse(w, http.StatusOK, response)
}

func (handler *Handler) Refresh(w http.ResponseWriter, r *http.Request) {

	var requestBody config.RefreshRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	accessToken, refreshToken, err := handler.Services.SessionService.Refresh(r.Context(), requestBody.UserID)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	var response = config.RefreshResponseBody{
		Message: "Refresh successful",
		Session: config.SessionBody{
			AccessToken:  *accessToken,
			RefreshToken: *refreshToken,
			ExpiresIn:    600,
		},
	}
	helpers.WriteResponse(w, http.StatusOK, response)
}

func (handler *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var requestBody config.LogoutRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	err = handler.Services.Users.Logout(r.Context(), requestBody.UserID)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	var response = config.LogoutResponseBody{
		Message: "Logout successful",
	}
	helpers.WriteResponse(w, http.StatusOK, response)
}
