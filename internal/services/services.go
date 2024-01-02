package services

type Services struct {
	UserService *UserService
	TodoService *TodoService
}

type PaginatedResponse[T any] struct {
	Data       []T    `json:"data"`
	TotalCount int64  `json:"count"`
	Next       string `json:"next"`
	Prev       string `json:"prev"`
}

// todo - Detail response struct for each entity
