package actions

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v25/github"
	"github.com/google/uuid"
	fb "github.com/huandu/facebook"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	oauthfb "golang.org/x/oauth2/facebook"
	oauthgithub "golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gopkg.in/ezzarghili/recaptcha-go.v3"

	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/utils"
	"github.com/mailbadger/app/validator"
)

// PostAuthenticate authenticates a user with the given username and password.
func PostAuthenticate(c *gin.Context) {
	body := &params.PostAuthenticate{}
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters, please try again",
		})
		return
	}

	if err := validator.Validate(body); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := storage.GetActiveUserByUsername(c, body.Username)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			logger.From(c).WithError(err).Error("Unable to fetch active user by username.")
		}

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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Invalid credentials.",
		})
		return
	}

	sessID, err := utils.GenerateRandomString(32)
	if err != nil {
		logger.From(c).WithError(err).Error("Cannot create session id.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to create session id.",
		})
		return
	}

	err = persistSession(c, user.ID, sessID)
	if err != nil {
		logger.From(c).WithError(err).Error("Cannot persist session id.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to create session id.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
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

	body := &params.PostSignUp{}
	err := c.ShouldBind(body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid parameters, please try again.",
		})
		return
	}

	if err = validator.Validate(body); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	_, err = storage.GetUserByUsername(c, body.Email)
	if err == nil {
		logger.From(c).WithField("email", body.Email).Warn("Duplicate account.")
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Unable to create an account.",
		})
		return
	}

	secret := os.Getenv("RECAPTCHA_SECRET")
	if secret != "" {
		captcha, err := recaptcha.NewReCAPTCHA(secret, recaptcha.V2, 10*time.Second)
		if err != nil {
			logger.From(c).WithError(err).Error("Recaptcha initialize error.")
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unable to create an account. Captcha is invalid.",
			})
			return
		}

		err = captcha.Verify(body.TokenResponse)
		if err != nil {
			logger.From(c).WithField("username", body.Email).WithError(err).Infof("recaptcha invalid response.")
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unable to create an account.",
			})
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to generate hash from password.")
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Unable to create an account.",
		})
		return
	}

	uuid := uuid.New()

	user := &entities.User{
		Username: body.Email,
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
		logger.From(c).WithField("username", body.Email).WithError(err).Error("Unable to persist user.")
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
		tokenStr, err := utils.GenerateRandomString(32)
		if err != nil {
			logger.From(c).WithError(err).Error("Unable to generate random string.")
		}
		t := &entities.Token{
			UserID:    user.ID,
			Token:     tokenStr,
			Type:      entities.VerifyEmailTokenType,
			ExpiresAt: time.Now().AddDate(0, 0, 1),
		}
		err = storage.CreateToken(c, t)
		if err != nil {
			logger.From(c).WithError(err).Error("Cannot create token.")
		} else {
			go func(c *gin.Context) {
				err := sendVerifyEmail(tokenStr, user.Username, sender)
				if err != nil {
					logger.From(c).WithError(err).Error("Unable to send verification email.")
				}
			}(c.Copy())
		}
	} else {
		logger.From(c).WithError(err).Warn("Unable to create SES sender.")
	}

	sessID, err := utils.GenerateRandomString(32)
	if err != nil {
		logger.From(c).WithField("user_id", user.ID).WithError(err).Error("Cannot create session id.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to create session id.",
		})
		return
	}

	err = persistSession(c, user.ID, sessID)
	if err != nil {
		logger.From(c).WithField("user_id", user.ID).WithError(err).Error("Cannot persist session id.")
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
		logger.From(c).WithError(err).Error("Github: unable to generate random string.")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service unavailable.",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("state", state)
	err = session.Save()
	if err != nil {
		logger.From(c).WithError(err).Error("Github: unable to save session.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to save session, please try again.",
		})
		return
	}

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
		logger.From(c).WithError(err).Warn("Github: unable to exchange code for access token.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	tc := conf.Client(ctx, ghToken)
	client := github.NewClient(tc)

	ghUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		logger.From(c).WithError(err).Error("Github: get user error.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	completeCallback(c, ghUser.GetEmail(), "github", host)
}

// GetGoogleAuth redirects the user to the google oauth authorization page.
func GetGoogleAuth(c *gin.Context) {
	host := os.Getenv("APP_URL")

	state, err := utils.GenerateRandomString(12)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to generate random string.")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service unavailable.",
		})
		return
	}

	session := sessions.Default(c)

	session.Set("state", state)
	err = session.Save()
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to save session.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to save session, please try again.",
		})
		return
	}

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
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	session.Clear()

	if s != state {
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	ctx := context.Background()
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		logger.From(c).WithError(err).Warn("Google: exchange token error.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	oauth2Service, err := googleoauth2.NewService(ctx, option.WithTokenSource(conf.TokenSource(ctx, tok)))
	if err != nil {
		logger.From(c).WithError(err).Error("Google: unable to instantiate oauth2 service.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	userInfoSvc := googleoauth2.NewUserinfoV2MeService(oauth2Service)
	gUser, err := userInfoSvc.Get().Do()
	if err != nil {
		logger.From(c).WithError(err).Error("Google: fetch user error.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	completeCallback(c, gUser.Email, "google", host)
}

// GetFacebookAuth redirects the user to the facebook oauth authorization page.
func GetFacebookAuth(c *gin.Context) {
	host := os.Getenv("APP_URL")
	state, err := utils.GenerateRandomString(12)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to generate random string.")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service unavailable.",
		})
		return
	}

	session := sessions.Default(c)

	session.Set("state", state)
	err = session.Save()
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to save session.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to save session, please try again.",
		})
		return
	}

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
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	session.Clear()

	if s != state {
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	fbToken, err := conf.Exchange(ctx, code)
	if err != nil {
		logger.From(c).WithError(err).Warn("FB: exchange token error.")
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
		logger.From(c).WithError(err).Error("FB: unable to get user.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	email, ok := res["email"]
	if !ok {
		logger.From(c).WithField("resp", res).Warn("FB: response does not include email.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		logger.From(c).WithField("email", email).Error("FB: cannot convert email to string.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=server-error")
		return
	}

	completeCallback(c, emailStr, "facebook", host)
}

// PostLogout deletes the current user session.
func PostLogout(c *gin.Context) {
	session := sessions.Default(c)
	sessID := session.Get("sess_id")
	s, ok := sessID.(string)
	if ok {
		err := storage.DeleteSession(c, s)
		if err != nil {
			logger.From(c).WithError(err).Error("Unable to delete session.")
		}
	}

	session.Delete("sess_id")
	err := session.Save()
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to save session.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to save session, please try again.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "You have been successfully logged out.",
	})
}

func completeCallback(c *gin.Context, email, source, host string) {
	u, err := storage.GetUserByUsername(c, email)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			logger.From(c).WithError(err).Error("Social auth callback: unable to fetch user by username.")
			c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
			return
		}

		uuid := uuid.New()

		u = &entities.User{
			UUID:     uuid.String(),
			Username: email,
			Active:   true,
			Verified: true,
			Source:   source,
		}

		err = storage.CreateUser(c, u)
		if err != nil {
			logger.From(c).WithError(err).Error("Unable to create user.")
			c.Redirect(http.StatusPermanentRedirect, host+"/login?message=register-failed")
			return
		}
	}

	if !u.Active {
		logger.From(c).WithField("user_id", u.ID).Warn("Inactive user sign in.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	sessID, err := utils.GenerateRandomString(32)
	if err != nil {
		logger.From(c).WithField("user_id", u.ID).WithError(err).Error("Cannot create session id.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	err = persistSession(c, u.ID, sessID)
	if err != nil {
		logger.From(c).WithField("user_id", u.ID).WithError(err).Error("Cannot persist session.")
		c.Redirect(http.StatusPermanentRedirect, host+"/login?message=forbidden")
		return
	}

	c.Redirect(http.StatusPermanentRedirect, host+"/dashboard")
}

func sendVerifyEmail(token, email string, sender emails.Sender) error {
	url := os.Getenv("APP_URL") + "/verify-email/" + token

	_, err := sender.SendTemplatedEmail(&ses.SendTemplatedEmailInput{
		Template:     aws.String("VerifyEmail"),
		Source:       aws.String(os.Getenv("SYSTEM_EMAIL_SOURCE")),
		TemplateData: aws.String(fmt.Sprintf(`{"url": "%s"}`, url)),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(email)},
		},
	})

	return err
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
