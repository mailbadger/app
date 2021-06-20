package actions_test

import (
	"net/http"
	"testing"

	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/s3"
)

func TestSubscribers(t *testing.T) {
	s := storage.New("sqlite3", ":memory:")

	s3mock := new(s3.MockS3Client)

	e := setup(t, s, s3mock)
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
		ValueEqual("message", "Invalid parameters, please try again")

	// test validator for params // using WithFormField because of  dif metadata encoding gin/httpexpect
	auth.POST("/api/subscribers").WithFormField("metadata[aaa aa]", "blabla").WithFormField("email", "sda").
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{"email": "Invalid email format", "metadata[aaa aa]": "Must consist only of alphanumeric and hyphen characters"})

	// test post subscriber
	auth.POST("/api/subscribers").WithForm(params.PostSubscriber{Name: "Djale", Email: "djale@email.com", Metadata: map[string]string{"test": "test"}, SegmentIDs: nil}).
		Expect().
		Status(http.StatusCreated)

	// test post subscriber ( insert more then 1 sub)
	auth.POST("/api/subscribers").WithForm(params.PostSubscriber{Name: "Foo", Email: "foo@email.com", Metadata: map[string]string{"test": "test"}, SegmentIDs: nil}).
		Expect().
		Status(http.StatusCreated)

	// test get subscribers collection length 2
	collection := auth.GET("/api/subscribers").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ValueEqual("total", 2)

	// test links from  collection
	collection.Value("links").Object().ContainsKey("previous").ContainsKey("next")

	// test collection[0] values objects
	collection.Value("collection").Array().Element(0).Object().
		ValueEqual("name", "Foo").
		ValueEqual("email", "foo@email.com").
		ValueEqual("blacklisted", false).
		ValueEqual("active", true)

	// test get subscribers by filter email like foo
	collection = auth.GET("/api/subscribers").WithQuery("scopes[email]", "foo").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ValueEqual("total", 1)

	// test collection[0] values objects filtered.
	collection.Value("collection").Array().Element(0).Object().
		ValueEqual("name", "Foo").
		ValueEqual("email", "foo@email.com").
		ValueEqual("blacklisted", false).
		ValueEqual("active", true)

	// test get subscriber by id
	auth.GET("/api/subscribers/2").
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual("name", "Foo").
		ValueEqual("email", "foo@email.com").
		ValueEqual("blacklisted", false).
		ValueEqual("active", true)

	// test put subscriber by id
	auth.PUT("/api/subscribers/2").WithForm(params.PutSubscriber{Name: "FooPutChange", Metadata: map[string]string{"test": "test"}}).
		Expect().
		Status(http.StatusOK)

	// test updated subscriber
	auth.GET("/api/subscribers/2").
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual("name", "FooPutChange").
		ValueEqual("email", "foo@email.com").
		ValueEqual("blacklisted", false).
		ValueEqual("active", true)

	// delete subscriber by id
	auth.DELETE("/api/subscribers/1").
		Expect().
		Status(http.StatusNoContent)

	// delete subscriber by id
	auth.DELETE("/api/subscribers/2").
		Expect().
		Status(http.StatusNoContent)
}
