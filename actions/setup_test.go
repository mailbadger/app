package actions_test

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/opa"
	"github.com/mailbadger/app/routes"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/s3"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setup(t *testing.T, s storage.Storage, s3Mock *s3.MockS3Client) *httpexpect.Expect {
	err := os.Setenv("SESSION_AUTH_KEY", "foo")
	if err != nil {
		t.FailNow()
	}
	err = os.Setenv("SESSION_ENCRYPT_KEY", "secretexmplkeythatis32characters")
	if err != nil {
		t.FailNow()
	}

	cookiestore := cookie.NewStore(
		[]byte(os.Getenv("SESSION_AUTH_KEY")),
		[]byte(os.Getenv("SESSION_ENCRYPT_KEY")),
	)
	secureCookie, _ := strconv.ParseBool(os.Getenv("SECURE_COOKIE"))
	cookiestore.Options(sessions.Options{
		Secure:   secureCookie,
		HttpOnly: true,
	})

	handler := gin.New()
	handler.Use(sessions.Sessions("mbsess", cookiestore))
	handler.Use(middleware.Storage(s))
	handler.Use(middleware.SetUser())
	handler.Use(middleware.S3Client(s3Mock))

	routes.SetGuestRoutes(handler)

	compiler, err := opa.NewCompiler()
	if err != nil {
		t.FailNow()
	}
	routes.SetAuthorizedRoutes(handler, compiler)

	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCompactPrinter(t),
		},
	})
}

func createAuthenticatedExpect(e *httpexpect.Expect, s storage.Storage) (*httpexpect.Expect, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte("hunter1"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	b, err := s.GetBoundariesByType("db_test")
	if err != nil {
		return nil, err
	}

	u := &entities.User{
		Active:   true,
		Username: "john",
		Password: sql.NullString{
			String: string(pass),
			Valid:  true,
		},
		Boundaries: b,
		Roles: []entities.Role{
			{Name: "admin"},
		},
	}
	err = s.CreateUser(u)
	if err != nil {
		return nil, err
	}

	c := e.POST("/api/authenticate").WithJSON(params.PostAuthenticate{
		Username: "john",
		Password: "hunter1",
	}).Expect().Status(http.StatusOK).Cookie("mbsess")

	return e.Builder(func(req *httpexpect.Request) {
		req.WithCookie(c.Name().Raw(), c.Value().Raw())
	}), nil
}
