package actions_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/s3"
)

func TestAuth(t *testing.T) {
	s := storage.New("sqlite3", ":memory:")

	s3mock := new(s3.MockS3Client)

	e := setup(t, s, s3mock)

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

	err := os.Setenv("ENABLE_SIGNUP", "true")
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
