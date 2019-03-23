package emails

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Sender interface {
	SendTemplatedEmail(input *ses.SendTemplatedEmailInput) (*ses.SendTemplatedEmailOutput, error)
	SendBulkTemplatedEmail(input *ses.SendBulkTemplatedEmailInput) (*ses.SendBulkTemplatedEmailOutput, error)
}

func NewSesSender(key, secret, region string) (Sender, error) {
	return NewSESClient(key, secret, region)
}

func NewSESClient(key, secret, region string) (*ses.SES, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
	})

	if err != nil {
		return nil, err
	}

	return ses.New(sess), nil
}
