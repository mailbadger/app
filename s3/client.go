package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewS3Client(key, secret, region string) (*s3.S3, error) {
	conf := &aws.Config{
		Region: aws.String(region),
	}

	if key != "" && secret != "" {
		conf.Credentials = credentials.NewStaticCredentials(key, secret, "")
	}

	sess, err := session.NewSession(conf)

	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}
