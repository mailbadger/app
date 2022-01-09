package actions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/validator"
)

func GetSignedURL(client s3iface.S3API, bucket string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := middleware.GetUser(c)

		body := &params.GetSignedURL{}
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

		req, _ := client.PutObjectRequest(&awss3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(fmt.Sprintf("subscribers/%s/%d/%s", body.Action, u.ID, body.Filename)),
			ContentType: aws.String(body.ContentType),
		})

		url, err := req.Presign(15 * time.Minute)
		if err != nil {
			logger.From(c).WithError(err).Error("Unable to sign s3 url.")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Unable to sign url.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"url":    url,
			"method": req.Operation.HTTPMethod,
			"headers": map[string]string{
				"content-type": body.ContentType,
			},
		})
	}
}
