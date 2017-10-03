package middleware

import (
	"net/http"

	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/storage"
	"github.com/FilipNikolovski/news-maily/utils/token"
	"github.com/gin-gonic/gin"
)

// SetUser fetches the token and then from the token fetches the user entity
// and sets it to the context.
func SetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *entities.User

		_, err := token.FromRequest(c.Request, func(t *token.Token) (string, error) {
			var err error
			user, err = storage.GetUserByUsername(c, t.Value)
			return user.AuthKey, err
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
			c.JSON(http.StatusUnauthorized, gin.H{"reason": "User not authorized"})
			c.Abort()
			return
		}
		_, ok = val.(*entities.User)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"reason": "User not authorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
