package actions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/validator"

	"github.com/sirupsen/logrus"
)

func GetSegments(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("cursor")
		if !ok {
			logger.From(c).Error("get groups: unable to fetch pagination cursor from context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch segments. Please try again.",
			})
			return
		}

		p, ok := val.(*storage.PaginationCursor)
		if !ok {
			logger.From(c).Error("get groups: unable to cast pagination cursor from context value")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch segments. Please try again.",
			})
			return
		}

		err := store.GetSegments(middleware.GetUser(c).ID, p)
		if err != nil {
			logger.From(c).WithError(err).Error("get groups: unable to fetch segments collection")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch segments. Please try again.",
			})
			return
		}

		c.JSON(http.StatusOK, p)
	}
}

func GetSegment(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer.",
			})
		}

		userID := middleware.GetUser(c).ID

		s, err := storage.GetSegment(id, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Segment not found.",
			})
			return
		}

		totalSubs, err := storage.GetTotalSubscribers(userID)
		if err != nil {
			logger.From(c).WithError(err).Error("get group: unable to fetch total subscribers")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch group. Please try again.",
			})
			return
		}

		subsInSeg, err := storage.GetTotalSubscribersBySegment(s.ID, userID)
		if err != nil {
			logger.From(c).WithError(err).Error("get group: Unable to fetch total subscribers in segment")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch group. Please try again.",
			})
			return
		}

		c.JSON(http.StatusOK, &entities.SegmentWithTotalSubs{
			Segment:          *s,
			TotalSubscribers: &totalSubs,
			SubscribersInSeg: subsInSeg,
		})
	}
}

func PostSegment(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		body := &params.Segment{}
		if err := c.ShouldBindJSON(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid parameters, please try again.",
			})
			return
		}

		if err := validator.Validate(body); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		l := &entities.Segment{
			Name:   body.Name,
			UserID: middleware.GetUser(c).ID,
		}

		_, err := storage.GetSegmentByName(body.Name, middleware.GetUser(c).ID)
		if err == nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Segment with that name already exists.",
			})
			return
		}

		if err := storage.CreateSegment(l); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to create segment.",
			})
			return
		}

		c.JSON(http.StatusCreated, l)
	}
}

func PutSegment(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer.",
			})
		}

		l, err := storage.GetSegment(id, middleware.GetUser(c).ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Segment not found.",
			})
			return
		}

		body := &params.Segment{}
		if err := c.ShouldBindJSON(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid parameters, please try again.",
			})
			return
		}

		if err := validator.Validate(body); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		l2, err := storage.GetSegmentByName(body.Name, middleware.GetUser(c).ID)
		if err == nil && l2.ID != l.ID {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Segment with that name already exists.",
			})
			return
		}

		l.Name = body.Name

		if err = storage.UpdateSegment(l); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to update segment.",
			})
			return
		}

		c.JSON(http.StatusOK, l)
	}
}

func DeleteSegment(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer.",
			})
		}

		user := middleware.GetUser(c)
		err = storage.DeleteSegment(id, user.ID)
		if err != nil {
			logger.From(c).WithError(err).Error("unable to delete segment")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to delete segment.",
			})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func PutSegmentSubscribers(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer.",
			})
		}

		user := middleware.GetUser(c)
		l, err := storage.GetSegment(id, user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Segment not found.",
			})
			return
		}

		body := &params.SegmentSubs{}
		if err := c.ShouldBindJSON(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid parameters, please try again.",
			})
			return
		}

		if err := validator.Validate(body); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		s, err := storage.GetSubscribersByIDs(body.Ids, user.ID)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{"ids": body.Ids}).WithError(err).
				Error("put subs in group: unable to find subscribers by the list of ids")

			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to add subscribers to the segment.",
			})
			return
		}

		l.Subscribers = s

		err = storage.AppendSubscribers(l)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{"ids": body.Ids}).WithError(err).
				Error("put subs in group: unable to create subscriber_segment associations by the list of ids")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to add the subscribers to the segment.",
			})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func GetSegmentsubscribers(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer.",
			})
		}

		val, ok := c.Get("cursor")
		if !ok {
			logger.From(c).Error("get group subs: unable to fetch pagination cursor from context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch subscribers. Please try again.",
			})
			return
		}

		p, ok := val.(*storage.PaginationCursor)
		if !ok {
			logger.From(c).Error("get group subs: unable to cast pagination cursor from context value")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch subscribers. Please try again.",
			})
			return
		}

		err = store.GetSubscribersBySegmentID(id, middleware.GetUser(c).ID, p)
		if err != nil {
			logger.From(c).WithError(err).Error("get group subs: unable to fetch subscribers for segment collection")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch subscribers. Please try again.",
			})
			return
		}

		c.JSON(http.StatusOK, p)
	}
}

func DetachSegmentSubscribers(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer",
			})
		}

		user := middleware.GetUser(c)
		l, err := storage.GetSegment(id, user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Segment not found.",
			})
			return
		}

		body := &params.SegmentSubs{}
		if err := c.ShouldBindJSON(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid parameters, please try again.",
			})
			return
		}

		if err := validator.Validate(body); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		s, err := storage.GetSubscribersByIDs(body.Ids, user.ID)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{"ids": body.Ids}).WithError(err).
				Error("detach subs: unable to find subscribers by the list of ids")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to detach subscribers from the segment.",
			})
			return
		}

		l.Subscribers = s

		if err = storage.DetachSubscribers(l); err != nil {
			logger.From(c).WithFields(logrus.Fields{"ids": body.Ids}).WithError(err).
				Error("detach subs: unable to remove subscriber_segment associations by the list of ids")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to detach subscribers from the segment.",
			})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func DetachSubscriber(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer.",
			})
		}

		subID, err := strconv.ParseInt(c.Param("sub_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Subscriber id must be an integer.",
			})
		}

		user := middleware.GetUser(c)
		l, err := storage.GetSegment(id, user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Segment not found.",
			})
			return
		}

		s, err := storage.GetSubscriber(subID, user.ID)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{"subscriber_id": subID, "segment_id": id}).WithError(err).
				Error("detach sub: unable to find subscriber by id")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to remove subscriber from the segment, the subscriber does not exist.",
			})
			return
		}

		l.Subscribers = []entities.Subscriber{*s}

		if err = storage.DetachSubscribers(l); err != nil {
			logger.From(c).WithFields(logrus.Fields{"subscriber_id": subID, "segment_id": id}).WithError(err).
				Error("detach sub: unable to remove subscriber from segment")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to remove subscriber from the segment.",
			})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
