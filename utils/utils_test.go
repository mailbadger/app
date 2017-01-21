package utils

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomBytes(t *testing.T) {
	b, err := GenerateRandomBytes(16)
	if err != nil {
		t.Error("could not generate random bytes: ", err)
	}

	assert.Len(t, b, 16, "they should be equal")
}

func TestGenerateRandomString(t *testing.T) {
	str, err := GenerateRandomString(16)
	if err != nil {
		t.Error("could not generate random string: ", err)
	}

	decoded, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		t.Error("could not decode random string from base64: ", err)
	}

	assert.Len(t, decoded, 16, "they should be equal")
}

func TestCamelCase(t *testing.T) {
	const testStr = "First Name"
	if camelCase := CamelCase(testStr); camelCase != "FirstName" {
		t.Errorf("could not transform string %v to %v", testStr, "FirstName")
	}
}
