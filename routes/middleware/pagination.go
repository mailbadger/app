package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/utils/pagination"
)

// PaginateWithCursor is a middleware that populates the cursor pagination object and sets it to the context.
// If the parameters are not valid the request is aborted.
func PaginateWithCursor() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := &pagination.Cursor{
			Path:       c.Request.URL.Path,
			PerPage:    pagination.DefaultPerPage,
			Collection: make([]interface{}, 0),
		}

		if len(c.Query("per_page")) > 0 {
			perpage, err := strconv.ParseInt(c.Query("per_page"), 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "per_page field must be an integer."})
				c.Abort()
				return
			}

			p.PerPage = perpage

			//Lock on 100 if the user requests more than 100 items per page
			if p.PerPage > 100 {
				p.PerPage = 100
			}
		}

		if len(c.Query("ending_before")) > 0 {
			endBefore, err := strconv.ParseInt(c.Query("ending_before"), 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "ending_before field must be an integer."})
				c.Abort()
				return
			}

			p.EndingBefore = endBefore

			c.Set("cursor", p)
			c.Next()
			return
		}

		if len(c.Query("starting_after")) > 0 {
			startAfter, err := strconv.ParseInt(c.Query("starting_after"), 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "starting_after field must be an integer."})
				c.Abort()
				return
			}

			p.StartingAfter = startAfter
		}

		c.Set("cursor", p)
		c.Next()
	}
}
