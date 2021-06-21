package middleware

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/sirupsen/logrus"
)

// CSRF returns a handler which performs checks for CSRF tokens.
func CSRF() gin.HandlerFunc {
	secureCookie, _ := strconv.ParseBool(os.Getenv("SECURE_COOKIE"))
	csrfMd := csrf.Protect([]byte(os.Getenv("SESSION_AUTH_KEY")),
		csrf.MaxAge(0),
		csrf.Secure(secureCookie),
		csrf.Path("/api"),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, err := w.Write([]byte(`{"message": "Forbidden - CSRF token invalid"}`))
			if err != nil {
				logrus.Error(err)
			}
		})),
	)

	return adapter.Wrap(csrfMd)
}
