package validator

import (
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

type FieldError struct {
	Err validator.FieldError
}

func (q FieldError) String() string {
	var sb strings.Builder

	sb.WriteString("Validation failed on field '" + q.Err.Field() + "'")

	switch q.Err.ActualTag() {
	case "email":
		sb.WriteString(", wrong " + q.Err.ActualTag() + " format")
	case "required":
		sb.WriteString(", field is " + q.Err.ActualTag())
	case "max":
		sb.WriteString(", max length allowed: " + q.Err.Param())
	case "min":
		sb.WriteString(", min length allowed: " + q.Err.Param())
	default:
		sb.WriteString(", condition: " + q.Err.ActualTag())

	}

	return sb.String()
}
