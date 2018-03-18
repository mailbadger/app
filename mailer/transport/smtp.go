package transport

import (
	"errors"

	gomail "gopkg.in/gomail.v2"
)

type Dialer interface {
	DialAndSend(m ...*gomail.Message) error
}

type SmtpTransport struct {
	Client Dialer
}

func (smtp *SmtpTransport) Send(m *Message) error {
	message, err := smtp.NewMessageFrom(m)
	if err != nil {
		return err
	}

	return smtp.Client.DialAndSend(message)
}

func (smtp *SmtpTransport) NewMessageFrom(m *Message) (*gomail.Message, error) {
	if m.Email.To.Email == "" {
		return nil, errors.New("email cannot be empty")
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.Email.From)
	message.SetHeader("To", m.Email.To.Email)
	message.SetHeader("Subject", m.Email.Subject)
	message.SetBody("text/html", m.Body)

	if len(m.Email.Cc) != 0 {
		message.SetHeader("Cc", m.Email.Cc...)
	}

	if len(m.Email.Bcc) != 0 {
		message.SetHeader("Bcc", m.Email.Bcc...)
	}

	return message, nil
}
