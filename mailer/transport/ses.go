package transport

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SesClient interface {
	SendEmail(m *ses.SendEmailInput) (*ses.SendEmailOutput, error)
}

type SesTransport struct {
	Client SesClient
}

func (t *SesTransport) Send(m *Message) (string, error) {
	message, err := t.NewMessageFrom(m)
	if err != nil {
		return "", err
	}

	out, err := t.Client.SendEmail(message)

	return *out.MessageId, err
}

func (t *SesTransport) NewMessageFrom(m *Message) (*ses.SendEmailInput, error) {
	if len(m.Email.To) == 0 {
		return nil, errors.New("to value is empty, cannot create message")
	}

	charset := aws.String("UTF-8")

	var toAddresses []string
	for _, to := range m.Email.To {
		toAddresses = append(toAddresses, to.Email)
	}

	return &ses.SendEmailInput{
		Source: aws.String(m.Email.From),

		Destination: &ses.Destination{
			ToAddresses:  aws.StringSlice(toAddresses),
			CcAddresses:  aws.StringSlice(m.Email.Cc),
			BccAddresses: aws.StringSlice(m.Email.Bcc),
		},

		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: charset,
				Data:    aws.String(m.Email.Subject),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: charset,
					Data:    aws.String(m.Body),
				},
			},
		},
	}, nil
}
