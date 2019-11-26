package token

import (
	"github.com/dgrijalva/jwt-go"
)

const (
	SignerAlgorithm = "HS256"
)

type SecretFunc func(*Token) (string, error)

type Token struct {
	Type  string
	Value string
}

func New(t, v string) *Token {
	return &Token{Type: t, Value: v}
}

// ParseToken parses the token from the raw string and returns it
func ParseToken(tokenStr string, secretFn SecretFunc) (*Token, error) {
	token := &Token{}
	parsedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {

		if t.Method.Alg() != SignerAlgorithm {
			return nil, jwt.ErrSignatureInvalid
		}

		claims := t.Claims.(jwt.MapClaims)
		typev, ok := claims["type"]
		if !ok {
			return nil, jwt.ValidationError{}
		}

		val, ok := claims["value"]
		if !ok {
			return nil, jwt.ValidationError{}
		}

		token.Type = typev.(string)
		token.Value = val.(string)

		secret, err := secretFn(token)
		return []byte(secret), err
	})

	if err != nil {
		return nil, err
	} else if !parsedToken.Valid {
		return nil, jwt.ValidationError{}
	}

	return token, nil
}

// SignWithExp signs the token using the given secret with an expiration date.
func (t *Token) SignWithExp(secret string, expiration int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["type"] = t.Type
	claims["value"] = t.Value
	if expiration > 0 {
		claims["exp"] = float64(expiration)
	}

	return token.SignedString([]byte(secret))
}
