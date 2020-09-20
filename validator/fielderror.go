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

	if q.Err.ActualTag() == "email" {
		sb.WriteString("Validation failed on field '" + q.Err.Field() + "'")
		sb.WriteString(", wrong " + q.Err.ActualTag() + " format")
		return sb.String()
	}

	sb.WriteString("Validation failed on field '" + q.Err.Field() + "'")
	sb.WriteString(", condition: " + q.Err.ActualTag())

	// Print condition parameters, e.g. max=191 -> max { 191 }
	if q.Err.Param() != "" {
		sb.WriteString(" " + q.Err.Param())
	}

	return sb.String()
}
