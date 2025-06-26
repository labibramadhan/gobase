package modeldto

type PageableDto struct {
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
	Page   *int `json:"page"`
}

type ListPageableDto[T any] struct {
	Data       []T         `json:"data"`
	Pagination PageableDto `json:"pagination"`
}
