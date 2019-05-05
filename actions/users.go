package actions

import (
	"database/sql"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/routes/middleware"
	"github.com/news-maily/api/storage"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func GetMe(c *gin.Context) {
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
	c.Bind(params)

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
		logrus.Errorf("Invalid credentials. %s", err)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "The password that you entered is incorrect.",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user": u.ID,
		}).Println(err)
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
		logrus.WithFields(logrus.Fields{
			"user": u.ID,
		}).Println(err)
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

}
