package response

import "net/http"

type SuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewSuccessResponse(message string) (int, SuccessResponse) {
	return http.StatusOK, SuccessResponse{Code: http.StatusOK, Message: message}
}

func NewSuccessResponseCreated(message string) (int, SuccessResponse) {
	return http.StatusCreated, SuccessResponse{Code: http.StatusCreated, Message: message}
}
