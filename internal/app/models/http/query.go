package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Query struct {
	Limit  int
	Offset int
}

func (q *Query) UpdateFromContext(c *gin.Context) {
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
