package actions_test

import (
	"net/http"
	"testing"

	"github.com/mailbadger/app/storage"
)

func TestUser(t *testing.T) {
	s := storage.New("sqlite3", ":memory:")

	e := setup(t, s)
	auth, err := createAuthenticatedExpect(e, s)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	e.GET("/api/users/me").
		Expect().
		Status(http.StatusUnauthorized).
		JSON().Object().ValueEqual("message", "User not authorized")

	auth.GET("/api/users/me").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("username", "john").
		ValueEqual("active", true)
}
