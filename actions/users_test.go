package actions_test

import (
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/mailbadger/app/actions"
	"github.com/mailbadger/app/routes/middleware"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
	os.Setenv("DATABASE_DRIVER", "sqlite3")
	os.Setenv("DATABASE_CONFIG", ":memory:")
	os.Setenv("SESSION_AUTH_KEY", "foo")
	os.Setenv("SESSION_AUTH_KEY", "secretexmplkeythatis32characters")
}

func TestUnauthorizedUser(t *testing.T) {
	store := cookie.NewStore(
		[]byte(os.Getenv("SESSION_AUTH_KEY")),
		[]byte(os.Getenv("SESSION_ENCRYPT_KEY")),
	)
	secureCookie, _ := strconv.ParseBool(os.Getenv("SECURE_COOKIE"))
	store.Options(sessions.Options{
		Secure:   secureCookie,
		HttpOnly: true,
	})

	handler := gin.New()
	handler.Use(sessions.Sessions("mbsess", store))
	handler.Use(middleware.Storage())
	handler.Use(middleware.SetUser())

	authorized := handler.Group("/api")
	authorized.Use(middleware.Authorized())
	{
		users := authorized.Group("/users")
		{
			users.GET("/me", actions.GetMe)
			users.POST("/password", actions.ChangePassword)
		}
	}
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	e.GET("/api/users/me").
		Expect().
		Status(http.StatusUnauthorized).
		JSON().Object().ValueEqual("message", "User not authorized")

}
