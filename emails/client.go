package emails

import (
	"github.com/news-maily/api/emails/transport"
)

type Mailer interface {
	Send(m *transport.Message) error
}

type mailerImpl struct {
	transporter transport.Transporter
}

func NewClient(t transport.Transporter) Mailer {
	return &mailerImpl{
		transporter: t,
	}
}

func (mailer *mailerImpl) Send(m *transport.Message) error {
	return mailer.transporter.Send(m)
}
