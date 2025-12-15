package handlers

import (
	"encoding/json"
	"net/http"
)

func (handler *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := handler.Services.Users.GetByEmail(r.Context(), "null")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonData, _ := json.Marshal(users)
	w.Write(jsonData)
	return
}
