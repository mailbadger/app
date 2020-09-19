package actions

import (
	"net/http"

	mbvalidator "github.com/mailbadger/app/validator"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

func AbortWithError(c *gin.Context, err error) {
	for _, fieldErr := range err.(validator.ValidationErrors) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": mbvalidator.FieldError{Err: fieldErr}.String(),
		})
	}
}
