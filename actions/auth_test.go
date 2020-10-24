package actions_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/storage"
)

func TestAuth(t *testing.T) {
	s := storage.New("sqlite3", ":memory:")

	e := setup(t, s)

	// test when signup is disabled
	e.POST("/api/signup").WithForm(params.PostSignUp{
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
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		Value("errors").Object().
		ValueEqual("email", "This field is required").
		ValueEqual("password", "This field is required")

	e.POST("/api/signup").WithForm(params.PostSignUp{
		Email:    "email",
		Password: "password",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		Value("errors").Object().
		ValueEqual("email", "Invalid email format")

	e.POST("/api/signup").WithForm(params.PostSignUp{
		Email:    "email",
		Password: "password",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		Value("errors").Object().
		ValueEqual("email", "Invalid email format")

	e.POST("/api/signup").WithForm(params.PostSignUp{
		Email:    "gl@mail.com",
		Password: "password",
	}).Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("user").Object().
		ValueEqual("username", "gl@mail.com").
		ValueEqual("source", "mailbadger.io").
		ValueEqual("active", true).
		ValueEqual("verified", false)

	e.POST("/api/signup").WithForm(params.PostSignUp{
		Email:    "gl@mail.com",
		Password: "password",
	}).Expect().
		Status(http.StatusForbidden).
		JSON().Object().
		ValueEqual("message", "Unable to create an account.")
}
