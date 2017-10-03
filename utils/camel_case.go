package utils

import (
	"bytes"
	"regexp"
)

var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

// CamelCase converts the given string to camel case.
// Ex: First Name -> FirstName etc..
func CamelCase(src string) string {
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	return string(bytes.Join(chunks, nil))
}
