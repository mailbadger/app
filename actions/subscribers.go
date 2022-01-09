package actions

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/services/boundaries"
	"github.com/mailbadger/app/services/reports"
	"github.com/mailbadger/app/services/subscribers"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/utils"
	"github.com/mailbadger/app/validator"
)

func GetSubscribers(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("cursor")
		if !ok {
			logger.From(c).Error("get subs: unable to fetch pagination cursor from context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch subscribers. Please try again.",
			})
			return
		}

		p, ok := val.(*storage.PaginationCursor)
		if !ok {
			logger.From(c).Error("get subs: unable to cast pagination cursor from context value")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch subscribers. Please try again.",
			})
			return
		}

		scopeMap := c.QueryMap("scopes")
		err := store.GetSubscribers(middleware.GetUser(c).ID, p, scopeMap)
		if err != nil {
			logger.From(c).WithError(err).Error("Unable to fetch subscribers collection.")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to fetch subscribers. Please try again.",
			})
			return
		}

		c.JSON(http.StatusOK, p)
	}
}

func GetSubscriber(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer.",
			})
		}

		s, err := storage.GetSubscriber(id, middleware.GetUser(c).ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Subscriber not found.",
			})
			return
		}

		c.JSON(http.StatusOK, s)
	}
}

func PostSubscriber(boundarysvc boundaries.Service, storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		body := &params.PostSubscriber{}

		if err = c.ShouldBindJSON(body); err != nil {
			logger.From(c).WithError(err).Error("post subscriber: unable to bind body")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid parameters, please try again.",
			})
			return
		}

		if err = validator.Validate(body); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		user := middleware.GetUser(c)

		limitexceeded, _, err := boundarysvc.SubscribersLimitExceeded(user)
		if err != nil {
			logger.From(c).WithError(err).Error("post subscriber: unable to check subscribers limit for user")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to check subscribers limit. Please try again.",
			})
			return
		}

		if limitexceeded {
			logger.From(c).Info("post subscriber: user has exceeded his subscribers limit")
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You have exceeded your subscribers limit, please upgrade to a bigger plan or contact support.",
			})
			return
		}

		s := &entities.Subscriber{
			Name:     body.Name,
			Email:    body.Email,
			Metadata: body.Metadata,
			Active:   true,
			UserID:   middleware.GetUser(c).ID,
		}

		s.Segments, err = storage.GetSegmentsByIDs(s.UserID, body.SegmentIDs)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Invalid data",
				"errors": map[string]string{
					"segments": "Unable to find the specified segments.",
				},
			})
			return
		}

		_, err = storage.GetSubscriberByEmail(s.Email, s.UserID)
		if err == nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Subscriber with that email already exists.",
			})
			return
		}

		s.MetaJSON, err = json.Marshal(body.Metadata)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to create subscriber, invalid metadata.",
			})
			return
		}

		if err := storage.CreateSubscriber(s); err != nil {
			logger.From(c).
				WithError(err).
				Error("post subscriber: unable to create subscriber")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to create subscriber",
			})
			return
		}

		c.JSON(http.StatusCreated, s)
	}
}

func PutSubscriber(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer.",
			})
			return
		}

		s, err := storage.GetSubscriber(id, middleware.GetUser(c).ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Subscriber not found.",
			})
			return
		}

		body := &params.PutSubscriber{}
		if err = c.ShouldBindJSON(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid parameters, please try again.",
			})
			return
		}

		if err = validator.Validate(body); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		segments, err := storage.GetSegmentsByIDs(s.UserID, body.SegmentIDs)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Invalid data",
				"errors": map[string]string{
					"segments": "Unable to find the specified segments.",
				},
			})
			return
		}

		metaJSON, err := json.Marshal(body.Metadata)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to create subscriber, invalid metadata.",
			})
			return
		}

		s.Name = body.Name
		s.MetaJSON = metaJSON
		s.Segments = segments

		if err = storage.UpdateSubscriber(s); err != nil {
			logger.From(c).
				WithError(err).
				WithField("subscriber_id", id).
				Error("put subscriber: unable to update subscriber")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to update subscriber.",
			})
			return
		}

		c.JSON(http.StatusOK, s)
	}
}

func DeleteSubscriber(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Id must be an integer.",
			})
		}
		user := middleware.GetUser(c)

		err = storage.DeleteSubscriber(id, user.ID)
		if err != nil {
			logger.From(c).WithError(err).Error("delete subscriber: unable to delete subscriber")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to delete subscriber.",
			})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func PostUnsubscribe(
	storage storage.Storage,
	unsubscribeSecret string,
	appURL string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		body := &params.PostUnsubscribe{}
		if err := c.ShouldBindJSON(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid parameters, please try again",
			})
			return
		}

		if err := validator.Validate(body); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		redirWithError := c.Request.Referer()

		params := url.Values{}
		params.Add("email", body.Email)
		params.Add("uuid", body.UUID)
		params.Add("t", body.Token)
		params.Add("failed", "true")

		redirWithError = redirWithError + "?" + params.Encode()

		u, err := storage.GetUserByUUID(body.UUID)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{
				"email": body.Email,
				"uuid":  body.UUID,
			}).WithError(err).Error("Unsubscribe: cannot find user by uuid")

			c.Redirect(http.StatusTemporaryRedirect, redirWithError)
			return
		}

		sub, err := storage.GetSubscriberByEmail(body.Email, u.ID)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{
				"email": body.Email,
				"uuid":  body.UUID,
			}).WithError(err).Error("Unsubscribe: unable to fetch subscriber by email")

			c.Redirect(http.StatusTemporaryRedirect, redirWithError)
			return
		}

		hash, err := sub.GenerateUnsubscribeToken(unsubscribeSecret)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{
				"email": body.Email,
				"uuid":  body.UUID,
			}).WithError(err).Error("Unsubscribe: unable to generate hash")

			c.Redirect(http.StatusTemporaryRedirect, redirWithError)
			return
		}

		if len(body.Token) != len(hash) {
			logger.From(c).WithFields(logrus.Fields{
				"email": body.Email,
				"uuid":  body.UUID,
			}).WithError(err).Warn("Unsubscribe: hashes are not the same length")

			c.Redirect(http.StatusTemporaryRedirect, redirWithError)
			return
		}

		if subtle.ConstantTimeCompare([]byte(body.Token), []byte(hash)) != 1 {
			logger.From(c).WithFields(logrus.Fields{
				"email": body.Email,
				"uuid":  body.UUID,
			}).WithError(err).Warn("Unsubscribe: hashes don't match.")

			c.Redirect(http.StatusTemporaryRedirect, redirWithError)
			return
		}

		err = storage.DeactivateSubscriber(u.ID, body.Email)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{
				"email": body.Email,
				"uuid":  body.UUID,
			}).WithError(err).Error("Unsubscribe: unable to deactivate subscriber")

			c.Redirect(http.StatusTemporaryRedirect, redirWithError)
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, appURL+"/unsubscribe-success.html")
	}
}

func ImportSubscribers(
	subscrsvc subscribers.Service,
	boundarysvc boundaries.Service,
	storage storage.Storage,
	s3Client s3iface.S3API,
	bucket string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := middleware.GetUser(c)

		reqParams := &params.ImportSubscribers{}
		err := c.ShouldBindJSON(reqParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid parameters, please try again.",
			})
			return
		}

		if err := validator.Validate(reqParams); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		var segs []entities.Segment
		if len(reqParams.SegmentIDs) > 0 {
			segs, err = storage.GetSegmentsByIDs(u.ID, reqParams.SegmentIDs)
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

		limitExceeded, count, err := boundarysvc.SubscribersLimitExceeded(u)
		if err != nil {
			logger.From(c).WithError(err).Error("import subscribers: unable to fetch user limits")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to import subscribers. Please try again.",
			})
			return
		}

		if limitExceeded {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You have exceeded your subscribers limit, please upgrade to a bigger plan or contact support.",
			})
			return
		}

		res, err := s3Client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fmt.Sprintf("subscribers/import/%d/%s", u.ID, reqParams.Filename)),
		})
		if err != nil {
			logger.From(c).WithError(err).Error("import subscribers: unable to fetch import file")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to import subscribers. Please try again.",
			})
			return
		}

		defer func() {
			err = res.Body.Close()
			if err != nil {
				logger.From(c).WithError(err).Error("import subscribers: unable to close body")
			}
		}()

		// duplicate stream read.
		var buf bytes.Buffer
		tee := io.TeeReader(res.Body, &buf)

		csvCount, err := utils.CountLines(tee)
		if err != nil {
			logger.From(c).WithError(err).Error("import subscribers: unable to count lines")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to import subscribers. Please try again.",
			})
			return
		}

		if count+int64(csvCount) > u.Boundaries.SubscribersLimit {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "With this import you will exceed the limit of your subscribers, update your plan or contact the support team.",
				"total":   count,
				"count":   csvCount,
			})
			return
		}

		go func(
			ctx context.Context,
			svc subscribers.Service,
			userID int64,
			segs []entities.Segment,
			r io.Reader,
		) {
			err := svc.ImportSubscribersFromFile(ctx, u.ID, segs, r)
			if err != nil {
				logger.From(ctx).WithFields(logrus.Fields{
					"segments": segs,
				}).WithError(err).Error("import subscribers: unable to import subscribers from file")
			}
		}(c.Copy(), subscrsvc, u.ID, segs, &buf)

		c.JSON(http.StatusOK, gin.H{
			"message": "We will begin processing the file shortly. As we import the subscribers, you will see them in the dashboard.",
		})
	}
}

func BulkRemoveSubscribers(svc subscribers.Service, s3Client s3iface.S3API, bucket string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := middleware.GetUser(c)

		body := &params.BulkRemoveSubscribers{}
		err := c.ShouldBindJSON(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid parameters, please try again",
			})
			return
		}

		if err := validator.Validate(body); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		res, err := s3Client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fmt.Sprintf("subscribers/remove/%d/%s", u.ID, body.Filename)),
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to remove subscribers. Please try again.",
			})
			return
		}

		go func(ctx context.Context, svc subscribers.Service, filename string, userID int64, r io.ReadCloser) {
			err := svc.RemoveSubscribersFromFile(ctx, filename, userID, r)
			if err != nil {
				logger.From(ctx).WithFields(logrus.Fields{
					"filename": filename,
				}).WithError(err).Error("delete subs: wnable to remove subscribers")
			}
		}(c, svc, body.Filename, u.ID, res.Body)

		c.JSON(http.StatusOK, gin.H{
			"message": "We will begin processing the file shortly.",
		})
	}
}

func ExportSubscribers(reportsvc reports.Service, bucket string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const (
			resource = "subscribers"
			note     = "Started the export process."
		)

		u := middleware.GetUser(c)

		report, err := reportsvc.CreateExportReport(c, u.ID, resource, note, time.Now())
		if err != nil {
			switch {
			case errors.Is(err, reports.ErrAnotherReportRunning):
				logger.From(c).WithFields(logrus.Fields{
					"user_id":  u.ID,
					"resource": resource,
					"note":     note,
				}).WithError(err).Info("There is a report already running for this user")
				c.JSON(http.StatusForbidden, gin.H{
					"message": "There is a report already running.",
				})
			case errors.Is(err, reports.ErrLimitReached):
				logger.From(c).WithFields(logrus.Fields{
					"user_id":  u.ID,
					"resource": resource,
					"note":     note,
				}).WithError(err).Info("This user reached the daily limit")
				c.JSON(http.StatusForbidden, gin.H{
					"message": "You reached the daily limit, unable to generate report.",
				})
			default:
				logger.From(c).WithFields(logrus.Fields{
					"user_id":  u.ID,
					"resource": resource,
					"note":     note,
				}).WithError(err).Error("Unable to create export report service")
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Unable to create export report.",
				})
			}
			return
		}

		go func(c context.Context, report *entities.Report) {
			report, err = reportsvc.GenerateExportReport(c, u.ID, report, bucket)
			if err != nil {
				logger.From(c).WithFields(logrus.Fields{
					"report": report,
				}).WithError(err).Error("Export failed")
			}
		}(c.Copy(), report)

		c.JSON(http.StatusOK, report)
	}
}

func DownloadSubscribersReport(storage storage.Storage, s3Client s3iface.S3API, bucket string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := middleware.GetUser(c)

		fileName := c.Query("filename")

		report, err := storage.GetReportByFilename(fileName, u.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Report not found.",
			})
			return
		}

		if report.Status == entities.StatusFailed {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Failed to generate report, please try again.",
				"status":  report.Status,
			})
			return
		}

		if report.Status == entities.StatusInProgress {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Generating report, please try again later.",
				"status":  report.Status,
			})
			return
		}

		if report.Status == entities.StatusDone {
			req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(fmt.Sprintf("subscribers/export/%d/%s", u.ID, fileName)),
			})

			pUrl, err := req.Presign(15 * time.Minute)
			if err != nil {
				logger.From(c).WithError(err).Warn("Unable to sign s3 url.")
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Unable to sign url.",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"url": pUrl,
			})
		}
	}
}
