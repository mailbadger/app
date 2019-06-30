package middleware

import (
	"crypto"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/news-maily/app/storage/secretprovider"

	"github.com/auroratechnologies/vangoh"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/storage"
	"github.com/news-maily/app/utils/token"
	log "github.com/sirupsen/logrus"
)

const org = "MB"

func SecretProvider() gin.HandlerFunc {
	return func(c *gin.Context) {
		sp := secretprovider.NewSecretProvider(storage.GetFromContext(c), c)
		secretprovider.SetToContext(c, sp)
		c.Next()
	}
}

// SetUser fetches the token and then from the token fetches the user entity
// and sets it to the context.
func SetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var authHeader = c.GetHeader("Authorization")

		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 {
				c.Next()
				return
			}

			if parts[0] == org {
				vg := vangoh.New()
				vg.SetAlgorithm(crypto.SHA256.New)
				err := vg.AddProvider(org, secretprovider.GetFromContext(c))
				if err != nil {
					log.WithError(err).Error("unable to add secret provider")
					c.Next()
					return
				}

				err = vg.AuthenticateRequest(c.Request)
				if err != nil {
					log.WithError(err).Error("unable to authenticate api request")
					c.Next()
					return
				}
			} else if parts[0] == "Bearer" {
				var user *entities.User
				_, err := token.ParseToken(parts[1], func(t *token.Token) (string, error) {
					var err error
					secret := os.Getenv("AUTH_SECRET")
					if secret == "" {
						log.Error("auth secret is empty, unable to validate jwt.")
						return "", errors.New("auth secret is empty, unable to validate jwt")
					}
					user, err = storage.GetActiveUserByUsername(c, t.Value)
					return secret, err
				})

				if err == nil {
					c.Set("user", user)
				}
			}
		}

		c.Next()
	}
}

// GetUser returns the user set in the context
func GetUser(c *gin.Context) *entities.User {
	val, ok := c.Get("user")
	if !ok {
		return nil
	}

	user, ok := val.(*entities.User)
	if !ok {
		return nil
	}

	return user
}

// Authorized is a middleware that checks if the user is authorized to do the
// requested action.
func Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
			c.Abort()
			return
		}
		_, ok = val.(*entities.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func extractAccessKey(authHeader string) (string, error) {
	orgSplit := strings.Split(authHeader, " ")
	if len(orgSplit) < 2 {
		return "", errors.New("invalid auth header")
	}
	keySplit := strings.Split(orgSplit[1], ":")
	if len(keySplit) < 2 {
		return "", errors.New("invalid auth header")
	}

	return keySplit[0], nil
}
