package actions_test

import (
	"net/http"
	"testing"

	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/s3"
)

func TestSegments(t *testing.T) {
	s := storage.New("sqlite3", ":memory:")

	s3mock := new(s3.MockS3Client)

	e := setup(t, s, s3mock)
	auth, err := createAuthenticatedExpect(e, s)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	e.POST("/api/segments").WithForm(params.Segment{Name: "djale"}).
		Expect().
		Status(http.StatusUnauthorized)

	auth.POST("/api/segments").WithForm(params.Segment{Name: ""}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{"name": "This field is required"})

	// test post segment
	auth.POST("/api/segments").WithForm(params.Segment{Name: "djale"}).
		Expect().
		Status(http.StatusCreated)

	// test inserted segment
	auth.GET("/api/segments/1").
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual("name", "djale").
		ValueEqual("subscribers_in_segment", 0).
		ValueEqual("total_subscribers", 0)

	// test put segments
	auth.PUT("/api/segments/1").WithForm(params.Segment{Name: "djaleputtest"}).
		Expect().
		Status(http.StatusOK)

	// delete segment by id
	auth.DELETE("/api/segments/1").
		Expect().
		Status(http.StatusNoContent)
}
