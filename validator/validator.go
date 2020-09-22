package validator

import (
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	once sync.Once
	validatorTagName = "validate"
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

func initValidator() {
	MBValidator = validator.New()
	MBValidator.SetTagName(validatorTagName)
}

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &DefaultValidator{}

func (v *DefaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}

	return nil
}

func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding") // Print JSON name on validator.FieldError.Field()
		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
			return name
		})
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
