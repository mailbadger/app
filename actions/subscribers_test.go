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
	awss3 "github.com/mailbadger/app/storage/s3"
	"github.com/stretchr/testify/mock"
)

func TestSubscribers(t *testing.T) {
	db := storage.New(config.Config{
		Database: config.Database{
			Driver:        "sqlite3",
			Sqlite3Source: ":memory:",
		},
	})
	s := storage.From(db)
	sess := session.New(s, "foo", "secretexmplkeythatis32characters", true)

	mockS3 := new(awss3.MockS3Client)
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

	e.POST("/api/subscribers").WithJSON(params.PostSubscriber{Name: "", Email: "", Metadata: map[string]string{"": ""}, SegmentIDs: nil}).
		Expect().
		Status(http.StatusUnauthorized)

	// test binding on post subscriber
	auth.POST("/api/subscribers").WithJSON(params.PostSubscriber{Name: "", Email: "", Metadata: map[string]string{"": ""}, SegmentIDs: nil}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again")

	// test validator for params // using WithJSONField because of  dif metadata encoding gin/httpexpect
	auth.POST("/api/subscribers").
		WithJSON(
			params.PostSubscriber{
				Name:     "",
				Email:    "aaa",
				Metadata: map[string]string{"aaa aa": "blabla"},
			}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{"email": "Invalid email format", "metadata[aaa aa]": "Must consist only of alphanumeric and hyphen characters"})

	// test post subscriber
	auth.POST("/api/subscribers").WithJSON(params.PostSubscriber{Name: "Djale", Email: "djale@email.com", Metadata: map[string]string{"test": "test"}, SegmentIDs: nil}).
		Expect().
		Status(http.StatusCreated)

	// test post subscriber ( insert more then 1 sub)
	auth.POST("/api/subscribers").WithJSON(params.PostSubscriber{Name: "Foo", Email: "foo@email.com", Metadata: map[string]string{"test": "test"}, SegmentIDs: nil}).
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
	auth.PUT("/api/subscribers/2").WithJSON(params.PutSubscriber{Name: "FooPutChange", Metadata: map[string]string{"test": "test"}}).
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
