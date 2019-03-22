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
	"github.com/sirupsen/logrus"
)

type subs struct {
	Ids []int64 `form:"ids[]"`
}

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
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
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

	_, err := storage.GetListByName(c, name, middleware.GetUser(c).Id)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": "List with that name already exists",
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
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
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

		l2, err := storage.GetListByName(c, l.Name, middleware.GetUser(c).Id)
		if err == nil && l2.Id != l.Id {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "List with that name already exists",
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
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
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

func PutListSubscribers(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		user := middleware.GetUser(c)
		l, err := storage.GetList(c, id, user.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "List not found",
			})
			return
		}

		subs := &subs{}
		c.Bind(subs)

		if len(subs.Ids) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Ids list is empty",
			})
			return
		}

		s, err := storage.GetSubscribersByIDs(c, subs.Ids, user.Id)
		if err != nil {
			logrus.Warn(err)
		}

		l.Subscribers = s

		if err = storage.AppendSubscribers(c, l); err != nil {
			logrus.Error(err)
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

func GetListSubscribers(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
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

		storage.GetSubscribersByListID(c, id, middleware.GetUser(c).Id, p)
		c.JSON(http.StatusOK, p)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Id must be an integer",
	})
	return
}

func DetachListSubscribers(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		user := middleware.GetUser(c)
		l, err := storage.GetList(c, id, user.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "List not found",
			})
			return
		}

		subs := &subs{}
		c.Bind(subs)

		if len(subs.Ids) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Ids list is empty",
			})
			return
		}

		s, err := storage.GetSubscribersByIDs(c, subs.Ids, user.Id)
		if err != nil {
			logrus.Warn(err)
		}

		l.Subscribers = s

		if err = storage.DetachSubscribers(c, l); err != nil {
			logrus.Error(err)
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
