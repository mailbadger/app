package actions_test

import (
	"database/sql"
	"net/http"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/open-policy-agent/opa/ast"

	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/mode"
	"github.com/mailbadger/app/routes"
	"github.com/mailbadger/app/services/boundaries"
	"github.com/mailbadger/app/services/reports"
	"github.com/mailbadger/app/services/subscribers"
	"github.com/mailbadger/app/services/templates"
	"github.com/mailbadger/app/session"
	"github.com/mailbadger/app/sqs"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/s3"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setup(
	t *testing.T,
	s storage.Storage,
	sess session.Session,
	s3Mock *s3.MockS3Client,
	pub sqs.PublisherAPI,
	emailSender emails.Sender,
	templatesvc templates.Service,
	boundarysvc boundaries.Service,
	subscrsvc subscribers.Service,
	reportsvc reports.Service,
	compiler *ast.Compiler,
	enableSignup bool,
	verifyEmail bool,
) *httpexpect.Expect {
	mode.SetMode("test")

	queueURL := "http://example.com/campaigns-queue"
	api := routes.New(
		sess,
		s,
		compiler,
		pub,
		s3Mock,
		emailSender,
		templatesvc,
		boundarysvc,
		subscrsvc,
		reportsvc,
		&queueURL,
		"/var/www/app",       // app dir
		"http://example.com", // app url
		"files-bucket",
		enableSignup,
		verifyEmail,
		"",                                 // recaptcha secret
		"secretexmplkeythatis32characters", // unsubscribe token secret
		"test@example.com",                 // system email
		config.Social{},
	)

	handler := api.Handler()

	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: "http://example.com",
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCurlPrinter(t),
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

	e = e.Builder(func(req *httpexpect.Request) {
		req.WithCookie(c.Name().Raw(), c.Value().Raw())
	})

	res := e.GET("/api/users/me").Expect()
	token := res.Header("X-CSRF-Token").Raw()
	csrfCookie := res.Cookie("_gorilla_csrf")

	return e.Builder(func(req *httpexpect.Request) {
		req.WithCookie(csrfCookie.Name().Raw(), csrfCookie.Value().Raw())
		req.WithHeader("X-CSRF-Token", token)
	}), nil
}
