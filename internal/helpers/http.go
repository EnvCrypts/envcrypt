package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	errorBody := ErrorResponse{
		Error: errorMessage,
	}

	resp, err := json.Marshal(errorBody)
	if err != nil {
		_, err := w.Write([]byte(errorMessage))
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
		return
	}
}

func WriteResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	resp, err := json.Marshal(response)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
