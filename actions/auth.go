package actions

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/storage"
	"github.com/news-maily/api/utils/token"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/ezzarghili/recaptcha-go.v3"
)

type tokenPayload struct {
	Access    string `json:"access_token,omitempty"`
	ExpiresIn int64  `json:"expires_in,omitempty"`
	Refresh   string `json:"refresh_token,omitempty"`
}

func PostAuthenticate(c *gin.Context) {
	username, password := c.PostForm("username"), c.PostForm("password")

	user, err := storage.GetActiveUserByUsername(c, username)
	if err != nil {
		logrus.Errorf("Invalid credentials. %s", err)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Invalid credentials.",
		})
		return
	}

	if !user.Password.Valid {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Invalid credentials. Most likely your account was created using one of the oauth providers. Try a different authentication method.",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
	if err != nil {
		logrus.Errorf("Invalid credentials. %s", err)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Invalid credentials.",
		})
		return
	}

	exp := time.Now().Add(time.Hour * 72).Unix()
	t := token.New(token.SessionToken, user.Username)
	tokenStr, err := t.SignWithExp(os.Getenv("AUTH_SECRET"), exp)
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

type signupParams struct {
	Username      string `form:"username" valid:"email,required~Email is blank or in invalid format"`
	Password      string `form:"password" valid:"required"`
	TokenResponse string `form:"token_response" valid:"optional"`
}

func PostSignup(c *gin.Context) {
	params := &signupParams{}
	c.Bind(params)

	v, err := valid.ValidateStruct(params)
	if !v {
		msg := "Unable to create account, some parameters are invalid."
		if err != nil {
			msg = err.Error()
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": msg,
		})
		return
	}

	if len(params.Password) < 8 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"new_password": "The password must be atleast 8 characters in length.",
		})
		return
	}

	_, err = storage.GetUserByUsername(c, params.Username)
	if err == nil {
		logrus.WithField("username", params.Username).Errorf("duplicate account %s", err)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Unable to create an account.",
		})
		return
	}

	secret := os.Getenv("RECAPTCHA_SECRET")
	if secret != "" {
		captcha, err := recaptcha.NewReCAPTCHA(secret, recaptcha.V2, 10*time.Second)
		if err != nil {
			logrus.Errorf("recaptcha initialize error %s", err.Error())
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unable to create an account.",
			})
			return
		}

		err = captcha.Verify(params.TokenResponse)
		if err != nil {
			logrus.WithField("username", params.Username).Errorf("recaptcha invalid response. %s", err)
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unable to create an account.",
			})
			return
		}
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		logrus.Errorf("unable to generate random uuid: %s", err.Error())
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Unable to create an account.",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("unable to generate hash from password %s", err.Error())
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Unable to create an account.",
		})
		return
	}

	user := &entities.User{
		Username: params.Username,
		UUID:     uuid.String(),
		Password: sql.NullString{
			String: string(hashedPassword),
			Valid:  true,
		},
		Active:   true,
		Verified: false,
		Source:   "mailbadger.io",
	}

	err = storage.CreateUser(c, user)
	if err != nil {
		logrus.WithField("username", params.Username).Errorf("unable to persist user in db %s", err.Error())
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Unable to create an account.",
		})
		return
	}

	exp := time.Now().Add(time.Hour * 72).Unix()
	t := token.New(token.SessionToken, user.Username)
	tokenStr, err := t.SignWithExp(os.Getenv("AUTH_SECRET"), exp)
	if err != nil {
		logrus.WithField("username", user.Username).Errorf("cannot create token %s", err)
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
