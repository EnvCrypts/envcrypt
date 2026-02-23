package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// AppError is a structured error that carries HTTP status, machine-readable code,
// human message, optional hint, and optional field-level validation errors.
type AppError struct {
	StatusCode int               `json:"-"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Hint       string            `json:"hint,omitempty"`
	Fields     map[string]string `json:"fields,omitempty"`
}

func (e *AppError) Error() string { return e.Message }


func ErrBadRequest(message, hint string) *AppError {
	return &AppError{StatusCode: http.StatusBadRequest, Code: "BAD_REQUEST", Message: message, Hint: hint}
}

func ErrValidation(fields map[string]string) *AppError {
	return &AppError{StatusCode: http.StatusBadRequest, Code: "VALIDATION_FAILED", Message: "Validation failed", Fields: fields}
}

func ErrUnauthorized(code, message, hint string) *AppError {
	return &AppError{StatusCode: http.StatusUnauthorized, Code: code, Message: message, Hint: hint}
}

func ErrForbidden(message, hint string) *AppError {
	return &AppError{StatusCode: http.StatusForbidden, Code: "PERMISSION_DENIED", Message: message, Hint: hint}
}

func ErrNotFound(resource, hint string) *AppError {
	code := strings.ToUpper(resource) + "_NOT_FOUND"
	return &AppError{StatusCode: http.StatusNotFound, Code: code, Message: fmt.Sprintf("%s not found", resource), Hint: hint}
}

func ErrConflict(message, hint string) *AppError {
	return &AppError{StatusCode: http.StatusConflict, Code: "CONFLICT", Message: message, Hint: hint}
}

func ErrInternal(message string) *AppError {
	return &AppError{StatusCode: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: message}
}


// errorEnvelope is the JSON envelope for error responses.
type errorEnvelope struct {
	Error *AppError `json:"error"`
}

// WriteError writes a structured JSON error response.
// If err is an *AppError its status/code/hint/fields are used.
// Otherwise the error is wrapped as a 500 INTERNAL_ERROR.
// The statusCode parameter is a fallback; when err is *AppError its own StatusCode takes precedence.
func WriteError(w http.ResponseWriter, statusCode int, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		appErr = &AppError{
			StatusCode: statusCode,
			Code:       "INTERNAL_ERROR",
			Message:    err.Error(),
		}
		if statusCode == http.StatusBadRequest {
			appErr.Code = "BAD_REQUEST"
		}
	}

	if appErr.StatusCode == 0 {
		appErr.StatusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(appErr.StatusCode)

	if encErr := json.NewEncoder(w).Encode(errorEnvelope{Error: appErr}); encErr != nil {
		log.Println("failed to write error response:", encErr)
	}
}

func WriteResponse(w http.ResponseWriter, statusCode int, response any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("failed to write response:", err)
	}
}
