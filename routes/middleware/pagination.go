package middleware

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/utils/pagination"
)

func Paginate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var p = new(pagination.Pagination)

		p.Page = 0
		p.PerPage = pagination.DefaultPerPage
		p.Total = math.MaxUint64
		p.Collection = make([]interface{}, 0)

		if len(c.Query("per_page")) > 0 {
			perpage, err := strconv.ParseUint(c.Query("per_page"), 10, 64)
			if err != nil {
				panic(fmt.Sprintf("Error parsing 'per_page': %s", err))
			}

			p.PerPage = uint(perpage)

			//Lock on 100 if the user requests more than 100 items per page
			if p.PerPage > 100 {
				p.PerPage = 100
			}
		}
		if len(c.Query("page")) > 0 {
			page, err := strconv.ParseUint(c.Query("page"), 10, 64)
			if err != nil {
				panic(fmt.Sprintf("Error parsing 'page': %s", err))
			}
			p.Page = uint(page)
			p.Offset = uint(page * uint64(p.PerPage))
		}

		c.Set("pagination", p)
		c.Next()
	}
}
