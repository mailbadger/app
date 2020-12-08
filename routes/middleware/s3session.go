package middleware

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/storage/s3"
)

var (
	AccessKeyID = os.Getenv("AWS_S3_ACCESS_KEY")
	SecretAccessKey= os.Getenv("AWS_S3_SECRET_KEY")
	MyRegion =os.Getenv("AWS_S3_REGION")
)

// S3Session is a middleware that inits the S3Session and attaches it to the context.
func S3Session() gin.HandlerFunc {
	return func(c *gin.Context) {
		s3.SetToContext(c, connectAws())
		c.Next()
	}
}

// connectAws creates aws session
func connectAws() *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(MyRegion),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}