package actions

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/logger"
	"github.com/news-maily/app/routes/middleware"
	"github.com/news-maily/app/storage"
	"github.com/sirupsen/logrus"
)

type subs struct {
	Ids []int64 `form:"ids[]"`
}

func GetSegments(c *gin.Context) {
	val, ok := c.Get("cursor")
	if !ok {
		logger.From(c).Error("Unable to fetch pagination cursor from context.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch segments. Please try again.",
		})
		return
	}

	p, ok := val.(*storage.PaginationCursor)
	if !ok {
		logger.From(c).Error("Unable to cast pagination cursor from context value.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch segments. Please try again.",
		})
		return
	}

	err := storage.GetSegments(c, middleware.GetUser(c).ID, p)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to fetch segments collection.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch segments. Please try again.",
		})
		return
	}

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
	name := strings.TrimSpace(c.PostForm("name"))

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
			"message": "Segment with that name already exists.",
		})
		return
	}

	if err := storage.CreateSegment(c, l); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to create segment.",
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

		l.Name = strings.TrimSpace(c.PostForm("name"))

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
				"message": "Unable to update segment.",
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
			logger.From(c).WithError(err).Error("Unable to delete segment.")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to delete segment.",
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
		err = c.Bind(subs)
		if err != nil {
			logger.From(c).WithError(err).Error("Unable to bind params")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Invalid parameters, please try again.",
			})
			return
		}

		if len(subs.Ids) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Ids list is empty",
			})
			return
		}

		s, err := storage.GetSubscribersByIDs(c, subs.Ids, user.ID)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{"ids": subs.Ids}).WithError(err).
				Warn("Unable to find subscribers by the list of ids.")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to add subscribers to the segment.",
			})
			return
		}

		l.Subscribers = s

		if err = storage.AppendSubscribers(c, l); err != nil {
			logger.From(c).WithFields(logrus.Fields{"ids": subs.Ids}).WithError(err).
				Error("Unable to create subscriber_segment associations by the list of ids.")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to add the subscribers to the segment.",
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
			logger.From(c).Error("Unable to fetch pagination cursor from context.")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch subscribers. Please try again.",
			})
			return
		}

		p, ok := val.(*storage.PaginationCursor)
		if !ok {
			logger.From(c).Error("Unable to cast pagination cursor from context value.")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch subscribers. Please try again.",
			})
			return
		}

		err := storage.GetSubscribersBySegmentID(c, id, middleware.GetUser(c).ID, p)
		if err != nil {
			logger.From(c).WithError(err).Error("Unable to fetch subscribers for segment collection.")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch segments. Please try again.",
			})
			return
		}

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
		err = c.Bind(subs)
		if err != nil {
			logger.From(c).WithError(err).Error("Unable to bind params")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Invalid parameters, please try again.",
			})
			return
		}

		if len(subs.Ids) == 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Ids list is empty",
			})
			return
		}

		s, err := storage.GetSubscribersByIDs(c, subs.Ids, user.ID)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{"ids": subs.Ids}).WithError(err).
				Error("Unable to find subscribers by the list of ids.")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to detach subscribers from the segment.",
			})
			return
		}

		l.Subscribers = s

		if err = storage.DetachSubscribers(c, l); err != nil {
			logger.From(c).WithFields(logrus.Fields{"ids": subs.Ids}).WithError(err).
				Error("Unable to remove subscriber_segment associations by the list of ids.")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to detach subscribers from the segment.",
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
