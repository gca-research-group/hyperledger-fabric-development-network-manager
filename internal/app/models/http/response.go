package http

import (
	"math"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/app/models/sql"
)

type Response[T any] struct {
	HasMore bool `json:"hasMore"`
	Total   int  `json:"total"`
	Page    int  `json:"page"`
	Pages   int  `json:"pages"`
	Data    T    `json:"data"`
}

func (r *Response[T]) NewResponse(data T, queryOptions sql.QueryOptions, total int) *Response[T] {
	r.HasMore = true

	if (queryOptions.Offset + queryOptions.Limit) >= total {
		r.HasMore = false
	}

	r.Total = int(total)
	r.Page = (queryOptions.Offset / queryOptions.Limit) + 1
	r.Pages = int(math.Ceil(float64(total) / float64(queryOptions.Limit)))
	r.Data = data

	return r
}
