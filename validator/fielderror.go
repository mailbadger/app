package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// GenericValidationError default error when casting validator error into validator.ValidationErrors
var GenericValidationError = errors.New("Invalid parameters, please try again")

type FieldErrors struct {
	Errors validator.ValidationErrors
}

func (q FieldErrors) Error() string {
	var sb strings.Builder

	for _, err := range q.Errors {
		sb.WriteString("Validation failed on field '" + err.Field() + "'")

		switch err.ActualTag() {
		case "email":
			sb.WriteString(", wrong email format")
		case "required":
			sb.WriteString(", field is required")
		case "max":
			sb.WriteString(", max length allowed: " + err.Param())
		case "min":
			sb.WriteString(", min length allowed: " + err.Param())
		default:
			sb.WriteString(", condition: " + err.ActualTag())

		}
	}

	return sb.String()
}
