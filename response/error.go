package response

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

func (e ErrResponse) HasErrorField(errorField string) bool {
	_, ok := e.Errors[errorField]
	return ok
}

func (e ErrResponse) HasErrorFields(errorFields []string) bool {
	for _, f := range errorFields {
		if _, ok := e.Errors[f]; !ok {
			return false
		}
	}
	return true
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

func NewErrResponseException(message string) (int, ErrResponse) {
	return http.StatusInternalServerError, ErrResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}
