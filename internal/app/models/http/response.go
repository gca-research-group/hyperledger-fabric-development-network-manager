package http

type Response[T any] struct {
	HasMore bool `json:"hasMore"`
	Total   int  `json:"total"`
	Page    int  `json:"page"`
	Pages   int  `json:"pages"`
	Data    T    `json:"data"`
}
