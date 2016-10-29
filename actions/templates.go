package actions

import (
	"errors"
	"net/http"

	"github.com/FilipNikolovski/news-maily/routes/middleware"
	"github.com/FilipNikolovski/news-maily/storage"
	"github.com/FilipNikolovski/news-maily/utils/pagination"
	"github.com/gin-gonic/gin"
)

func GetTemplates(c *gin.Context) {
	val, ok := c.Get("pagination")
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("cannot create pagination object"))
		return
	}

	p, ok := val.(*pagination.Pagination)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("cannot cast pagination object"))
		return
	}

	storage.GetTemplates(c, middleware.GetUser(c).Id, p)
	c.JSON(http.StatusOK, p)
}
