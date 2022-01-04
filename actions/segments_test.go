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

func TestSegments(t *testing.T) {
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

	e.POST("/api/segments").WithJSON(params.Segment{Name: "djale"}).
		Expect().
		Status(http.StatusUnauthorized)

	auth.POST("/api/segments").WithJSON(params.Segment{Name: ""}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{"name": "This field is required"})

	// test post segment
	auth.POST("/api/segments").WithJSON(params.Segment{Name: "djale"}).
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
	auth.PUT("/api/segments/1").WithJSON(params.Segment{Name: "djaleputtest"}).
		Expect().
		Status(http.StatusOK)

	// delete segment by id
	auth.DELETE("/api/segments/1").
		Expect().
		Status(http.StatusNoContent)
}
