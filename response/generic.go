package response

import "net/http"

type Response[T any] struct {
	Data T `json:"data"`
}

func NewResponse[T any](model T) (int, Response[T]) {
	return http.StatusOK, Response[T]{Data: model}
}

func NewResponseCreated[T any](model T) (int, Response[T]) {
	return http.StatusCreated, Response[T]{Data: model}
}

type ManyResponsePaginated[T any] struct {
	Data       []T `json:"data"`
	Pagination struct {
		Total int64 `json:"total"`
	} `json:"pagination"`
}

func NewManyResponsePaginated[T any](models []T, total int64) (int, ManyResponsePaginated[T]) {
	resp := ManyResponsePaginated[T]{Data: models}
	resp.Pagination.Total = total
	return http.StatusOK, resp
}
