package middleware

import (
	"fmt"
	"math"
	"strconv"

	"github.com/FilipNikolovski/news-maily/utils/pagination"
	"github.com/gin-gonic/gin"
)

func Paginate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var p pagination.Pagination

		p.Page = 0
		p.PerPage = pagination.DefaultPerPage
		p.Total = math.MaxUint32
		p.Collection = make([]interface{}, 0)

		if len(c.Query("per_page")) > 0 {
			if len(c.Query("per_page")) > 1 {
				panic("More than one per_page parameter attached to get url")
			}
			perpage, err := strconv.ParseUint(c.Query("per_page"), 10, 32)
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
			if len(c.Query("page")) > 1 {
				panic("More than one page parameter attached to get url")
			}
			page, err := strconv.ParseUint(c.Query("page"), 10, 32)
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
