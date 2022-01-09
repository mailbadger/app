package actions

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/templates"
	"github.com/mailbadger/app/utils"
	"github.com/mailbadger/app/validator"
)

func GetMe(c *gin.Context) {
	c.Header("X-CSRF-Token", csrf.Token(c.Request))
	c.JSON(http.StatusOK, middleware.GetUser(c))
}

func ChangePassword(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := middleware.GetUser(c)
		if u == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unable to fetch user.",
			})
			return
		}

		body := &params.ChangePassword{}
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

		err := bcrypt.CompareHashAndPassword([]byte(u.Password.String), []byte(body.Password))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "The password that you entered is incorrect.",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			logger.From(c).WithError(err).Error("change pass: unable to generate hash from password")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to update your password. Please try again.",
			})
			return
		}

		u.Password = sql.NullString{
			String: string(hashedPassword),
			Valid:  true,
		}

		err = storage.UpdateUser(u)
		if err != nil {
			logger.From(c).WithError(err).Error("change pass: Unable to update user's password")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to update your password. Please try again.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Your password was updated successfully.",
		})
	}
}

func PostForgotPassword(
	storage storage.Storage,
	emailSender emails.Sender,
	systemEmailSource string,
	appURL string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		body := &params.ForgotPassword{}
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

		// always send a success message
		c.JSON(http.StatusOK, gin.H{
			"message": "Email will be sent to you with the information on how to update your password.",
		})

		u, err := storage.GetUserByUsername(body.Email)
		if err != nil {
			logger.From(c).WithError(err).Warn("user tried to enter a non-existent email")
			return
		}

		tokenStr, err := utils.GenerateRandomString(32)
		if err != nil {
			logger.From(c).WithError(err).Error("forgot pass: unable to generate random string")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "We are unable to process this request, please try again.",
			})
			return
		}
		t := &entities.Token{
			UserID:    u.ID,
			Token:     tokenStr,
			Type:      entities.ForgotPasswordTokenType,
			ExpiresAt: time.Now().Add(time.Hour * 1),
		}
		err = storage.CreateToken(t)
		if err != nil {
			logger.From(c).WithError(err).Error("forgot pass: cannot create token")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "We are unable to process this request, please try again.",
			})
			return
		}

		go func(c *gin.Context) {
			err := sendForgotPasswordEmail(tokenStr, u.Username, emailSender, systemEmailSource, appURL)
			if err != nil {
				logger.From(c).WithError(err).Error("forgot pass: unable to send email")
			}
		}(c.Copy())
	}
}

func sendForgotPasswordEmail(
	token string,
	email string,
	sender emails.Sender,
	systemEmailSource string,
	appURL string,
) error {
	var html bytes.Buffer
	emailTmpls := templates.GetEmailTemplates()
	url := fmt.Sprintf("%s/forgot-password/%s", appURL, token)

	err := emailTmpls.ExecuteTemplate(&html, "forgot-password.html", map[string]string{
		"url": url,
	})
	if err != nil {
		return fmt.Errorf("send forgot password email: exec template: %w", err)
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
				Data:    aws.String(string("Reset your password")),
			},
		},
		Source: aws.String(fmt.Sprintf("%s <%s>", "Mailbadger.io", systemEmailSource)),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(email)},
		},
	})

	return err
}

func PutForgotPassword(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Param("token")

		t, err := storage.GetToken(tokenStr)
		if err != nil || t.Type != entities.ForgotPasswordTokenType {
			logger.From(c).WithError(err).Error("forgot pass: token is invalid")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to update your password. The token is invalid.",
			})
			return
		}

		body := &params.PutForgotPassword{}
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

		user, err := storage.GetUser(t.UserID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Unable to update your password. The user associated with the token is not found.",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.From(c).WithError(err).Error("Unable to generate hash from password.")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to update your password. Please try again.",
			})
			return
		}

		user.Password = sql.NullString{
			String: string(hashedPassword),
			Valid:  true,
		}

		err = storage.UpdateUser(user)
		if err != nil {
			logger.From(c).WithError(err).Error("forgot pass: unable to update user")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to update your password. Please try again.",
			})
			return
		}

		err = storage.DeleteToken(tokenStr)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{
				"token": tokenStr,
			}).WithError(err).Error("forgot pass: unable to delete token")
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Your password has been updated successfully.",
		})
	}
}

func PutVerifyEmail(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Param("token")

		t, err := storage.GetToken(tokenStr)
		if err != nil || t.Type != entities.VerifyEmailTokenType {
			logger.From(c).WithError(err).Error("verify email: token is invalid")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to verify your email. The token is invalid.",
			})
			return
		}

		user, err := storage.GetUser(t.UserID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Unable to verify your email, please try again.",
			})
			return
		}

		if !user.Verified {
			user.Verified = true
			err = storage.UpdateUser(user)
			if err != nil {
				logger.From(c).WithError(err).Error("verify email: unable to update user status to verified")
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Unable to verify your email. Please try again.",
				})
				return
			}
		}

		err = storage.DeleteToken(tokenStr)
		if err != nil {
			logger.From(c).WithFields(logrus.Fields{
				"token": tokenStr,
			}).WithError(err).Error("verify email: unable to delete token")
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Your email has been verified.",
		})
	}
}
