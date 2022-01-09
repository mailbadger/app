package s3

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewClient() (*s3.S3, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}
