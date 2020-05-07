package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/storage"
)

// PaginateWithCursor is a middleware that populates the cursor pagination object and sets it to the context.
// If the parameters are not valid the request is aborted.
func PaginateWithCursor() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := storage.NewPaginationCursor(c.Request.URL.Path, storage.DefaultPerPage)

		if len(c.Query("per_page")) > 0 {
			perpage, err := strconv.ParseInt(c.Query("per_page"), 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "per_page field must be an integer."})
				return
			}

			p.SetPerPage(perpage)
		}

		if len(c.Query("ending_before")) > 0 {
			endBefore, err := strconv.ParseInt(c.Query("ending_before"), 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ending_before field must be an integer."})
				return
			}

			p.SetEndingBefore(endBefore)

			c.Set("cursor", p)
			c.Next()
			return
		}

		if len(c.Query("starting_after")) > 0 {
			startAfter, err := strconv.ParseInt(c.Query("starting_after"), 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "starting_after field must be an integer."})
				return
			}

			p.SetStartingAfter(startAfter)
		}

		c.Set("cursor", p)
		c.Next()
	}
}
