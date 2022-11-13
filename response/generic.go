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
