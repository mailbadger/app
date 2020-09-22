package params

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	myvalidator "github.com/mailbadger/app/validator"
)

// RequestBody represents type for request body structures to simplify validation
type RequestBody interface {
	TrimSpaces()
}

// Validate is generic function for validation request body
func Validate(body RequestBody, v *validator.Validate) *myvalidator.FieldErrors {
	body.TrimSpaces()

	// Validate the instance
	if err := v.Struct(body); err != nil {
		logrus.Error(err)
		if fieldErrors, ok := err.(validator.ValidationErrors); ok {
			return &myvalidator.FieldErrors{
				Errors: fieldErrors,
			}
		}
		return new(myvalidator.FieldErrors)
	}

	return nil
}
