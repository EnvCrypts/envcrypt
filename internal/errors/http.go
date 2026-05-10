package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type APIError struct {
	Code    Code              `json:"code"`
	Message string            `json:"message"`
	Hint    string            `json:"hint,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type APIResponse struct {
	Error APIError `json:"error"`
}

func ToHTTP(err error) (int, APIResponse) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		status := appErr.Status
		if status == 0 {
			status = http.StatusInternalServerError
		}
		return status, APIResponse{
			Error: APIError{
				Code:    appErr.Code,
				Message: appErr.Msg,
				Hint:    appErr.Hint,
				Fields:  appErr.Fields,
			},
		}
	}

	internal := Internal(err)
	return internal.Status, APIResponse{
		Error: APIError{Code: internal.Code, Message: internal.Msg},
	}
}

func Render(w http.ResponseWriter, err error, debug bool) {
	status, payload := ToHTTP(err)
	if debug {
		payload.Error.Message = payload.Error.Message + " (debug: " + err.Error() + ")"
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
