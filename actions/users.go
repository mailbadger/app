package actions

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	"github.com/news-maily/app/emails"
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/logger"
	"github.com/news-maily/app/routes/middleware"
	"github.com/news-maily/app/storage"
	"github.com/news-maily/app/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func GetMe(c *gin.Context) {
	c.Header("X-CSRF-Token", csrf.Token(c.Request))
	c.JSON(http.StatusOK, middleware.GetUser(c))
}

type changePassParams struct {
	Password    string `form:"password" valid:"required"`
	NewPassword string `form:"new_password" valid:"required"`
}

func ChangePassword(c *gin.Context) {
	u := middleware.GetUser(c)
	if u == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unable to fetch user.",
		})
		return
	}

	params := &changePassParams{}
	err := c.Bind(params)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to bind params.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid parameters, please try again.",
		})
		return
	}

	v, err := valid.ValidateStruct(params)
	if !v {
		msg := "Unable to change password, invalid request parameters."
		if err != nil {
			msg = err.Error()
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": msg,
		})
		return
	}

	if len(params.NewPassword) < 8 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"new_password": "The new password must be atleast 8 characters.",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password.String), []byte(params.Password))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "The password that you entered is incorrect.",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.NewPassword), bcrypt.DefaultCost)
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

type forgotPassParams struct {
	Email string `form:"email" valid:"email"`
}

func PostForgotPassword(c *gin.Context) {
	params := &forgotPassParams{}
	err := c.Bind(params)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to bind params.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid parameters, please try again.",
		})
		return
	}

	v, err := valid.ValidateStruct(params)
	if !v {
		emailError := valid.ErrorByField(err, "Email")
		if emailError == "" {
			emailError = "Email must be in valid format."
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": emailError,
		})
		return
	}

	u, err := storage.GetUserByUsername(c, params.Email)
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

type putForgotPassParams struct {
	Password string `form:"password" valid:"required"`
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

	params := &putForgotPassParams{}
	err = c.Bind(params)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to bind params.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid parameters, please try again.",
		})
		return
	}

	v, err := valid.ValidateStruct(params)
	if !v {
		passError := valid.ErrorByField(err, "Password")
		if passError == "" {
			passError = "The password must not be empty."
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": passError,
		})
		return
	}

	if len(params.Password) < 8 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"password": "The new password must be atleast 8 characters.",
		})
		return
	}

	user, err := storage.GetUser(c, t.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Unable to update your password. The user associated with the token is not found.",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
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
