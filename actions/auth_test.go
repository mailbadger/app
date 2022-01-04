package actions_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/opa"
	"github.com/mailbadger/app/session"
	"github.com/mailbadger/app/sqs"
	"github.com/mailbadger/app/storage"
	s3mock "github.com/mailbadger/app/storage/s3"
)

func TestAuth(t *testing.T) {
	db := storage.New(config.Config{
		Database: config.Database{
			Driver:        "sqlite3",
			Sqlite3Source: ":memory:",
		},
	})
	s := storage.From(db)
	sess := session.New(s, "foo", "secretexmplkeythatis32characters", true)

	mockS3 := new(s3mock.MockS3Client)
	mockS3.On("PutObject", mock.AnythingOfType("*s3.PutObjectInput")).Twice().Return(&s3.PutObjectAclOutput{}, nil)

	mockPub := new(sqs.MockPublisher)

	compiler, err := opa.NewCompiler()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	e := setup(t, s, sess, mockS3, mockPub, compiler)

	// test when signup is disabled
	e.POST("/api/signup").WithJSON(params.PostSignUp{
		Email:    "foo@bar.com",
		Password: "test1234",
	}).
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().
		Value("message").
		Equal("Sign up is disabled.")

	err = os.Setenv("ENABLE_SIGNUP", "true")
	assert.Nil(t, err)

	e.POST("/api/signup").
		Expect().
		Status(http.StatusUnprocessableEntity).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again.")

	e.POST("/api/signup").WithJSON(params.PostSignUp{
		Email:    "email",
		Password: "password",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		Value("errors").Object().
		ValueEqual("email", "Invalid email format")

	e.POST("/api/signup").WithJSON(params.PostSignUp{
		Email:    "email",
		Password: "password",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		Value("errors").Object().
		ValueEqual("email", "Invalid email format")

	userObj := e.POST("/api/signup").WithJSON(params.PostSignUp{
		Email:    "gl@mail.com",
		Password: "password",
	}).Expect().
		Status(http.StatusOK).
		JSON().Object()

	userObj.Value("user").Object().
		ValueEqual("username", "gl@mail.com").
		ValueEqual("source", "mailbadger.io").
		ValueEqual("active", true).
		ValueEqual("verified", false).
		Value("boundaries").Object().
		ValueEqual("type", entities.BoundaryTypeFree)

	userObj.Value("user").Object().
		Value("roles").
		Array().
		NotEmpty().
		ContainsOnly(entities.Role{ID: 1, Name: entities.AdminRole})

	e.POST("/api/signup").WithJSON(params.PostSignUp{
		Email:    "gl@mail.com",
		Password: "password",
	}).Expect().
		Status(http.StatusForbidden).
		JSON().Object().
		ValueEqual("message", "Unable to create an account.")

	e.POST("/api/authenticate").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again")

	e.POST("/api/authenticate").WithJSON(params.PostAuthenticate{
		Username: "username",
		Password: "password",
	}).Expect().
		Status(http.StatusForbidden).
		JSON().Object().
		ValueEqual("message", "Invalid credentials.")

	e.POST("/api/authenticate").WithJSON(params.PostAuthenticate{
		Username: "gl@mail.com",
		Password: "badpassword",
	}).Expect().
		Status(http.StatusForbidden).
		JSON().Object().
		ValueEqual("message", "Invalid credentials.")

	e.POST("/api/authenticate").WithJSON(params.PostAuthenticate{
		Username: "gl@mail.com",
		Password: "password",
	}).Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("user").Object().
		ValueEqual("username", "gl@mail.com").
		ValueEqual("source", "mailbadger.io").
		ValueEqual("active", true)

	e.GET("/api/auth/github").
		Expect().
		Status(http.StatusTemporaryRedirect)

	e.GET("/api/auth/github/callback").
		Expect().
		Status(http.StatusPermanentRedirect)

	e.GET("/api/auth/google").
		Expect().
		Status(http.StatusTemporaryRedirect)

	e.GET("/api/auth/google/callback").
		Expect().
		Status(http.StatusPermanentRedirect)

	e.GET("/api/auth/facebook").
		Expect().
		Status(http.StatusTemporaryRedirect)

	e.GET("/api/auth/facebook/callback").
		Expect().
		Status(http.StatusPermanentRedirect)
}
