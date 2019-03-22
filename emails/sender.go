package emails

import (
	"github.com/aws/aws-sdk-go/service/ses"
)

type Sender interface {
	SendTemplatedEmail(input *ses.SendTemplatedEmailInput) (*ses.SendTemplatedEmailOutput, error)
	SendBulkTemplatedEmail(input *ses.SendBulkTemplatedEmailInput) (*ses.SendBulkTemplatedEmailOutput, error)
}

func NewSesSender(key, secret, region string) (Sender, error) {
	return NewSESClient(key, secret, region)
}
