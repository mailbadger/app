package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	mbvalidator "github.com/mailbadger/app/validator"
)

func AbortWithError(c *gin.Context, err error) {
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fieldErrors.Error(),
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"message": mbvalidator.ErrGeneric.Error(),
	})
}
