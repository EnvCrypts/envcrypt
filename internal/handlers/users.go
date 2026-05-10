package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/errors"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (handler *Handler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := handler.Services.Users.GetAllUsers(r.Context())
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, users)
	return nil
}

func (handler *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.CreateRequestBody

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	user, err := handler.Services.Users.Create(r.Context(), requestBody)
	if err != nil {
		return err
	}

	accessToken, refreshToken, err := handler.Services.SessionService.Refresh(r.Context(), user.Id)
	if err != nil {
		return err
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
	return nil
}

func (handler *Handler) LoginUser(w http.ResponseWriter, r *http.Request) error {
	var loginRequestBody config.LoginRequestBody
	if err := json.NewDecoder(r.Body).Decode(&loginRequestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	user, err := handler.Services.Users.Login(r.Context(), loginRequestBody.Email, loginRequestBody.Password)
	if err != nil {
		return err
	}

	accessToken, refreshToken, err := handler.Services.SessionService.Refresh(r.Context(), user.Id)
	if err != nil {
		return err
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
	return nil
}

func (handler *Handler) GetUserPublicKey(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.UserKeyRequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	userId, publicKey, err := handler.Services.Users.GetUserPublicKey(r.Context(), requestBody.Email)
	if err != nil {
		return err
	}

	var response = config.UserKeyResponseBody{
		PublicKey: publicKey,
		UserId:    userId,
		Message:   "User public key successfully retrieved",
	}

	helpers.WriteResponse(w, http.StatusOK, response)
	return nil
}

func (handler *Handler) Refresh(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.RefreshRequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	accessToken, refreshToken, err := handler.Services.SessionService.Refresh(r.Context(), requestBody.UserID)
	if err != nil {
		return err
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
	return nil
}

func (handler *Handler) Logout(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.LogoutRequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	if err := handler.Services.Users.Logout(r.Context(), requestBody.UserID); err != nil {
		return err
	}

	var response = config.LogoutResponseBody{
		Message: "Logout successful",
	}
	helpers.WriteResponse(w, http.StatusOK, response)
	return nil
}

func (handler *Handler) RecoveryInit(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.RecoveryInitRequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	responseBody, err := handler.Services.Users.RecoveryInit(r.Context(), requestBody.Email)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, responseBody)
	return nil
}

func (handler *Handler) RecoveryComplete(w http.ResponseWriter, r *http.Request) error {
	var requestBody config.RecoveryCompleteRequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	if err := handler.Services.Users.RecoveryComplete(r.Context(), requestBody); err != nil {
		return err
	}

	response := config.RecoveryCompleteResponseBody{
		Message: "Recovery completed successfully",
	}
	helpers.WriteResponse(w, http.StatusOK, response)
	return nil
}
