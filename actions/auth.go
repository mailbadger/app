package actions

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v25/github"
	"github.com/google/uuid"
	fb "github.com/huandu/facebook"
	"github.com/jinzhu/gorm"
	"github.com/news-maily/app/emails"
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/storage"
	"github.com/news-maily/app/utils"
	"github.com/news-maily/app/utils/token"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	oauthfb "golang.org/x/oauth2/facebook"
	oauthgithub "golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gopkg.in/ezzarghili/recaptcha-go.v3"
)

// PostAuthenticate authenticates a user with the given username and password.
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
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Invalid credentials.",
		})
		return
	}

	sessID, err := utils.GenerateRandomString(32)
	if err != nil {
		logrus.WithField("user_id", user.ID).WithError(err).Error("Cannot create session id.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to create session id.",
		})
		return
	}

	err = persistSession(c, user.ID, sessID)
	if err != nil {
		logrus.WithField("user_id", user.ID).WithError(err).Error("Cannot persist session id.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to create session id.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

type signupParams struct {
	Email         string `form:"email" valid:"email,required~Email is blank or in invalid format"`
	Password      string `form:"password" valid:"required"`
	TokenResponse string `form:"token_response" valid:"optional"`
}

// PostSignup validates and creates a user account by the given
// user parameters. The handler also sends a verification email
func PostSignup(c *gin.Context) {
	enableSignup, _ := strconv.ParseBool(os.Getenv("ENABLE_SIGNUP"))
	if !enableSignup {
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
		logrus.WithField("username", params.Email).Error("duplicate account")
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Unable to create an account.",
		})
		return
	}

	secret := os.Getenv("RECAPTCHA_SECRET")
	if secret != "" {
		captcha, err := recaptcha.NewReCAPTCHA(secret, recaptcha.V2, 10*time.Second)
		if err != nil {
			logrus.WithError(err).Error("recaptcha initialize error")
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unable to create an account. Captcha is invalid.",
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
		logrus.WithError(err).Error("unable to generate random uuid")
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Unable to create an account.",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithError(err).Error("unable to generate hash from password")
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
		logrus.WithField("username", params.Email).WithError(err).Error("unable to persist user in db")
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

	sessID, err := utils.GenerateRandomString(32)
	if err != nil {
		logrus.WithField("user_id", user.ID).WithError(err).Error("Cannot create session id.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to create session id.",
		})
		return
	}

	err = persistSession(c, user.ID, sessID)
	if err != nil {
		logrus.WithField("user_id", user.ID).WithError(err).Error("Cannot persist session id.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to create session id.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// GetGithubAuth redirects the user to the github oauth authorization page.
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

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GithubCallback fetches the github user by the given access code and creates a new session.
// The callback creates a new user if the email does not exist in our system.
// If the sign in is successful it redirects the user to the dashboard, if it fails we redirect the user
// to the login screen with an error message.
func GithubCallback(c *gin.Context) {
	host := os.Getenv("APP_URL")

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
		logrus.WithFields(logrus.Fields{
			"state":          s,
			"callback_state": state,
		}).Error("state mismatch")
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
		if !gorm.IsRecordNotFoundError(err) {
			logrus.WithError(err).Error("github social auth")
			c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
			return
		}

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
			Verified: true,
			Source:   "github",
		}

		err = storage.CreateUser(c, u)
		if err != nil {
			logrus.WithError(err).Error("github register failed")
			c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
			return
		}

	}

	if !u.Active {
		logrus.WithField("user_id", u.ID).Warn("inactive user sign in via github")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	sessID, err := utils.GenerateRandomString(32)
	if err != nil {
		logrus.WithField("user_id", u.ID).WithError(err).Error("Cannot create session id.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	err = persistSession(c, u.ID, sessID)
	if err != nil {
		logrus.WithField("user_id", u.ID).WithError(err).Error("Cannot persist session.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	c.Redirect(http.StatusPermanentRedirect, host+"/dashboard")
}

// GetGoogleAuth redirects the user to the google oauth authorization page.
func GetGoogleAuth(c *gin.Context) {
	host := os.Getenv("APP_URL")

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

	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  host + "/api/auth/google/callback",
		Scopes: []string{
			googleoauth2.UserinfoEmailScope,
		},
		Endpoint: google.Endpoint,
	}

	url := conf.AuthCodeURL(state)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback fetches the google user by the given access code and creates a new session.
// The callback creates a new user if the email does not exist in our system.
// If the sign in is successful it redirects the user to the dashboard, if it fails we redirect the user
// to the login screen with an error message.
func GoogleCallback(c *gin.Context) {
	host := os.Getenv("APP_URL")
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  host + "/api/auth/google/callback",
		Scopes: []string{
			googleoauth2.UserinfoEmailScope,
		},
		Endpoint: google.Endpoint,
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
		logrus.WithFields(logrus.Fields{
			"state":          s,
			"callback_state": state,
		}).Error("state mismatch")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	ctx := context.Background()
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		logrus.WithError(err).Error("exchange token error")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	oauth2Service, err := googleoauth2.NewService(ctx, option.WithTokenSource(conf.TokenSource(ctx, tok)))
	if err != nil {
		logrus.WithError(err).Error("unable to instantiate oauth2 service")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	userInfoSvc := googleoauth2.NewUserinfoV2MeService(oauth2Service)
	gUser, err := userInfoSvc.Get().Do()
	if err != nil {
		logrus.WithError(err).Error("fetch user error")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	u, err := storage.GetUserByUsername(c, gUser.Email)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			logrus.WithError(err).Error("google social auth")
			c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
			return
		}
		uuid, err := uuid.NewRandom()
		if err != nil {
			logrus.WithError(err).Error("unable to generate random uuid")
			c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
			return
		}

		u = &entities.User{
			UUID:     uuid.String(),
			Username: gUser.Email,
			Active:   true,
			Verified: true,
			Source:   "google",
		}

		err = storage.CreateUser(c, u)
		if err != nil {
			logrus.WithError(err).Error("google register failed")
			c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
			return
		}
	}

	if !u.Active {
		logrus.WithField("user_id", u.ID).Warn("inactive user sign in via google")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	sessID, err := utils.GenerateRandomString(32)
	if err != nil {
		logrus.WithField("user_id", u.ID).WithError(err).Error("Cannot create session id.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	err = persistSession(c, u.ID, sessID)
	if err != nil {
		logrus.WithField("user_id", u.ID).WithError(err).Error("Cannot persist session.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	c.Redirect(http.StatusPermanentRedirect, host+"/dashboard")
}

// GetFacebookAuth redirects the user to the facebook oauth authorization page.
func GetFacebookAuth(c *gin.Context) {
	host := os.Getenv("APP_URL")
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

	url := fmt.Sprintf("https://www.facebook.com/v3.3/dialog/oauth?client_id=%s&scope=email&redirect_uri=%s&state=%s",
		os.Getenv("FACEBOOK_CLIENT_ID"),
		host+"/api/auth/facebook/callback",
		state,
	)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// FacebookCallback fetches the facebook user by the given access code and creates a new session.
// The callback creates a new user if the email does not exist in our system.
// If the sign in is successful it redirects the user to the dashboard, if it fails we redirect the user
// to the login screen with an error message.
func FacebookCallback(c *gin.Context) {
	host := os.Getenv("APP_URL")

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
		ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		Scopes:       []string{"email"},
		Endpoint:     oauthfb.Endpoint,
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
		logrus.WithFields(logrus.Fields{
			"state":          s,
			"callback_state": state,
		}).Error("state mismatch")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	fbToken, err := conf.Exchange(ctx, code)
	if err != nil {
		logrus.WithError(err).Error("exchange token error")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	tc := conf.Client(ctx, fbToken)
	sess := &fb.Session{
		HttpClient: tc,
		Version:    "v3.3",
	}

	res, err := sess.Get("/me", nil)
	if err != nil {
		logrus.WithError(err).Error("fb client error")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	email, ok := res["email"]
	if !ok {
		logrus.WithField("resp", res).Error("fb response does not include email")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		logrus.WithField("email", email).Error("cannot convert email to string")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	u, err := storage.GetUserByUsername(c, emailStr)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			logrus.WithError(err).Error("facebook social auth")
			c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
			return
		}

		uuid, err := uuid.NewRandom()
		if err != nil {
			logrus.WithError(err).Error("unable to generate random uuid")
			c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
			return
		}

		u = &entities.User{
			UUID:     uuid.String(),
			Username: emailStr,
			Active:   true,
			Verified: true,
			Source:   "facebook",
		}

		err = storage.CreateUser(c, u)
		if err != nil {
			logrus.WithError(err).Error("facebook register failed")
			c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
			return
		}
	}

	if !u.Active {
		logrus.WithField("user_id", u.ID).Warn("inactive user sign in via facebook")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	sessID, err := utils.GenerateRandomString(32)
	if err != nil {
		logrus.WithField("user_id", u.ID).WithError(err).Error("Cannot create session id.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	err = persistSession(c, u.ID, sessID)
	if err != nil {
		logrus.WithField("user_id", u.ID).WithError(err).Error("Cannot persist session.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	c.Redirect(http.StatusPermanentRedirect, host+"/dashboard")
}

// PostLogout deletes the current user session.
func PostLogout(c *gin.Context) {
	session := sessions.Default(c)
	sessID := session.Get("sess_id")
	s, ok := sessID.(string)
	if ok {
		err := storage.DeleteSession(c, s)
		if err != nil {
			logrus.WithError(err).Error("Unable to delete session.")
		}
	}

	session.Delete("sess_id")
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "You have been successfully logged out.",
	})
}

func sendVerifyEmail(token, email string, sender emails.Sender) {
	url := os.Getenv("APP_URL") + "/verify-email/" + token

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

func persistSession(c *gin.Context, userID int64, sessID string) error {
	err := storage.CreateSession(c, &entities.Session{
		UserID:    userID,
		SessionID: sessID,
	})
	if err != nil {
		return err
	}

	session := sessions.Default(c)
	exp := time.Now().Add(time.Hour*72).Unix() - time.Now().Unix()
	secureCookie, _ := strconv.ParseBool(os.Getenv("SECURE_COOKIE"))
	session.Options(sessions.Options{
		HttpOnly: true,
		MaxAge:   int(exp),
		Secure:   secureCookie,
		Path:     "/api",
	})
	session.Set("sess_id", sessID)

	return session.Save()
}
