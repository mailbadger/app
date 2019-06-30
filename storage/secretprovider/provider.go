package secretprovider

import (
	"context"

	"github.com/auroratechnologies/vangoh"
	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/storage"
	"github.com/sirupsen/logrus"
)

const key = "vg"

type secretProviderImpl struct {
	s storage.Storage
}

func NewSecretProvider(s storage.Storage) vangoh.SecretProvider {
	return &secretProviderImpl{s}
}

// SetToContext sets the secret provider to the context
func SetToContext(c *gin.Context, sp vangoh.SecretProvider) {
	c.Set(key, sp)
}

// GetFromContext returns the SecretProvider associated with the context
func GetFromContext(c context.Context) vangoh.SecretProvider {
	return c.Value(key).(vangoh.SecretProvider)
}

func (sp *secretProviderImpl) GetSecret(identifier []byte) ([]byte, error) {
	key, err := sp.s.GetAccessKey(string(identifier))
	if err != nil {
		logrus.WithField("identifier", string(identifier)).WithError(err).Warn("access keys not found")
		// do not return err here, the vangoh library depends on the err to be nil
		return nil, nil
	}

	return []byte(key.SecretKey), nil
}
