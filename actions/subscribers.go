package actions

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/mailbadger/app/services/subscribers/bulkremover"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/routes/middleware"
	awss3 "github.com/mailbadger/app/s3"
	"github.com/mailbadger/app/services/subscribers/importer"
	"github.com/mailbadger/app/storage"
	"github.com/sirupsen/logrus"
)

func GetSubscribers(c *gin.Context) {
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

	err := storage.GetSubscribers(c, middleware.GetUser(c).ID, p)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to fetch subscribers collection.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch subscribers. Please try again.",
		})
		return
	}

	c.JSON(http.StatusOK, p)
}

func GetSubscriber(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		if s, err := storage.GetSubscriber(c, id, middleware.GetUser(c).ID); err == nil {
			c.JSON(http.StatusOK, s)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Subscriber not found",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Id must be an integer",
	})
}

type segmentsParam struct {
	Ids []int64 `form:"segments[]"`
}

func PostSubscriber(c *gin.Context) {
	name, email := strings.TrimSpace(c.PostForm("name")), strings.TrimSpace(c.PostForm("email"))
	meta := c.PostFormMap("metadata")
	segments := &segmentsParam{}
	err := c.Bind(segments)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid data",
			"errors": map[string]string{
				"segments": "The segments array is in an invalid format.",
			},
		})
		return
	}

	s := &entities.Subscriber{
		Name:     name,
		Email:    email,
		Metadata: meta,
		Active:   true,
		UserID:   middleware.GetUser(c).ID,
	}

	segs, err := storage.GetSegmentsByIDs(c, s.UserID, segments.Ids)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid data",
			"errors": map[string]string{
				"segments": "Unable to find the specified segments.",
			},
		})
		return
	}

	s.Segments = segs

	if !s.Validate() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid data",
			"errors":  s.Errors,
		})
		return
	}

	_, err = storage.GetSubscriberByEmail(c, s.Email, s.UserID)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Subscriber with that email already exists.",
		})
		return
	}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to create subscriber, invalid metadata.",
		})
		return
	}
	s.MetaJSON = metaJSON

	if err := storage.CreateSubscriber(c, s); err != nil {
		logger.From(c).
			WithError(err).
			Warn("Unable to create subscriber.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to create subscriber",
		})
		return
	}

	c.JSON(http.StatusCreated, s)
}

func PutSubscriber(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
		return
	}

	s, err := storage.GetSubscriber(c, id, middleware.GetUser(c).ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Subscriber not found",
		})
		return
	}

	s.Name = strings.TrimSpace(c.PostForm("name"))
	s.Metadata = c.PostFormMap("metadata")

	segments := &segmentsParam{}
	err = c.Bind(segments)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid data",
			"errors": map[string]string{
				"segments": "The segments array is in an invalid format.",
			},
		})
		return
	}

	segs, err := storage.GetSegmentsByIDs(c, s.UserID, segments.Ids)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid data",
			"errors": map[string]string{
				"segments": "Unable to find the specified segments.",
			},
		})
		return
	}

	s.Segments = segs

	if !s.Validate() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid data",
			"errors":  s.Errors,
		})
		return
	}

	metaJSON, err := json.Marshal(s.Metadata)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to create subscriber, invalid metadata.",
		})
		return
	}
	s.MetaJSON = metaJSON

	if err = storage.UpdateSubscriber(c, s); err != nil {
		logger.From(c).
			WithError(err).
			WithField("subscriber_id", id).
			Warn("Unable to update subscriber.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to update subscriber.",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteSubscriber(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		user := middleware.GetUser(c)
		_, err := storage.GetSubscriber(c, id, user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Subscriber not found",
			})
			return
		}

		err = storage.DeleteSubscriber(c, id, user.ID)
		if err != nil {
			logger.From(c).WithError(err).Warn("Unable to delete subscriber.")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to delete subscriber.",
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

func PostUnsubscribe(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	uuid := strings.TrimSpace(c.PostForm("uuid"))
	token := strings.TrimSpace(c.PostForm("t"))

	redirWithError := c.Request.Referer()

	params := url.Values{}
	params.Add("email", email)
	params.Add("uuid", uuid)
	params.Add("t", token)
	params.Add("failed", "true")

	redirWithError = redirWithError + "?" + params.Encode()

	if token == "" || email == "" || uuid == "" {
		c.Redirect(http.StatusPermanentRedirect, redirWithError)
		return
	}

	u, err := storage.GetUserByUUID(c, uuid)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"email": email,
			"uuid":  uuid,
		}).WithError(err).Warn("Unsubscribe: cannot find user by uuid.")
		c.Redirect(http.StatusPermanentRedirect, redirWithError)
		return
	}

	sub, err := storage.GetSubscriberByEmail(c, email, u.ID)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"email": email,
			"uuid":  uuid,
		}).WithError(err).Warn("Unsubscribe: unable to fetch subscriber by email.")
		c.Redirect(http.StatusPermanentRedirect, redirWithError)
		return
	}

	hash, err := sub.GenerateUnsubscribeToken(os.Getenv("UNSUBSCRIBE_SECRET"))
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"email": email,
			"uuid":  uuid,
		}).WithError(err).Error("Unsubscribe: unable to generate hash.")
		c.Redirect(http.StatusPermanentRedirect, redirWithError)
		return
	}

	if subtle.ConstantTimeCompare([]byte(token), []byte(hash)) != 1 {
		logger.From(c).WithFields(logrus.Fields{
			"email": email,
			"uuid":  uuid,
		}).WithError(err).Warn("Unsubscribe: hashes don't match.")
		c.Redirect(http.StatusPermanentRedirect, redirWithError)
		return
	}

	err = storage.DeactivateSubscriber(c, u.ID, sub.Email)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"email": email,
			"uuid":  uuid,
		}).WithError(err).Warn("Unsubscribe: unable to update subscriber's status.")
		c.Redirect(http.StatusPermanentRedirect, redirWithError)
		return
	}

	c.Redirect(http.StatusPermanentRedirect, os.Getenv("APP_URL")+"/unsubscribe-success.html")
}

func ImportSubscribers(c *gin.Context) {
	u := middleware.GetUser(c)

	segments := &segmentsParam{}
	err := c.Bind(segments)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid data",
			"errors": map[string]string{
				"segments": "The segments array is in an invalid format.",
			},
		})
		return
	}

	var segs []entities.Segment
	if len(segments.Ids) > 0 {
		segs, err = storage.GetSegmentsByIDs(c, u.ID, segments.Ids)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Invalid data",
				"errors": map[string]string{
					"segments": "Unable to find the specified segments.",
				},
			})
			return
		}
	}

	filename := strings.TrimSpace(c.PostForm("filename"))
	if filename == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "The filename must not be empty.",
		})
	}

	client, err := awss3.NewS3Client(
		os.Getenv("AWS_S3_ACCESS_KEY"),
		os.Getenv("AWS_S3_SECRET_KEY"),
		os.Getenv("AWS_S3_REGION"),
	)
	if err != nil {
		logger.From(c).WithError(err).Error("Import subs: unable to create s3 client.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to import subscribers. Please try again.",
		})
		return
	}

	go func(ctx context.Context, client s3iface.S3API, filename string, userID int64, segs []entities.Segment) {
		imp := importer.NewS3SubscribersImporter(client)
		err := imp.ImportSubscribersFromFile(ctx, filename, userID, segs)
		if err != nil {
			logger.From(ctx).WithFields(logrus.Fields{
				"filename": filename,
				"segments": segs,
			}).WithError(err).Warn("Unable to import subscribers.")
		}
	}(c, client, filename, u.ID, segs)

	c.JSON(http.StatusOK, gin.H{
		"message": "We will begin processing the file shortly. As we import the subscribers, you will see them in the dashboard.",
	})
}

func BulkRemoveSubscribers(c *gin.Context) {
	u := middleware.GetUser(c)

	filename := strings.TrimSpace(c.PostForm("filename"))
	if filename == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "The filename must not be empty.",
		})
	}

	client, err := awss3.NewS3Client(
		os.Getenv("AWS_S3_ACCESS_KEY"),
		os.Getenv("AWS_S3_SECRET_KEY"),
		os.Getenv("AWS_S3_REGION"),
	)
	if err != nil {
		logger.From(c).WithError(err).Error("Remove subs: unable to create s3 client.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to import subscribers. Please try again.",
		})
		return
	}

	go func(ctx context.Context, client s3iface.S3API, filename string, userID int64) {
		svc := bulkremover.NewS3SubscribersBulkRemover(client)
		err := svc.RemoveSubscribersFromFile(ctx, filename, userID)
		if err != nil {
			logger.From(ctx).WithFields(logrus.Fields{
				"filename": filename,
			}).WithError(err).Warn("Unable to remove subscribers.")
		}
	}(c, client, filename, u.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "We will begin processing the file shortly.",
	})
}
