package actions

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

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
	"gorm.io/gorm"

	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/session"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/templates"
	"github.com/mailbadger/app/utils"
	"github.com/mailbadger/app/validator"
)

// PostAuthenticate authenticates a user with the given username and password.
func PostAuthenticate(storage storage.Storage, sess session.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		body := &params.PostAuthenticate{}
		if err := c.ShouldBindJSON(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid parameters, please try again",
			})
			return
		}

		if err := validator.Validate(body); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		user, err := storage.GetActiveUserByUsername(body.Username)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
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

		err = sess.CreateUserSession(c, user.ID)
		if err != nil {
			logger.From(c).WithError(err).Error("Cannot persist session id.")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to create session.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

// PostSignup validates and creates a user account by the given
// user parameters. The handler also sends a verification email.
func PostSignup(
	storage storage.Storage,
	sess session.Session,
	emailSender emails.Sender,
	enableSignup bool,
	verifyEmail bool,
	recaptchaSecret string,
	systemEmailSource string,
	appURL string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enableSignup {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Sign up is disabled.",
			})
			return
		}

		body := &params.PostSignUp{}
		err := c.ShouldBindJSON(body)
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

		_, err = storage.GetUserByUsername(body.Email)
		if err == nil {
			logger.From(c).WithField("email", body.Email).Warn("Duplicate account.")
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unable to create an account.",
			})
			return
		}

		if recaptchaSecret != "" {
			captcha, err := recaptcha.NewReCAPTCHA(recaptchaSecret, recaptcha.V2, 10*time.Second)
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

		b, err := storage.GetBoundariesByType(entities.BoundaryTypeFree)
		if err != nil {
			logger.From(c).WithError(err).Error("signup: unable to fetch boundary")
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unable to create an account.",
			})
			return
		}

		r, err := storage.GetRole(entities.AdminRole)
		if err != nil {
			logger.From(c).WithError(err).Error("signup: unable to fetch admin role")
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unable to create an account.",
			})
			return
		}

		uuid := uuid.NewString()

		user := &entities.User{
			Username: body.Email,
			UUID:     uuid,
			Password: sql.NullString{
				String: string(hashedPassword),
				Valid:  true,
			},
			Active:     true,
			Verified:   false,
			Boundaries: b,
			Roles:      []entities.Role{*r},
			Source:     "mailbadger.io",
		}

		err = storage.CreateUser(user)
		if err != nil {
			logger.From(c).WithField("username", body.Email).WithError(err).Error("Unable to persist user.")
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Unable to create an account.",
			})
			return
		}

		if verifyEmail {
			go func(c *gin.Context) {
				err := sendVerifyEmail(storage, emailSender, user, systemEmailSource, appURL)
				if err != nil {
					logger.From(c).WithError(err).Error("Unable to send verification email.")
				}
			}(c.Copy())
		}

		err = sess.CreateUserSession(c, user.ID)
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
}

// GetGithubAuth redirects the user to the github oauth authorization page.
func GetGithubAuth(clientID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := utils.GenerateRandomString(12)
		if err != nil {
			logger.From(c).WithError(err).Error("Github: unable to generate random string")
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"message": "Service unavailable.",
			})
			return
		}

		session := sessions.Default(c)
		session.Set("state", state)
		err = session.Save()
		if err != nil {
			logger.From(c).WithError(err).Error("Github: unable to save session")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to save session, please try again.",
			})
			return
		}

		url := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=user:email&state=%s", clientID, state)

		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// GithubCallback fetches the github user by the given access code and creates a new session.
// The callback creates a new user if the email does not exist in our system.
// If the sign in is successful it redirects the user to the dashboard, if it fails we redirect the user
// to the login screen with an error message.
func GithubCallback(
	storage storage.Storage,
	sess session.Session,
	clientID string,
	clientSecret string,
	appURL string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"user:email"},
			Endpoint:     oauthgithub.Endpoint,
		}

		code := c.Query("code")
		state := c.Query("state")

		session := sessions.Default(c)

		sessState := session.Get("state")
		s, ok := sessState.(string)
		if !ok {
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		session.Clear()

		if s != state {
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		ghToken, err := conf.Exchange(c, code)
		if err != nil {
			logger.From(c).WithError(err).Error("Github: unable to exchange code for access token")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		tc := conf.Client(c, ghToken)
		client := github.NewClient(tc)

		ghUser, _, err := client.Users.Get(c, "")
		if err != nil {
			logger.From(c).WithError(err).Error("Github: get user error")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		completeCallback(c, storage, sess, ghUser.GetEmail(), "github", appURL)
	}
}

// GetGoogleAuth redirects the user to the google oauth authorization page.
func GetGoogleAuth(clientID, clientSecret, appURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := utils.GenerateRandomString(12)
		if err != nil {
			logger.From(c).WithError(err).Error("Google: unable to generate random string")
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"message": "Service unavailable.",
			})
			return
		}

		session := sessions.Default(c)

		session.Set("state", state)
		err = session.Save()
		if err != nil {
			logger.From(c).WithError(err).Error("Google: unable to save session")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to save session, please try again.",
			})
			return
		}

		conf := &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  appURL + "/api/auth/google/callback",
			Scopes: []string{
				googleoauth2.UserinfoEmailScope,
			},
			Endpoint: google.Endpoint,
		}

		url := conf.AuthCodeURL(state)

		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// GoogleCallback fetches the google user by the given access code and creates a new session.
// The callback creates a new user if the email does not exist in our system.
// If the sign in is successful it redirects the user to the dashboard, if it fails we redirect the user
// to the login screen with an error message.
func GoogleCallback(
	storage storage.Storage,
	sess session.Session,
	clientID string,
	clientSecret string,
	appURL string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  appURL + "/api/auth/google/callback",
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
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		session.Clear()

		if s != state {
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		tok, err := conf.Exchange(c, code)
		if err != nil {
			logger.From(c).WithError(err).Error("Google: exchange token error")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		oauth2Service, err := googleoauth2.NewService(c, option.WithTokenSource(conf.TokenSource(c, tok)))
		if err != nil {
			logger.From(c).WithError(err).Error("Google: unable to instantiate oauth2 service")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		userInfoSvc := googleoauth2.NewUserinfoV2MeService(oauth2Service)
		gUser, err := userInfoSvc.Get().Do()
		if err != nil {
			logger.From(c).WithError(err).Error("Google: fetch user error")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		completeCallback(c, storage, sess, gUser.Email, "google", appURL)
	}
}

// GetFacebookAuth redirects the user to the facebook oauth authorization page.
func GetFacebookAuth(clientID, appURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := utils.GenerateRandomString(12)
		if err != nil {
			logger.From(c).WithError(err).Error("Facebook: unable to generate random string")
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"message": "Service unavailable.",
			})
			return
		}

		session := sessions.Default(c)

		session.Set("state", state)
		err = session.Save()
		if err != nil {
			logger.From(c).WithError(err).Error("Facebook: unable to save session")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to save session, please try again.",
			})
			return
		}

		url := fmt.Sprintf("https://www.facebook.com/v3.3/dialog/oauth?client_id=%s&scope=email&redirect_uri=%s&state=%s",
			clientID,
			appURL+"/api/auth/facebook/callback",
			state,
		)

		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// FacebookCallback fetches the facebook user by the given access code and creates a new session.
// The callback creates a new user if the email does not exist in our system.
// If the sign in is successful it redirects the user to the dashboard, if it fails we redirect the user
// to the login screen with an error message.
func FacebookCallback(
	storage storage.Storage,
	sess session.Session,
	clientID string,
	clientSecret string,
	appURL string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")

		session := sessions.Default(c)

		sessState := session.Get("state")
		s, ok := sessState.(string)
		if !ok {
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		session.Clear()

		if s != state {
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		conf := &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"email"},
			Endpoint:     oauthfb.Endpoint,
		}
		fbToken, err := conf.Exchange(c, code)
		if err != nil {
			logger.From(c).WithError(err).Error("Facebook: exchange token error")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		tc := conf.Client(c, fbToken)
		fbsess := &fb.Session{
			HttpClient: tc,
			Version:    "v3.3",
		}

		res, err := fbsess.Get("/me", nil)
		if err != nil {
			logger.From(c).WithError(err).Error("Facebook: unable to get user")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		email, ok := res["email"]
		if !ok {
			logger.From(c).WithField("resp", res).Error("Facebook: response does not include email")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		emailStr, ok := email.(string)
		if !ok {
			logger.From(c).WithField("email", email).Error("Facebook: cannot convert email to string")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=server-error")
			return
		}

		completeCallback(c, storage, sess, emailStr, "facebook", appURL)
	}
}

// PostLogout deletes the current user session.
func PostLogout(sess session.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := sess.DeleteUserSession(c)
		if err != nil {
			logger.From(c).WithError(err).Error("logout: unable to delete session.")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "We are unable to process the request, please try again.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "You have been successfully logged out.",
		})
	}
}

func completeCallback(
	c *gin.Context,
	storage storage.Storage,
	sess session.Session,
	email string,
	source string,
	appURL string,
) {
	u, err := storage.GetUserByUsername(email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.From(c).WithError(err).Error("social auth callback: unable to fetch user by username")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=register-failed")
			return
		}

		b, err := storage.GetBoundariesByType(entities.BoundaryTypeFree)
		if err != nil {
			logger.From(c).WithError(err).Error("social auth callback: unable to fetch boundary")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=register-failed")
			return
		}

		r, err := storage.GetRole(entities.AdminRole)
		if err != nil {
			logger.From(c).WithError(err).Error("social auth callback: unable to fetch admin role")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=register-failed")
			return
		}

		uuid := uuid.NewString()

		u = &entities.User{
			UUID:       uuid,
			Username:   email,
			Active:     true,
			Verified:   true,
			Source:     source,
			Boundaries: b,
			Roles:      []entities.Role{*r},
		}

		err = storage.CreateUser(u)
		if err != nil {
			logger.From(c).WithError(err).Error("social auth callback: unable to create user")
			c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=register-failed")
			return
		}
	}

	if !u.Active {
		logger.From(c).WithField("user_id", u.ID).Warn("social auth callback: inactive user sign in")
		c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=forbidden")
		return
	}

	err = sess.CreateUserSession(c, u.ID)
	if err != nil {
		logger.From(c).WithField("user_id", u.ID).WithError(err).Error("Cannot persist session.")
		c.Redirect(http.StatusTemporaryRedirect, appURL+"/login?message=forbidden")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, appURL+"/dashboard")
}

func sendVerifyEmail(
	storage storage.Storage,
	sender emails.Sender,
	u *entities.User,
	systemEmailSource string,
	appURL string,
) error {
	token, err := utils.GenerateRandomString(32)
	if err != nil {
		return fmt.Errorf("send verify email: gen token: %w", err)
	}

	t := &entities.Token{
		UserID:    u.ID,
		Token:     token,
		Type:      entities.VerifyEmailTokenType,
		ExpiresAt: time.Now().AddDate(0, 0, 1),
	}

	err = storage.CreateToken(t)
	if err != nil {
		return fmt.Errorf("send verify email: create token: %w", err)
	}

	var html bytes.Buffer
	emailTmpls := templates.GetEmailTemplates()
	url := fmt.Sprintf("%s/verify-email/%s", appURL, token)

	err = emailTmpls.ExecuteTemplate(&html, "verify-email.html", map[string]string{
		"url": url,
	})
	if err != nil {
		return fmt.Errorf("send verify email: exec template: %w", err)
	}

	charset := aws.String("UTF-8")
	_, err = sender.SendEmail(&ses.SendEmailInput{
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: charset,
					Data:    aws.String(html.String()),
				},
			},
			Subject: &ses.Content{
				Charset: charset,
				Data:    aws.String(string("Verify your email address")),
			},
		},
		Source: aws.String(fmt.Sprintf("%s <%s>", "Mailbadger.io", systemEmailSource)),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(u.Username)},
		},
	})
	if err != nil {
		return fmt.Errorf("send verify email: %w", err)
	}

	return nil
}
