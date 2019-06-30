package secretprovider

import (
	"context"
	"net/http"

	"github.com/news-maily/app/entities"

	"github.com/auroratechnologies/vangoh"
	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/storage"
	"github.com/sirupsen/logrus"
)

const key = "vg"

type secretProviderImpl struct {
	s storage.Storage
	c *gin.Context
}

func NewSecretProvider(s storage.Storage, c *gin.Context) vangoh.SecretProviderWithCallback {
	return &secretProviderImpl{s, c}
}

// SetToContext sets the secret provider to the context
func SetToContext(c *gin.Context, sp vangoh.SecretProviderWithCallback) {
	c.Set(key, sp)
}

// GetFromContext returns the SecretProvider associated with the context
func GetFromContext(c context.Context) vangoh.SecretProvider {
	return c.Value(key).(vangoh.SecretProvider)
}

func (sp *secretProviderImpl) GetSecret(identifier []byte, cbPayload *vangoh.CallbackPayload) ([]byte, error) {
	key, err := sp.s.GetAccessKey(string(identifier))
	if err != nil {
		logrus.WithField("identifier", string(identifier)).WithError(err).Warn("access keys not found")
		// do not return err here, the vangoh library depends on the err to be nil
		return nil, nil
	}

	cbPayload.SetPayload(&key.User)

	return []byte(key.SecretKey), nil
}

func (sp *secretProviderImpl) SuccessCallback(r *http.Request, cbPayload *vangoh.CallbackPayload) {
	user, ok := cbPayload.GetPayload().(*entities.User)
	if !ok {
		logrus.Error("unable to retreive user from callback payload")
		return
	}

	sp.c.Set("user", user)
}
