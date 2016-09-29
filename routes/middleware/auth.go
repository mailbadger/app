package middleware

import (
	"strconv"

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

		t, err := token.FromRequest(c.Request)

		if err == nil {
			if id, err := strconv.ParseInt(t.Value, 10, 64); err == nil {
				if user, err = storage.GetUser(c, id); err == nil {
					c.Set("user", user)
				}
			}
		}

		c.Next()
	}
}

// Authorized is a middleware that checks if the user is authorized to do the
// requested action.
func Authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("user")
		if !ok {
			c.String(401, "User not authorized")
			c.Abort()
		}
		_, ok = val.(*entities.User)
		if !ok {
			c.String(401, "User not authorized")
			c.Abort()
		}

		c.Next()
	}
}
