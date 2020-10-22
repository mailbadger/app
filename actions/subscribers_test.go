package actions_test

import (
	"net/http"
	"testing"

	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/storage"
)

func TestSubscribers(t *testing.T) {
	s := storage.New("sqlite3", ":memory:")

	e := setup(t, s)
	auth, err := createAuthenticatedExpect(e, s)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	e.POST("/api/subscribers").WithForm(params.PostSubscriber{Name: "", Email: "", Metadata: map[string]string{"": ""}, SegmentIDs: nil}).
		Expect().
		Status(http.StatusUnauthorized)

	// test binding on post subscriber
	auth.POST("/api/subscribers").WithForm(params.PostSubscriber{Name: "", Email: "", Metadata: map[string]string{"": ""}, SegmentIDs: nil}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{"email": "This field is required"})

	// test validate email
	auth.POST("/api/subscribers").WithForm(params.PostSubscriber{Name: "Djale", Email: "sda", Metadata: map[string]string{"": ""}, SegmentIDs: nil}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{"email": "Invalid email format"})

	// test post subscriber
	auth.POST("/api/subscribers").WithForm(params.PostSubscriber{Name: "Djale", Email: "djale@email.com", Metadata: map[string]string{"test": "test"}, SegmentIDs: nil}).
		Expect().
		Status(http.StatusCreated)

	// test post subscriber ( insert more then 1 sub)
	auth.POST("/api/subscribers").WithForm(params.PostSubscriber{Name: "Foo", Email: "foo@email.com", Metadata: map[string]string{"test": "test"}, SegmentIDs: nil}).
		Expect().
		Status(http.StatusCreated)

	// test get subscribers
	auth.GET("/api/subscribers").
		Expect().
		Status(http.StatusOK)
}
