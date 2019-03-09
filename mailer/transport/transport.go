package transport

import (
	"errors"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	gomail "gopkg.in/gomail.v2"
)

type Transporter interface {
	Send(m *Message) error
}

type Message struct {
	Email Email
	Body  string
}

type Email struct {
	Subject      string                 `json:"subject"`
	From         string                 `json:"from"`
	Name         string                 `json:"name"`
	Template     string                 `json:"template"`
	TemplateData map[string]interface{} `json:"templateData"`
	Cc           []string               `json:"cc"`
	Bcc          []string               `json:"bcc"`
	To           []To                   `json:"to"`
}

type To struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Transport error messages.
var (
	ErrInvalidDriver = errors.New("cannot create transport, invalid email driver")
	ErrInvalidHost   = errors.New("cannot create config, mail host not found")
	ErrInvalidPort   = errors.New("cannot create config, mail port not found")
)

// NewSMTPTransport returns new SmtpTransport object.
func NewSMTPTransport(conf map[string]string) (Transporter, error) {
	port, err := strconv.Atoi(conf["port"])
	if err != nil {
		return nil, err
	}

	client := gomail.NewPlainDialer(conf["host"], port, conf["username"], conf["password"])

	return &SmtpTransport{
		Client: client,
	}, nil
}

// NewSesTransport returns new SesTransport object.
func NewSesTransport(conf map[string]string) (Transporter, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf["region"]),
		Credentials: credentials.NewStaticCredentials(conf["key"], conf["secret"], ""),
	})

	if err != nil {
		return nil, err
	}

	return &SesTransport{
		Client: ses.New(sess),
	}, nil
}

// MakeTransport returns new Transporter object based on the specified driver.
// Available drivers are: "mailgun", "smtp".
// Returns ErrInvalidDriver if the driver specified is invalid.
func MakeTransport(driver string, conf map[string]string) (Transporter, error) {
	switch driver {
	case "smtp":
		return NewSMTPTransport(conf)
	case "ses":
		return NewSesTransport(conf)
	default:
		return nil, ErrInvalidDriver
	}
}

// NewEnvConfig returns config map based on the specified driver.
// Available drivers are: "mailgun", "smtp".
// Returns ErrInvalidDriver if the driver specified is invalid.
func NewEnvConfig(driver string) (map[string]string, error) {
	switch driver {
	case "smtp":
		host := os.Getenv("MAIL_HOST")
		if host == "" {
			return nil, ErrInvalidHost
		}

		port := os.Getenv("MAIL_PORT")
		if port == "" {
			return nil, ErrInvalidPort
		}

		return map[string]string{
			"host": host,
			"port": port,
		}, nil
	default:
		return nil, ErrInvalidDriver
	}
}
