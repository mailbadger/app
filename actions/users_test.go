package actions_test

import (
	"net/http"
	"testing"

	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/s3"
)

func TestUser(t *testing.T) {
	db := storage.New(config.Config{
		Database: config.Database{
			Driver:        "sqlite3",
			Sqlite3Source: ":memory:",
		},
	})
	s := storage.From(db)

	s3mock := new(s3.MockS3Client)

	e := setup(t, s, s3mock)
	auth, err := createAuthenticatedExpect(e, s)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	e.GET("/api/users/me").
		Expect().
		Status(http.StatusUnauthorized).
		JSON().Object().ValueEqual("message", "You are not authorized to perform this request.")

	auth.GET("/api/users/me").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("username", "john").
		ValueEqual("active", true)

	auth.POST("/api/users/password").WithJSON(params.ChangePassword{
		Password:    "foo",
		NewPassword: "hunter2",
	}).Expect().Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		Value("errors").Object().ValueEqual("new_password", "Must be at least 8 character long")

	auth.POST("/api/users/password").WithJSON(params.ChangePassword{
		Password:    "foo",
		NewPassword: "hunter2foobar",
	}).Expect().Status(http.StatusForbidden).
		JSON().Object().
		ValueEqual("message", "The password that you entered is incorrect.")

	auth.POST("/api/users/password").WithJSON(params.ChangePassword{
		Password:    "hunter1",
		NewPassword: "hunter2foobar",
	}).Expect().Status(http.StatusOK).
		JSON().Object().
		ValueEqual("message", "Your password was updated successfully.")
}
