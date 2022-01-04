package actions_test

import (
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/opa"
	"github.com/mailbadger/app/session"
	"github.com/mailbadger/app/sqs"
	"github.com/mailbadger/app/storage"
	s3mock "github.com/mailbadger/app/storage/s3"
	"github.com/stretchr/testify/mock"
)

func TestUser(t *testing.T) {
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
