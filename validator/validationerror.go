package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// genericMessage default message when casting validator error into validator.ValidationErrors
var genericMessage = "Invalid parameters, please try again"

// ValidationError represents custom validation error structure
// using this error when validating request body
type ValidationError struct {
	Message          string                     `json:"message"`
	Errors           map[string]string          `json:"errors,omitempty"`
	ValidationErrors validator.ValidationErrors `json:"-"`
}

// NewValidationError creates new ValidationError
func NewValidationError(ve validator.ValidationErrors) *ValidationError {
	validationError := &ValidationError{
		Message: genericMessage,
	}

	if ve != nil {
		validationError.ValidationErrors = ve
		validationError.FormatErrors()
	}

	return validationError
}

// FormatErrors creates key and message for each validation error
// be careful when using err.Param() only use it on tags with param value (ex: max=1)
func (q *ValidationError) FormatErrors() {
	q.Errors = make(map[string]string)

	for _, err := range q.ValidationErrors {
		switch err.ActualTag() {
		case "email":
			q.Errors[err.Field()] = "Invalid email format"
		case "required":
			q.Errors[err.Field()] = "This field is required"
		case "max":
			q.Errors[err.Field()] = "Max length allowed is " + err.Param()
		case "min":
			q.Errors[err.Field()] = "Must be at least " + err.Param() + " character long"
		case "alphanum":
			q.Errors[err.Field()] = "Only alphanumeric characters allowed"
		case "html":
			q.Errors[err.Field()] = "Content must be html"
		case tagAlphanumericHyphen:
			q.Errors[err.Field()] = "Must consist only of alphanumeric and hyphen characters"
		default:
			q.Errors[err.Field()] = "Validation failed on condition: " + err.ActualTag()
		}
	}
}

// Error overriding this func just to implement error interface
func (q ValidationError) Error() string {
	if q.ValidationErrors != nil {
		return fmt.Sprintf("validator: %s", q.ValidationErrors)
	}

	return fmt.Sprintf("validator: %s", q.Message)
}
