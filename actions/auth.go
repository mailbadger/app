package actions

import (
	"net/http"
	"time"

	"github.com/FilipNikolovski/news-maily/routes/middleware"
	"github.com/FilipNikolovski/news-maily/storage"

	"github.com/FilipNikolovski/news-maily/utils/token"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type tokenPayload struct {
	Access    string `json:"access_token,omitempty"`
	ExpiresIn int64  `json:"expires_in,omitempty"`
	Refresh   string `json:"refresh_token,omitempty"`
}

func PostLogin(c *gin.Context) {
	username, password := c.Request.PostFormValue("username"), c.Request.PostFormValue("password")

	user, err := storage.GetUserByUsername(c, username)
	if err != nil {
		logrus.Errorf("Invalid credentials. %s", err)
		c.JSON(http.StatusForbidden, gin.H{
			"reason": "Invalid credentials.",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logrus.Errorf("Invalid credentials. %s", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"reason": "Invalid credentials.",
		})
		return
	}

	exp := time.Now().Add(time.Hour * 72).Unix()
	t := token.New(token.SessionToken, user.Username)
	tokenStr, err := t.SignWithExp(user.AuthKey, exp)
	if err != nil {
		logrus.Errorf("cannot create token for %s. %s", user.Username, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Unable to create token.",
		})
		return
	}

	c.JSON(http.StatusOK, &tokenPayload{
		Access:    tokenStr,
		ExpiresIn: exp - time.Now().Unix(), //seconds
	})
}

func GetMe(c *gin.Context) {
	c.JSON(http.StatusOK, middleware.GetUser(c))
}
