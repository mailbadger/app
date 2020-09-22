package params

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	myvalidator "github.com/mailbadger/app/validator"
)

type PostSESKeys struct {
	AccessKey string `form:"access_key" validate:"required,alphanum"`
	SecretKey string `form:"secret_key" validate:"required"`
	Region    string `form:"region" validate:"required"`
}

func (p *PostSESKeys) trimSpaces() {
	p.AccessKey = strings.TrimSpace(p.AccessKey)
	p.SecretKey = strings.TrimSpace(p.SecretKey)
	p.Region = strings.TrimSpace(p.Region)
}

func (p *PostSESKeys) Validate(v *validator.Validate) error {
	p.trimSpaces()

	// Validate the instance
	if err := v.Struct(p); err != nil {
		logrus.Error(err)
		if fieldErrors, ok := err.(validator.ValidationErrors); ok {
			return myvalidator.FieldErrors{
				Errors: fieldErrors,
			}
		}
		return myvalidator.ErrGeneric
	}

	return nil
}
