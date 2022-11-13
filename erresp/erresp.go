package erresp

import "net/http"

// ErrResponse is a generic error response
type ErrResponse struct {
	// Code is the error code
	Code int `json:"code"`
	// Message is the error message
	Message string `json:"message"`
	// Validation errors
	Errors map[string]string `json:"errors,omitempty"`
}

func NewErrResponse(e error) (int, ErrResponse) {
	return http.StatusInternalServerError, ErrResponse{
		Code:    http.StatusInternalServerError,
		Message: e.Error(),
	}
}

func NewErrResponseForbidden(message string) (int, ErrResponse) {
	return http.StatusForbidden, ErrResponse{
		Code:    http.StatusForbidden,
		Message: message,
	}
}

func NewErrResponseNotFound(message string) (int, ErrResponse) {
	return http.StatusNotFound, ErrResponse{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewErrResponseBadRequest(message string) (int, ErrResponse) {
	return http.StatusBadRequest, ErrResponse{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewErrResponseValidationErrors(message string, errors map[string]string) (int, ErrResponse) {
	return http.StatusBadRequest, ErrResponse{
		Code:    http.StatusBadRequest,
		Message: message,
		Errors:  errors,
	}
}
