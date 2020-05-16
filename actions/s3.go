package actions

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/s3"
)

func GetSignedURL(c *gin.Context) {
	u := middleware.GetUser(c)

	filename := strings.TrimSpace(c.PostForm("filename"))
	contentType := strings.TrimSpace(c.PostForm("contentType"))
	action := strings.ToLower(strings.TrimSpace(c.PostForm("action")))

	if action != "import" && action != "export" && action != "remove" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to sign url. Invalid action.",
		})
		return
	}

	client, err := s3.NewS3Client(
		os.Getenv("AWS_S3_ACCESS_KEY"),
		os.Getenv("AWS_S3_SECRET_KEY"),
		os.Getenv("AWS_S3_REGION"),
	)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to create s3 client.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to sign url.",
		})
		return
	}
	req, _ := client.PutObjectRequest(&awss3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:         aws.String(fmt.Sprintf("subscribers/%s/%d/%s", action, u.ID, filename)),
		ContentType: aws.String(contentType),
	})

	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		logger.From(c).WithError(err).Warn("Unable to sign s3 url.")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to sign url.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":    url,
		"method": req.Operation.HTTPMethod,
		"headers": map[string]string{
			"content-type": contentType,
		},
	})
}
