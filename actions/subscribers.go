package actions

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/routes/middleware"
	"github.com/news-maily/api/storage"
	"github.com/news-maily/api/utils/pagination"
)

func GetSubscribers(c *gin.Context) {
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

	storage.GetSubscribers(c, middleware.GetUser(c).Id, p)
	c.JSON(http.StatusOK, p)
}

func GetSubscriber(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		if s, err := storage.GetSubscriber(c, id, middleware.GetUser(c).Id); err == nil {
			c.JSON(http.StatusOK, s)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"reason": "Subscriber not found",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Id must be an integer",
	})
	return
}

func PostSubscriber(c *gin.Context) {
	name, email := c.PostForm("name"), c.PostForm("email")

	s := &entities.Subscriber{
		Name:   name,
		Email:  email,
		UserId: middleware.GetUser(c).Id,
	}

	if !s.Validate() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": "Invalid data",
			"errors": s.Errors,
		})
		return
	}

	_, err := storage.GetSubscriberByEmail(c, s.Email, s.UserId)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": "Subscriber with that email already exists",
		})
		return
	}

	if err := storage.CreateSubscriber(c, s); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, s)
	return

}

func PutSubscriber(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		s, err := storage.GetSubscriber(c, id, middleware.GetUser(c).Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Subscriber not found",
			})
			return
		}

		s.Name = c.PostForm("name")
		s.Email = c.PostForm("email")

		if !s.Validate() {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Invalid data",
				"errors": s.Errors,
			})
			return
		}

		if err = storage.UpdateSubscriber(c, s); err != nil {
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

func DeleteSubscriber(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		user := middleware.GetUser(c)
		_, err := storage.GetSubscriber(c, id, user.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Subscriber not found",
			})
			return
		}

		err = storage.DeleteSubscriber(c, id, user.Id)
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
