package actions

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-contrib/sessions"
	"github.com/jinzhu/gorm"
	"github.com/news-maily/api/emails"
	"github.com/news-maily/api/utils"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v25/github"
	"github.com/google/uuid"
	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/storage"
	"github.com/news-maily/api/utils/token"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	oauthgithub "golang.org/x/oauth2/github"
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
	Email         string `form:"email" valid:"email,required~Email is blank or in invalid format"`
	Password      string `form:"password" valid:"required"`
	TokenResponse string `form:"token_response" valid:"optional"`
}

func PostSignup(c *gin.Context) {
	enableSignup := os.Getenv("ENABLE_SIGNUP")
	if enableSignup == "" || enableSignup == "false" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Sign up is disabled.",
		})
		return
	}

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

	_, err = storage.GetUserByUsername(c, params.Email)
	if err == nil {
		logrus.WithField("username", params.Email).Errorf("duplicate account %s", err)
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
			logrus.WithField("username", params.Email).Errorf("recaptcha invalid response. %s", err)
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
		Username: params.Email,
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
		logrus.WithField("username", params.Email).Errorf("unable to persist user in db %s", err.Error())
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Unable to create an account.",
		})
		return
	}

	sender, err := emails.NewSesSender(
		os.Getenv("AWS_SES_ACCESS_KEY"),
		os.Getenv("AWS_SES_SECRET_KEY"),
		os.Getenv("AWS_SES_REGION"),
	)
	if err == nil {
		exp := time.Now().Add(time.Hour * 24).Unix()
		t := token.New(token.VerifyEmailToken, user.UUID)
		tokenStr, err := t.SignWithExp(os.Getenv("EMAILS_TOKEN_SECRET"), exp)
		if err != nil {
			logrus.WithError(err).Error("cannot create token")
		} else {
			go sendVerifyEmail(tokenStr, user.Username, sender)
		}
	} else {
		logrus.WithError(err).Error("unable to instantiate ses sender")
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

func GetGithubAuth(c *gin.Context) {
	state, err := utils.GenerateRandomString(12)
	if err != nil {
		logrus.WithError(err).Error("unable to generate random string")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service unavailable.",
		})
		return
	}

	session := sessions.Default(c)

	session.Set("state", state)
	session.Save()

	url := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=user:email&state=%s",
		os.Getenv("GITHUB_CLIENT_ID"),
		state,
	)

	c.Redirect(http.StatusPermanentRedirect, url)
}

func GithubCallback(c *gin.Context) {
	host := os.Getenv("DOMAIN_URL")

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user:email"},
		Endpoint:     oauthgithub.Endpoint,
	}

	code := c.Query("code")
	state := c.Query("state")

	session := sessions.Default(c)

	sessState := session.Get("state")
	s, ok := sessState.(string)
	if !ok {
		logrus.Error("unable to fetch state from session")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	session.Clear()

	if s != state {
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	ghToken, err := conf.Exchange(ctx, code)
	if err != nil {
		logrus.WithError(err).Error("exchange token error")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	tc := conf.Client(ctx, ghToken)
	client := github.NewClient(tc)

	ghUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		logrus.WithError(err).Error("fetch user error")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	u, err := storage.GetUserByUsername(c, ghUser.GetEmail())
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			uuid, err := uuid.NewRandom()
			if err != nil {
				logrus.WithError(err).Error("unable to generate random uuid")
				c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
				return
			}

			u = &entities.User{
				UUID:     uuid.String(),
				Username: ghUser.GetEmail(),
				Active:   true,
				Source:   "github",
			}

			err = storage.CreateUser(c, u)
			if err != nil {
				logrus.WithError(err).Error("github register failed")
				c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
				return
			}

			sender, err := emails.NewSesSender(
				os.Getenv("AWS_SES_ACCESS_KEY"),
				os.Getenv("AWS_SES_SECRET_KEY"),
				os.Getenv("AWS_SES_REGION"),
			)
			if err == nil {
				exp := time.Now().Add(time.Hour * 24).Unix()
				t := token.New(token.VerifyEmailToken, u.UUID)
				tokenStr, err := t.SignWithExp(os.Getenv("EMAILS_TOKEN_SECRET"), exp)
				if err != nil {
					logrus.WithError(err).Error("cannot create token")
				} else {
					go sendVerifyEmail(tokenStr, u.Username, sender)
				}
			} else {
				logrus.WithError(err).Error("unable to instantiate ses sender")
			}
		} else {
			logrus.WithError(err).Error("github social auth")
			c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
			return
		}
	}

	if !u.Active {
		logrus.WithField("user_id", u.ID).Warn("inactive user sign in via github")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	exp := time.Now().Add(time.Hour * 72).Unix()
	t := token.New(token.SessionToken, u.Username)
	tokenStr, err := t.SignWithExp(os.Getenv("AUTH_SECRET"), exp)
	if err != nil {
		logrus.WithField("user_id", u.ID).WithError(err).Error("cannot create token")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	c.Redirect(http.StatusPermanentRedirect, host+"/login/callback?t="+tokenStr+"&exp="+strconv.Itoa(int(exp)))
}

func sendVerifyEmail(token, email string, sender emails.Sender) {
	url := os.Getenv("DOMAIN_URL") + "/verify-email/" + token

	_, err := sender.SendTemplatedEmail(&ses.SendTemplatedEmailInput{
		Template:     aws.String("VerifyEmail"),
		Source:       aws.String(os.Getenv("SYSTEM_EMAIL_SOURCE")),
		TemplateData: aws.String(fmt.Sprintf(`{"url": "%s"}`, url)),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(email)},
		},
	})

	if err != nil {
		logrus.WithError(err).Error("email verification - send email failure")
	}
}
