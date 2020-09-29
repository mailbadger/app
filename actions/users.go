package actions

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
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
	"github.com/mailbadger/app/utils"
	"github.com/mailbadger/app/validator"
)

func GetMe(c *gin.Context) {
	c.Header("X-CSRF-Token", csrf.Token(c.Request))
	c.JSON(http.StatusOK, middleware.GetUser(c))
}

func ChangePassword(c *gin.Context) {
	u := middleware.GetUser(c)
	if u == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unable to fetch user.",
		})
		return
	}

	body := &params.ChangePassword{}
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

	err := bcrypt.CompareHashAndPassword([]byte(u.Password.String), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "The password that you entered is incorrect.",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to generate hash from password.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update your password. Please try again.",
		})
		return
	}

	u.Password = sql.NullString{
		String: string(hashedPassword),
		Valid:  true,
	}

	err = storage.UpdateUser(c, u)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to update user's password.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update your password. Please try again.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your password was updated successfully.",
	})
}

func PostForgotPassword(c *gin.Context) {
	body := &params.ForgotPassword{}
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

	u, err := storage.GetUserByUsername(c, body.Email)
	if err == nil {
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
				UserID:    u.ID,
				Token:     tokenStr,
				Type:      entities.ForgotPasswordTokenType,
				ExpiresAt: time.Now().Add(time.Hour * 1),
			}
			err = storage.CreateToken(c, t)
			if err != nil {
				logger.From(c).WithError(err).Error("Cannot create token.")
			} else {
				go func(c *gin.Context) {
					err := sendForgotPasswordEmail(tokenStr, u.Username, sender)
					if err != nil {
						logger.From(c).WithError(err).Error("Unable to send forgot pass email.")
					}
				}(c.Copy())
			}
		} else {
			logger.From(c).WithError(err).Warn("Unable to create SES sender.")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email will be sent to you with the information on how to update your password.",
	})
}

func sendForgotPasswordEmail(token, email string, sender emails.Sender) error {
	url := os.Getenv("APP_URL") + "/forgot-password/" + token

	_, err := sender.SendTemplatedEmail(&ses.SendTemplatedEmailInput{
		Template:     aws.String("ForgotPassword"),
		Source:       aws.String(os.Getenv("SYSTEM_EMAIL_SOURCE")),
		TemplateData: aws.String(fmt.Sprintf(`{"url": "%s"}`, url)),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(email)},
		},
	})

	return err
}


func PutForgotPassword(c *gin.Context) {
	tokenStr := c.Param("token")

	t, err := storage.GetToken(c, tokenStr)
	if err != nil || t.Type != entities.ForgotPasswordTokenType {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update your password. The token is invalid.",
		})
		return
	}

	body := &params.PutForgotPassword{}
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


	user, err := storage.GetUser(c, t.UserID)
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

	err = storage.UpdateUser(c, user)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to update user.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update your password. Please try again.",
		})
		return
	}

	err = storage.DeleteToken(c, tokenStr)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"token": tokenStr,
		}).WithError(err).Error("Unable to delete token.")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your password has been updated successfully.",
	})
}

func PutVerifyEmail(c *gin.Context) {
	tokenStr := c.Param("token")

	t, err := storage.GetToken(c, tokenStr)
	if err != nil || t.Type != entities.VerifyEmailTokenType {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to verify your email. The token is invalid.",
		})
		return
	}

	user, err := storage.GetUser(c, t.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Unable to verify your email. The user associated with the token is not found.",
		})
		return
	}

	if !user.Verified {
		user.Verified = true
		err = storage.UpdateUser(c, user)
		if err != nil {
			logger.From(c).WithError(err).Error("Unable to update user status to verified.")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to verify your email. Please try again.",
			})
			return
		}
	}

	err = storage.DeleteToken(c, tokenStr)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"token": tokenStr,
		}).WithError(err).Error("Unable to delete token.")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your email has been verified.",
	})
}
