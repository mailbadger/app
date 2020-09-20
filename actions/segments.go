package actions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/storage"
)

type subs struct {
	Ids []int64 `form:"ids[]" binding:"required"`
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
	}

	userID := middleware.GetUser(c).ID

	s, err := storage.GetSegment(c, id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Segment not found.",
		})
		return
	}

	totalSubs, err := storage.GetTotalSubscribers(c, userID)
	if err != nil {
		logger.From(c).WithError(err).Warn("Unable to fetch total subscribers.")
	}

	subsInSeg, err := storage.GetTotalSubscribersBySegment(c, s.ID, userID)
	if err != nil {
		logger.From(c).WithError(err).Warn("Unable to fetch total subscribers in segment.")
	}

	c.JSON(http.StatusOK, &entities.SegmentWithTotalSubs{
		Segment:          *s,
		TotalSubscribers: &totalSubs,
		SubscribersInSeg: subsInSeg,
	})
}

type segmentsParams struct {
	Name string `form:"name" binding:"required,max=191"`
}

func PostSegment(c *gin.Context) {
	params := &segmentsParams{}
	if err := c.ShouldBind(params); err != nil {
		AbortWithError(c, err)
		return
	}

	l := &entities.Segment{
		Name:   params.Name,
		UserID: middleware.GetUser(c).ID,
	}

	_, err := storage.GetSegmentByName(c, params.Name, middleware.GetUser(c).ID)
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
		params := &segmentsParams{}
		if err := c.ShouldBind(params); err != nil {
			AbortWithError(c, err)
			return
		}
		l.Name = params.Name

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
		if err := c.ShouldBind(subs); err != nil {
			AbortWithError(c, err)
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
		if err := c.ShouldBind(subs); err != nil {
			AbortWithError(c, err)
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

func DetachSubscriber(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
	}

	subID, err := strconv.ParseInt(c.Param("sub_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Subscriber id must be an integer",
		})
	}

	user := middleware.GetUser(c)
	l, err := storage.GetSegment(c, id, user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Segment not found",
		})
		return
	}

	s, err := storage.GetSubscriber(c, subID, user.ID)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{"subscriber_id": subID, "segment_id": id}).WithError(err).
			Warn("Unable to find subscriber by id.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to remove subscriber from the segment, the subscriber does not exist.",
		})
		return
	}

	l.Subscribers = []entities.Subscriber{*s}

	if err = storage.DetachSubscribers(c, l); err != nil {
		logger.From(c).WithFields(logrus.Fields{"subscriber_id": subID, "segment_id": id}).WithError(err).
			Warn("Unable to remove subscriber from segment.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to remove subscriber from the segment.",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
