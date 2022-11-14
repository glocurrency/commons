package response

type SingleResponse[T any] struct {
	Data T `json:"data"`
}

func NewSingleResponse[T any](model T) SingleResponse[T] {
	return SingleResponse[T]{Data: model}
}

type ManyResponse[T any] struct {
	Data []T `json:"data"`
}

func NewManyResponse[T any](models []T) ManyResponse[T] {
	return ManyResponse[T]{Data: models}
}

type ManyResponsePaginated[T any] struct {
	Data       []T `json:"data"`
	Pagination struct {
		Total int64 `json:"total"`
	} `json:"pagination"`
}

func NewManyResponsePaginated[T any](models []T, total int64) ManyResponsePaginated[T] {
	resp := ManyResponsePaginated[T]{Data: models}
	resp.Pagination.Total = total
	return resp
}
