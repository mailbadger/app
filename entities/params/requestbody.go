package params

import "github.com/go-playground/validator/v10"

type RequestBody interface {
	trimSpaces()
	Validate(v validator.Validate) (string, bool)
}