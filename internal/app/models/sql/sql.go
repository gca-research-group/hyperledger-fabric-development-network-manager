package sql

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type QueryOptions struct {
	Limit  int
	Offset int
}

func (q *QueryOptions) UpdateFromContext(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	page, _ := strconv.Atoi(c.Query("page"))

	if pageSize == 0 {
		pageSize = 20
	}

	if page == 0 {
		page = 1
	}

	offset := (page - 1) * pageSize

	q.Limit = pageSize
	q.Offset = offset
}
