package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// ErrGeneric default error when casting validator error into validator.ValidationErrors
var ErrGeneric = errors.New("Invalid parameters, please try again")

type FieldErrors struct {
	Errors      validator.ValidationErrors
}

// FormatErrors creates key and message for each validation error
func (q FieldErrors) FormatErrors() map[string]string {
	errMessages := make(map[string]string)

	for _, err := range q.Errors {
		switch err.ActualTag() {
		case "email":
			errMessages[err.Field()] = "Invalid email format"
		case "required":
			errMessages[err.Field()] = "This field is required"
		case "max":
			errMessages[err.Field()] = "Max length allowed is " + err.Param()
		case "min":
			errMessages[err.Field()] = "Must be at least " + err.Param() + " character long"
		default:
			errMessages[err.Field()] = "Validation failed on condition: " + err.ActualTag()
		}
	}

	return errMessages
}

// Override this func just to implement error interface
func (q FieldErrors) Error() string {
	return "Invalid parameters, please try again"
}
