package actions

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
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
	awss3 "github.com/mailbadger/app/s3"
	"github.com/mailbadger/app/services/exporters"
	"github.com/mailbadger/app/services/reports"
	"github.com/mailbadger/app/services/subscribers/bulkremover"
	"github.com/mailbadger/app/services/subscribers/importer"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/validator"
)

const resource = "subscribers"

var (
	note = "Started the export process."
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

	scopeMap := c.QueryMap("scopes")
	err := storage.GetSubscribers(c, middleware.GetUser(c).ID, p, scopeMap)
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

func PostSubscriber(c *gin.Context) {
	var err error
	body := &params.PostSubscriber{}

	if err = c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters, please try again",
		})
		return
	}

	body.Metadata = c.PostFormMap("metadata")

	if err = validator.Validate(body); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	s := &entities.Subscriber{
		Name:     body.Name,
		Email:    body.Email,
		Metadata: body.Metadata,
		Active:   true,
		UserID:   middleware.GetUser(c).ID,
	}

	s.Segments, err = storage.GetSegmentsByIDs(c, s.UserID, body.SegmentIDs)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid data",
			"errors": map[string]string{
				"segments": "Unable to find the specified segments.",
			},
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

	s.MetaJSON, err = json.Marshal(body.Metadata)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to create subscriber, invalid metadata.",
		})
		return
	}

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

	body := &params.PutSubscriber{}
	if err = c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters, please try again",
		})
		return
	}

	body.Metadata = c.PostFormMap("metadata")
	if err = validator.Validate(body); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	segments, err := storage.GetSegmentsByIDs(c, s.UserID, body.SegmentIDs)
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
	body := &params.PostUnsubscribe{}
	if err := c.ShouldBind(body); err != nil {
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

	u, err := storage.GetUserByUUID(c, body.UUID)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"email": body.Email,
			"uuid":  body.UUID,
		}).WithError(err).Warn("Unsubscribe: cannot find user by uuid.")
		c.Redirect(http.StatusPermanentRedirect, redirWithError)
		return
	}

	sub, err := storage.GetSubscriberByEmail(c, body.Email, u.ID)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"email": body.Email,
			"uuid":  body.UUID,
		}).WithError(err).Warn("Unsubscribe: unable to fetch subscriber by email.")
		c.Redirect(http.StatusPermanentRedirect, redirWithError)
		return
	}

	hash, err := sub.GenerateUnsubscribeToken(os.Getenv("UNSUBSCRIBE_SECRET"))
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"email": body.Email,
			"uuid":  body.UUID,
		}).WithError(err).Error("Unsubscribe: unable to generate hash.")
		c.Redirect(http.StatusPermanentRedirect, redirWithError)
		return
	}

	if subtle.ConstantTimeCompare([]byte(body.Token), []byte(hash)) != 1 {
		logger.From(c).WithFields(logrus.Fields{
			"email": body.Email,
			"uuid":  body.UUID,
		}).WithError(err).Warn("Unsubscribe: hashes don't match.")
		c.Redirect(http.StatusPermanentRedirect, redirWithError)
		return
	}

	err = storage.DeactivateSubscriber(c, u.ID, sub.Email)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"email": body.Email,
			"uuid":  body.UUID,
		}).WithError(err).Warn("Unsubscribe: unable to update subscriber's status.")
		c.Redirect(http.StatusPermanentRedirect, redirWithError)
		return
	}

	c.Redirect(http.StatusPermanentRedirect, os.Getenv("APP_URL")+"/unsubscribe-success.html")
}

func ImportSubscribers(c *gin.Context) {
	u := middleware.GetUser(c)

	body := &params.ImportSubscribers{}
	err := c.ShouldBind(body)
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

	var segs []entities.Segment
	if len(body.SegmentIDs) > 0 {
		segs, err = storage.GetSegmentsByIDs(c, u.ID, body.SegmentIDs)
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
	}(c, client, body.Filename, u.ID, segs)

	c.JSON(http.StatusOK, gin.H{
		"message": "We will begin processing the file shortly. As we import the subscribers, you will see them in the dashboard.",
	})
}

func BulkRemoveSubscribers(c *gin.Context) {
	u := middleware.GetUser(c)

	body := &params.BulkRemoveSubscribers{}
	err := c.ShouldBind(body)
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
	}(c, client, body.Filename, u.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "We will begin processing the file shortly.",
	})
}

func ExportSubscribers(c *gin.Context) {
	u := middleware.GetUser(c)

	s3, err := awss3.NewS3Client(
		os.Getenv("AWS_S3_ACCESS_KEY"),
		os.Getenv("AWS_S3_SECRET_KEY"),
		os.Getenv("AWS_S3_REGION"),
	)
	if err != nil {
		logger.From(c).WithError(err).Error("Import subs: unable to create s3 client.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to export subscribers. Please try again.",
		})
		return
	}

	exporter, err := exporters.NewExporter("subscribers", s3)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable do create subscribers exporter")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to export subscribers. Please try again.",
		})
		return
	}

	reportSvc := reports.NewReportService(exporter)

	report, err := reportSvc.CreateExportReport(c, u.ID, resource, note, time.Now())
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
		report, err = reportSvc.GenerateExportReport(c, u.ID, report)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{
				"report": report,
			}).WithError(err).Errorf("Export failed")
		}
	}(c.Copy(), report)

	c.JSON(http.StatusOK, report)
}

func DownloadSubscribersReport(c *gin.Context) {
	u := middleware.GetUser(c)

	fileName := c.Query("filename")

	report, err := storage.GetReportByFilename(c, fileName, u.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Report not found.",
		})
		return
	}

	if report.Status == entities.StatusFailed {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Failed to generate report, please try again.",
		})
		return
	}

	if report.Status == entities.StatusInProgress {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Generating report, please try again later.",
		})
		return
	}

	if report.Status == entities.StatusDone {
		client, err := awss3.NewS3Client(
			os.Getenv("AWS_S3_ACCESS_KEY"),
			os.Getenv("AWS_S3_SECRET_KEY"),
			os.Getenv("AWS_S3_REGION"),
		)
		if err != nil {
			logger.From(c).WithError(err).Error("Unable to create s3 client.")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to sign url.",
			})
			return
		}

		req, _ := client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
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
