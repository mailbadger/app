package validator

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"

	"github.com/mailbadger/app/entities/params"
)

var (
	once                           sync.Once
	validatorTagName               = "validate"
	regexPatternAlphanumericHyphen = "^[\\w-]*$"
	tagAlphanumericHyphen          = "alphanumhyphen"
)

// MBValidator global validator
var MBValidator *validator.Validate

// Validator is a constructor for our validator
// if our validator is once created it returns it else it creates it
func Validator() *validator.Validate {
	once.Do(func() {
		initValidator()
	})

	return MBValidator
}

// Validate is generic function for validation request body
func Validate(body params.RequestBody) error {
	body.TrimSpaces()

	// Validate the instance
	if err := Validator().Struct(body); err != nil {
		if fieldErrors, ok := err.(validator.ValidationErrors); ok {
			return NewValidationError(fieldErrors)
		}
		return NewValidationError(nil)
	}

	return nil
}

func initValidator() {
	MBValidator = validator.New()
	MBValidator.SetTagName(validatorTagName)
	MBValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		return name
	})

	err := MBValidator.RegisterValidation(tagAlphanumericHyphen, func(fl validator.FieldLevel) bool {
		matched, _ := regexp.MatchString(regexPatternAlphanumericHyphen, fl.Field().String())
		return matched
	})
	if err != nil {
		// if register validation fails panic
		panic(err)
	}
}
