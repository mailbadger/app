package actions

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/routes/middleware"
	"github.com/FilipNikolovski/news-maily/storage"
	"github.com/FilipNikolovski/news-maily/utils/pagination"
	"github.com/gin-gonic/gin"
)

func GetLists(c *gin.Context) {
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

	storage.GetLists(c, middleware.GetUser(c).Id, p)
	c.JSON(http.StatusOK, p)
}

func GetList(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		if l, err := storage.GetList(c, id, middleware.GetUser(c).Id); err == nil {
			c.JSON(http.StatusOK, l)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"reason": "List not found",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Id must be an integer",
	})
	return
}

func PostList(c *gin.Context) {
	name := c.PostForm("name")

	l := &entities.List{
		Name:   name,
		UserId: middleware.GetUser(c).Id,
	}

	if !l.Validate() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": "Invalid data",
			"errors": l.Errors,
		})
		return
	}

	if err := storage.CreateList(c, l); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, l)
	return

}

func PutList(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		l, err := storage.GetList(c, id, middleware.GetUser(c).Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "List not found",
			})
			return
		}

		l.Name = c.PostForm("name")

		if !l.Validate() {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Invalid data",
				"errors": l.Errors,
			})
			return
		}

		if err = storage.UpdateList(c, l); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": err.Error(),
			})
			return
		}

		c.Status(http.StatusNoContent)
		return

	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Id must be an integer",
	})
	return
}

func DeleteList(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		user := middleware.GetUser(c)
		_, err := storage.GetList(c, id, user.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "List not found",
			})
			return
		}

		err = storage.DeleteList(c, id, user.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": err.Error(),
			})
			return
		}

		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Id must be an integer",
	})
	return
}
