package actions

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/routes/middleware"
	"github.com/news-maily/app/storage"
	"github.com/news-maily/app/utils/pagination"
	"github.com/sirupsen/logrus"
)

type subs struct {
	Ids []int64 `form:"ids[]"`
}

func GetSegments(c *gin.Context) {
	val, ok := c.Get("cursor")
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("cannot create pagination object"))
		return
	}

	p, ok := val.(*pagination.Cursor)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("cannot cast pagination object"))
		return
	}

	storage.GetSegments(c, middleware.GetUser(c).ID, p)
	c.JSON(http.StatusOK, p)
}

func GetSegment(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		if l, err := storage.GetSegment(c, id, middleware.GetUser(c).ID); err == nil {
			c.JSON(http.StatusOK, l)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Segment not found",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Id must be an integer",
	})
}

func PostSegment(c *gin.Context) {
	name := c.PostForm("name")

	l := &entities.Segment{
		Name:   name,
		UserID: middleware.GetUser(c).ID,
	}

	if !l.Validate() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid data",
			"errors":  l.Errors,
		})
		return
	}

	_, err := storage.GetSegmentByName(c, name, middleware.GetUser(c).ID)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Segment with that name already exists",
		})
		return
	}

	if err := storage.CreateSegment(c, l); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, l)
}

func PutSegment(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		l, err := storage.GetSegment(c, id, middleware.GetUser(c).ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Segment not found",
			})
			return
		}

		l.Name = c.PostForm("name")

		if !l.Validate() {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Invalid data",
				"errors":  l.Errors,
			})
			return
		}

		l2, err := storage.GetSegmentByName(c, l.Name, middleware.GetUser(c).ID)
		if err == nil && l2.ID != l.ID {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Segment with that name already exists",
			})
			return
		}

		if err = storage.UpdateSegment(c, l); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.Status(http.StatusNoContent)
		return

	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Id must be an integer",
	})
}

func DeleteSegment(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		user := middleware.GetUser(c)
		_, err := storage.GetSegment(c, id, user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Segment not found",
			})
			return
		}

		err = storage.DeleteSegment(c, id, user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Id must be an integer",
	})
}

func PutSegmentSubscribers(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		user := middleware.GetUser(c)
		l, err := storage.GetSegment(c, id, user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Segment not found",
			})
			return
		}

		subs := &subs{}
		c.Bind(subs)

		if len(subs.Ids) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Ids list is empty",
			})
			return
		}

		s, err := storage.GetSubscribersByIDs(c, subs.Ids, user.ID)
		if err != nil {
			logrus.Warn(err)
		}

		l.Subscribers = s

		if err = storage.AppendSubscribers(c, l); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Id must be an integer",
	})
}

func GetSegmentsubscribers(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		val, ok := c.Get("cursor")
		if !ok {
			c.AbortWithError(http.StatusInternalServerError, errors.New("cannot create pagination object"))
			return
		}

		p, ok := val.(*pagination.Cursor)
		if !ok {
			c.AbortWithError(http.StatusInternalServerError, errors.New("cannot cast pagination object"))
			return
		}

		storage.GetSubscribersBySegmentID(c, id, middleware.GetUser(c).ID, p)
		c.JSON(http.StatusOK, p)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Id must be an integer",
	})
}

func DetachSegmentSubscribers(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		user := middleware.GetUser(c)
		l, err := storage.GetSegment(c, id, user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Segment not found",
			})
			return
		}

		subs := &subs{}
		c.Bind(subs)

		if len(subs.Ids) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Ids list is empty",
			})
			return
		}

		s, err := storage.GetSubscribersByIDs(c, subs.Ids, user.ID)
		if err != nil {
			logrus.Warn(err)
		}

		l.Subscribers = s

		if err = storage.DetachSubscribers(c, l); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Id must be an integer",
	})
}
