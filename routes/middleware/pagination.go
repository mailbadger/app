package middleware

import (
	"fmt"
	"math"
	"strconv"

	"github.com/FilipNikolovski/news-maily/utils"
	"github.com/gin-gonic/gin"
)

func Paginate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pagination utils.Pagination

		pagination.Page = 0
		pagination.PerPage = utils.DefaultPerPage
		pagination.Total = math.MaxUint32
		pagination.Collection = make([]interface{}, 0)

		if len(c.Query("per_page")) > 0 {
			if len(c.Query("per_page")) > 1 {
				panic("More than one per_page parameter attached to get url")
			}
			perpage, err := strconv.ParseUint(c.Query("per_page"), 10, 32)
			if err != nil {
				panic(fmt.Sprintf("Error parsing 'per_page': %s", err))
			}

			pagination.PerPage = uint(perpage)

			//Lock on 100 if the user requests more than 100 items per page
			if pagination.PerPage > 100 {
				pagination.PerPage = 100
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
			pagination.Page = uint(page)
			pagination.Offset = uint(page * uint64(pagination.PerPage))
		}

		c.Set("pagination", pagination)
		c.Next()
	}
}
