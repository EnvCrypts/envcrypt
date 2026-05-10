package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Code string

const (
	CodeBadRequest   Code = "BAD_REQUEST"
	CodeValidation   Code = "VALIDATION_FAILED"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeForbidden    Code = "FORBIDDEN"
	CodeNotFound     Code = "NOT_FOUND"
	CodeConflict     Code = "CONFLICT"
	CodeInternal     Code = "INTERNAL_ERROR"
)

// AppError is a structured error for API responses and internal logging.
type AppError struct {
	Status int               `json:"-"`
	Code   Code              `json:"code"`
	Msg    string            `json:"message"`
	Hint   string            `json:"hint,omitempty"`
	Fields map[string]string `json:"fields,omitempty"`
	Err    error             `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err == nil {
		return e.Msg
	}
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code Code, status int, msg string) *AppError {
	return &AppError{Status: status, Code: code, Msg: msg}
}

func Wrap(code Code, status int, msg string, err error) *AppError {
	return &AppError{Status: status, Code: code, Msg: msg, Err: err}
}

func BadRequest(msg, hint string) *AppError {
	return &AppError{Status: http.StatusBadRequest, Code: CodeBadRequest, Msg: msg, Hint: hint}
}

func Validation(fields map[string]string) *AppError {
	return &AppError{Status: http.StatusBadRequest, Code: CodeValidation, Msg: "Validation failed", Fields: fields}
}

func Unauthorized(code Code, msg, hint string) *AppError {
	if code == "" {
		code = CodeUnauthorized
	}
	return &AppError{Status: http.StatusUnauthorized, Code: code, Msg: msg, Hint: hint}
}

func Forbidden(msg, hint string) *AppError {
	return &AppError{Status: http.StatusForbidden, Code: CodeForbidden, Msg: msg, Hint: hint}
}

func NotFound(resource, hint string) *AppError {
	code := strings.ToUpper(resource) + "_NOT_FOUND"
	return &AppError{Status: http.StatusNotFound, Code: Code(code), Msg: fmt.Sprintf("%s not found", resource), Hint: hint}
}

func Conflict(msg, hint string) *AppError {
	return &AppError{Status: http.StatusConflict, Code: CodeConflict, Msg: msg, Hint: hint}
}

func Internal(err error) *AppError {
	return Wrap(CodeInternal, http.StatusInternalServerError, "Internal server error", err)
}

func InternalMessage(msg string, err error) *AppError {
	return Wrap(CodeInternal, http.StatusInternalServerError, msg, err)
}

func IsCode(err error, code Code) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}
