package emails

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Sender interface {
	SendTemplatedEmail(input *ses.SendTemplatedEmailInput) (*ses.SendTemplatedEmailOutput, error)
	SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error)
	SendBulkTemplatedEmail(input *ses.SendBulkTemplatedEmailInput) (*ses.SendBulkTemplatedEmailOutput, error)
	CreateConfigurationSet(input *ses.CreateConfigurationSetInput) (*ses.CreateConfigurationSetOutput, error)
	DescribeConfigurationSet(input *ses.DescribeConfigurationSetInput) (*ses.DescribeConfigurationSetOutput, error)
	CreateConfigurationSetEventDestination(input *ses.CreateConfigurationSetEventDestinationInput) (*ses.CreateConfigurationSetEventDestinationOutput, error)
	DeleteConfigurationSet(input *ses.DeleteConfigurationSetInput) (*ses.DeleteConfigurationSetOutput, error)
	GetSendQuota(input *ses.GetSendQuotaInput) (*ses.GetSendQuotaOutput, error)
}

type senderImpl struct {
	*ses.SES
}

// SES Notification Types
const (
	SendType             = "Send"
	ClickType            = "Click"
	OpenType             = "Open"
	BounceType           = "Bounce"
	DeliveryType         = "Delivery"
	ComplaintType        = "Complaint"
	RenderingFailureType = "Rendering Failure"
	SubConfirmationType  = "SubscriptionConfirmation"

	ConfigurationSetName = "MailbadgerEvents"
)

func NewSESClient(key, secret, region string) (*ses.SES, error) {
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

	return ses.New(sess), nil
}

func NewSesSender(key, secret, region string) (Sender, error) {
	client, err := NewSESClient(key, secret, region)
	if err != nil {
		return nil, err
	}

	return &senderImpl{client}, nil
}
