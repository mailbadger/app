package middleware

import (
	"errors"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/session"
	"github.com/mailbadger/app/storage"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

// Authorization header prefixes.
const (
	APIKeyAuth = "X-API-Key"
	userKey    = "user"
)

// GetUser returns the user set in the context
func GetUser(c *gin.Context) *entities.User {
	val, ok := c.Get(userKey)
	if !ok {
		return nil
	}

	user, ok := val.(*entities.User)
	if !ok {
		return nil
	}

	return user
}

// Authorized is a middleware that checks if the user is authorized to do the
// requested action.
func Authorized(sess session.Session, compiler *ast.Compiler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u *entities.User

		authHeader := c.GetHeader(APIKeyAuth)
		if authHeader != "" {
			key, err := storage.GetAPIKey(c, authHeader)
			if err != nil {
				logger.From(c).WithError(err).Error("unable to fetch api key")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this request."})
				return
			}

			u = &key.User
			// When using api keys it's ok to skip the csrf token
			// since we are not using cookies to authenticate the user
			c.Request = csrf.UnsafeSkipCheck(c.Request)
		} else {
			s, err := sess.GetUserSession(c)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) && !errors.Is(err, session.ErrNotFound) {
					logrus.WithError(err).Error("authorized: unable to get user session")
				}
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this request."})
				return
			}

			u = &s.User
		}

		// Create a new query that uses the compiled policy from above.
		input := map[string]interface{}{
			"roles":  u.RoleNames(),
			"method": c.Request.Method,
			"path":   c.FullPath(),
		}
		rego := rego.New(
			rego.Query("data.rbac.authz.allow"),
			rego.Compiler(compiler),
			rego.Input(input),
		)

		// Run auth evaluation.
		rs, err := rego.Eval(c)

		if err != nil {
			logger.From(c).WithField("input", input).WithError(err).Error("auth: unable to evaluate opa decision")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this request."})
			return
		}

		if len(rs) == 0 {
			logger.From(c).WithField("input", input).Error("auth: undefined opa decision")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this request."})
			return
		}

		if decision, ok := rs[0].Expressions[0].Value.(bool); !ok || len(rs) > 1 {
			logger.From(c).WithField("input", input).Error("auth: non-boolean opa decision")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this request."})
			return
		} else if !decision {
			logger.From(c).WithField("input", input).Info("auth: user does not have the required permissions")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to perform this request."})
			return
		}

		c.Set(userKey, u)

		entry := logger.From(c).WithField("user_id", u.ID)
		logger.SetToContext(c, entry)

		c.Next()
	}
}
