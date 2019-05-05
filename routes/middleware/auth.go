package middleware

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/storage"
	"github.com/news-maily/api/utils/token"
	log "github.com/sirupsen/logrus"
)

// SetUser fetches the token and then from the token fetches the user entity
// and sets it to the context.
func SetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *entities.User

		_, err := token.FromRequest(c.Request, func(t *token.Token) (string, error) {
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
