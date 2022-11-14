package response

import "net/http"

type SingleResponse[T any] struct {
	Data T `json:"data"`
}

func NewSingleResponse[T any](model T) (int, SingleResponse[T]) {
	return http.StatusOK, SingleResponse[T]{Data: model}
}

type ManyResponse[T any] struct {
	Data []T `json:"data"`
}

func NewManyResponse[T any](models []T) (int, ManyResponse[T]) {
	return http.StatusOK, ManyResponse[T]{Data: models}
}

func NewManyResponseCreated[T any](models []T) (int, ManyResponse[T]) {
	return http.StatusCreated, ManyResponse[T]{Data: models}
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
