package http

import (
	"math"
)

type Response[T any] struct {
	HasMore bool `json:"hasMore"`
	Total   int  `json:"total"`
	Page    int  `json:"page"`
	Pages   int  `json:"pages"`
	Data    T    `json:"data"`
}

func (r *Response[T]) NewResponse(data T, query Query, total int) *Response[T] {
	r.HasMore = true

	if (query.Offset + query.Limit) >= total {
		r.HasMore = false
	}

	r.Total = int(total)
	r.Page = (query.Offset / query.Limit) + 1
	r.Pages = int(math.Ceil(float64(total) / float64(query.Limit)))
	r.Data = data

	return r
}
