package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/sirupsen/logrus"
)

// CSRF returns a handler which performs checks for CSRF tokens.
func CSRF(authKey string, secure bool) gin.HandlerFunc {
	csrfMd := csrf.Protect([]byte(authKey),
		csrf.Secure(secure),
		csrf.Path("/api"),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := csrf.FailureReason(r)
			if err != nil {
				logrus.WithError(err).Error("csrf failure reason")
			}
			w.WriteHeader(http.StatusForbidden)
			_, err = w.Write([]byte(`{"message": "Forbidden - CSRF token invalid"}`))
			if err != nil {
				logrus.WithError(err).Error("csrf: unable to write message to writer")
			}
		})),
	)

	return adapter.Wrap(csrfMd)
}
