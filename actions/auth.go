package actions

import (
	"net/http"
	"time"

	"github.com/news-maily/api/routes/middleware"
	"github.com/news-maily/api/storage"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/utils/token"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type tokenPayload struct {
	Access    string `json:"access_token,omitempty"`
	ExpiresIn int64  `json:"expires_in,omitempty"`
	Refresh   string `json:"refresh_token,omitempty"`
}

func PostLogin(c *gin.Context) {
	username, password := c.PostForm("username"), c.PostForm("password")

	user, err := storage.GetUserByUsername(c, username)
	if err != nil {
		logrus.Errorf("Invalid credentials. %s", err)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Invalid credentials.",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logrus.Errorf("Invalid credentials. %s", err)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Invalid credentials.",
		})
		return
	}

	exp := time.Now().Add(time.Hour * 72).Unix()
	t := token.New(token.SessionToken, user.Username)
	tokenStr, err := t.SignWithExp(user.AuthKey, exp)
	if err != nil {
		logrus.Errorf("cannot create token for %s. %s", user.Username, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to create token.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": &tokenPayload{
			Access:    tokenStr,
			ExpiresIn: exp - time.Now().Unix(), //seconds
		},
		"user": user,
	})
}

func GetMe(c *gin.Context) {
	c.JSON(http.StatusOK, middleware.GetUser(c))
}
