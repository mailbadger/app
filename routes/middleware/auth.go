package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/storage"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

// Authorization header prefixes.
const (
	APIKeyAuth = "X-API-Key"
)

// SetUser fetches the token and then from the token fetches the user entity
// and sets it to the context.
func SetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var authHeader = c.GetHeader(APIKeyAuth)

		if authHeader != "" {
			key, err := storage.GetAPIKey(c, authHeader)
			if err != nil {
				logger.From(c).WithError(err).Error("unable to fetch api key")
				c.Next()
				return
			}

			c.Set("user", &key.User)

			// When using api keys it's ok to skip the csrf token
			// since we are not using cookies to authenticate the user
			c.Request = csrf.UnsafeSkipCheck(c.Request)
			c.Next()
			return
		}

		session := sessions.Default(c)
		v := session.Get("sess_id")
		if v == nil {
			c.Next()
			return
		}
		sessID := v.(string)
		s, err := storage.GetSession(c, sessID)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user", &s.User)

		c.Next()
	}
}

// GetUser returns the user set in the context
func GetUser(c *gin.Context) *entities.User {
	val, ok := c.Get("user")
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
func Authorized(compiler *ast.Compiler) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("user")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
			return
		}
		u, ok := val.(*entities.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
			return
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

		// Run evaluation.
		rs, err := rego.Eval(c)

		if err != nil {
			logger.From(c).WithField("input", input).WithError(err).Error("auth: unable to evaluate opa decision")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
			return
		}

		if len(rs) == 0 {
			logger.From(c).WithField("input", input).Error("auth: undefined opa decision")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
			return
		}

		if decision, ok := rs[0].Expressions[0].Value.(bool); !ok || len(rs) > 1 {
			logger.From(c).WithField("input", input).Error("auth: non-boolean opa decision")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
			return
		} else if !decision {
			logger.From(c).WithField("input", input).Info("auth: user does not have required permissions")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not authorized"})
			return
		}

		c.Next()
	}
}
