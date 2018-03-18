package transport_test

import (
	"encoding/json"
	"fmt"
	"testing"

	gomail "gopkg.in/gomail.v2"

	"github.com/FilipNikolovski/news-maily/mailer/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SmtpMock struct {
	mock.Mock
}

func (c *SmtpMock) DialAndSend(m ...*gomail.Message) error {
	args := c.Called(m)
	return args.Error(0)
}

func TestUnitNewMessageFrom(t *testing.T) {
	jsonMsg := `
	{
		"subject": "Welcome to InPlayer!",
		"from": "noreply@inplayer.com",
		"name": "InPlayer",
		"template": "consumer-account-registration",
		"templateData": {
		"name": "filipn",
		"confirmation_token": "",
		"client_activate_url": ""
		},
		"to": {
			"email": "filip.nikolovski+3@inplayer.com",
			"name": "filipn"
		},
		"cc": [],
		"bcc": []
	}
	`

	assert := assert.New(t)
	email := transport.Email{}
	msg := &transport.Message{}

	err := json.Unmarshal([]byte(jsonMsg), &email)
	assert.Nil(err)

	msg.Email = email
	smtp := &transport.SmtpTransport{}

	message, err := smtp.NewMessageFrom(msg)
	assert.Nil(err)
	assert.Equal(message.GetHeader("From")[0], msg.Email.From)
	assert.Equal(message.GetHeader("Subject")[0], msg.Email.Subject)

	assert.Equal(message.GetHeader("To")[0], msg.Email.To.Email)

}

func TestUnitSend(t *testing.T) {
	jsonMsg := `
	{
		"subject": "Welcome to InPlayer!",
		"from": "noreply@inplayer.com",
		"name": "InPlayer",
		"merchantId": 21,
		"template": "consumer-account-registration",
		"templateData": {
		"name": "filipn",
		"confirmation_token": "",
		"client_activate_url": ""
		},
		"to": {
			"email": "filip.nikolovski+3@inplayer.com",
			"name": "filipn"
		},
		"cc": [],
		"bcc": []
	}
	`

	assert := assert.New(t)
	email := transport.Email{}
	msg := &transport.Message{}

	err := json.Unmarshal([]byte(jsonMsg), &email)
	assert.Nil(err)

	msg.Email = email
	clientMock := new(SmtpMock)
	smtp := &transport.SmtpTransport{
		Client: clientMock,
	}

	fmt.Println(msg)
	clientMock.On("DialAndSend", mock.AnythingOfType("[]*gomail.Message")).Return(nil)
	err = smtp.Send(msg)
	assert.Nil(err)

	msg.Email.To = transport.To{}

	err = smtp.Send(msg)
	assert.Error(err)
}
