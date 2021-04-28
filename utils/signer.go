package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// SignData signs data string using HMAC SHA256 signer algorithm with the provided key.
func SignData(data, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	_, err := h.Write([]byte(data))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
