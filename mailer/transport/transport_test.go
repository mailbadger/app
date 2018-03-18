package transport_test

import (
	"os"
	"testing"

	"github.com/news-maily/api/mailer/transport"
	"github.com/stretchr/testify/assert"
)

func TestUnitEnvConfig(t *testing.T) {
	os.Setenv("MAIL_HOST", "")
	os.Setenv("MAIL_PORT", "")

	assert := assert.New(t)

	//Smtp test
	_, err := transport.NewEnvConfig("smtp")
	assert.Error(err)

	os.Setenv("MAIL_HOST", "inplayer.com")

	_, err = transport.NewEnvConfig("smtp")
	assert.Error(err)

	os.Setenv("MAIL_PORT", "25")

	conf, err := transport.NewEnvConfig("smtp")
	assert.Nil(err)
	assert.Equal(conf["host"], "inplayer.com")
	assert.Equal(conf["port"], "25")

	//Invalid driver
	_, err = transport.NewEnvConfig("foo")
	assert.EqualError(err, transport.ErrInvalidDriver.Error())
}

func TestUnitMakeTransport(t *testing.T) {
	assert := assert.New(t)

	//Smtp
	os.Setenv("MAIL_HOST", "inplayer.com")
	os.Setenv("MAIL_PORT", "abc")
	conf, err := transport.NewEnvConfig("smtp")
	trans, err := transport.MakeTransport("smtp", conf)
	assert.Error(err)

	os.Setenv("MAIL_PORT", "25")
	conf, err = transport.NewEnvConfig("smtp")
	assert.Nil(err)

	trans, err = transport.MakeTransport("smtp", conf)
	assert.Nil(err)
	assert.NotNil(trans)

	//Invalid driver
	_, err = transport.MakeTransport("foo", conf)
	assert.EqualError(err, transport.ErrInvalidDriver.Error())
}
