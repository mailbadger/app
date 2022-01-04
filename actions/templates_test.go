package actions_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/mock"

	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/opa"
	"github.com/mailbadger/app/session"
	"github.com/mailbadger/app/sqs"
	"github.com/mailbadger/app/storage"
	s3mock "github.com/mailbadger/app/storage/s3"
)

func TestTemplates(t *testing.T) {
	db := storage.New(config.Config{
		Database: config.Database{
			Driver:        "sqlite3",
			Sqlite3Source: ":memory:",
		},
	})
	s := storage.From(db)
	sess := session.New(s, "foo", "secretexmplkeythatis32characters", true)

	mockS3 := new(s3mock.MockS3Client)

	readCloser := ioutil.NopCloser(strings.NewReader("hello world"))

	mockS3.On("PutObject", mock.AnythingOfType("*s3.PutObjectInput")).Once().Return(nil, errors.New("error"))
	mockS3.On("PutObject", mock.AnythingOfType("*s3.PutObjectInput")).Twice().Return(&s3.PutObjectAclOutput{}, nil)
	mockS3.On("GetObject", mock.AnythingOfType("*s3.GetObjectInput")).Once().Return(nil, awserr.New(s3.ErrCodeNoSuchKey, "no such ky", errors.New("key not found")))
	mockS3.On("GetObject", mock.AnythingOfType("*s3.GetObjectInput")).Once().Return(nil, awserr.New(s3.ErrCodeInvalidObjectState, "invalid object state", errors.New("invalid object state")))
	mockS3.On("GetObject", mock.AnythingOfType("*s3.GetObjectInput")).Once().Return(nil, errors.New("some error"))
	mockS3.On("GetObject", mock.AnythingOfType("*s3.GetObjectInput")).Once().Return(&s3.GetObjectOutput{
		Body: readCloser,
	}, nil)
	mockS3.On("DeleteObject", mock.AnythingOfType("*s3.DeleteObjectInput")).Twice().Return(&s3.DeleteObjectOutput{}, nil)

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

	// test post template unauthorized
	e.POST("/api/templates").WithJSON(params.PostTemplate{
		Name:        "",
		HTMLPart:    "",
		TextPart:    "",
		SubjectPart: "",
	}).Expect().
		Status(http.StatusUnauthorized)

	// test binding on post template
	auth.POST("/api/templates").WithJSON(params.PostTemplate{
		Name:        "",
		HTMLPart:    "",
		TextPart:    "",
		SubjectPart: "",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{
			"html_part":    "This field is required",
			"name":         "This field is required",
			"subject_part": "This field is required",
			"text_part":    "This field is required",
		})

	// TODO fix test for html validation
	/*auth.POST("/api/templates").WithJSON(params.PostTemplate{
		Name:        "template 1",
		HTMLPart:    "test 1 template<div>hjhjhj",
		TextPart:    "template {{.number}} 123",
		SubjectPart: "hello {{.name}}",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{
			"html_part":    "This field is required",
		})*/

	// test filed to parse text part on post template
	auth.POST("/api/templates").WithJSON(params.PostTemplate{
		Name:        "template 3",
		HTMLPart:    "<span>test 2 template<span>",
		TextPart:    "template {{{number}} 223",
		SubjectPart: "hello {{.name}}",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Unable to create template, failed to parse text_part")

	// test filed to parse subject part on post template
	auth.POST("/api/templates").WithJSON(params.PostTemplate{
		Name:        "template 3",
		HTMLPart:    "<span>test 2 template<span>",
		TextPart:    "template {{number}} 223",
		SubjectPart: "hello {{{name}}",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Unable to create template, failed to parse subject_part")

	// test filed to parse html part on post template
	auth.POST("/api/templates").WithJSON(params.PostTemplate{
		Name:        "template 3",
		HTMLPart:    "<span>test 2 template<span>{{{tesT}}",
		TextPart:    "template {{number}} 223",
		SubjectPart: "hello {{name}}",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Unable to create template, failed to parse html_part")

	// test post subscriber with error on PutObject (this template is saved in database only the html part is not saved)
	auth.POST("/api/templates").WithJSON(params.PostTemplate{
		Name:        "template 1",
		HTMLPart:    "<span>test 1 template<span>",
		TextPart:    "template {{.number}} 123",
		SubjectPart: "hello {{.name}}",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Unable to create template, please try again.")

	// test post template
	id := auth.POST("/api/templates").WithJSON(params.PostTemplate{
		Name:        "template 2",
		HTMLPart:    "<span>test 2 template<span>",
		TextPart:    "template {{.number}} 223",
		SubjectPart: "hello {{.name}}",
	}).Expect().
		Status(http.StatusCreated).
		JSON().Object().Value("id")

	// test post template with name that exists
	auth.POST("/api/templates").WithJSON(params.PostTemplate{
		Name:        "template 2",
		HTMLPart:    "<span>test 2 template<span>",
		TextPart:    "template {{.number}} 223",
		SubjectPart: "hello {{.name}}",
	}).Expect().
		Status(http.StatusUnprocessableEntity).
		JSON().Object().
		ValueEqual("message", "Template with that name already exists")

	idStr := strconv.FormatFloat(id.Raw().(float64), 'f', 0, 64)

	// test invalid parameters on put template
	auth.PUT("/api/templates/"+idStr).WithJSON(params.PutTemplate{
		Name:        "",
		HTMLPart:    "",
		TextPart:    "",
		SubjectPart: "",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{
			"html_part":    "This field is required",
			"name":         "This field is required",
			"subject_part": "This field is required",
			"text_part":    "This field is required",
		})

	// test put template for non existing template
	auth.PUT("/api/templates/9933209").WithJSON(params.PutTemplate{
		Name:        "",
		HTMLPart:    "",
		TextPart:    "",
		SubjectPart: "",
	}).Expect().
		Status(http.StatusNotFound).
		JSON().Object().
		ValueEqual("message", "Template not found")

	// test put template with non integer id
	auth.PUT("/api/templates/2.2").WithJSON(params.PutTemplate{
		Name:        "",
		HTMLPart:    "",
		TextPart:    "",
		SubjectPart: "",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Id must be an integer")

	// test put template
	auth.PUT("/api/templates/" + idStr).WithJSON(params.PutTemplate{
		HTMLPart:    "<span>test 2 template updated<span>",
		TextPart:    "template {{.number}} 223 updated",
		SubjectPart: "hello {{.name}}",
		Name:        "template 2 updated",
	}).Expect().
		Status(http.StatusOK)

	// test put template with name that already exists
	auth.PUT("/api/templates/"+idStr).WithJSON(params.PutTemplate{
		HTMLPart:    "<span>test 2 template updated<span>",
		TextPart:    "template {{.number}} 223 updated",
		SubjectPart: "hello {{.name}}",
		Name:        "template 1",
	}).Expect().
		Status(http.StatusUnprocessableEntity).
		JSON().Object().
		ValueEqual("message", "Template with that name already exists")

	// TODO fix test for html validation
	/*auth.PUT("/api/templates/" + idStr).WithJSON(params.PostTemplate{
		Name:        "template 2 uuds",
		HTMLPart:    "test 1 template<div>hjhjhj",
		TextPart:    "template {{.number}} 123",
		SubjectPart: "hello {{.name}}",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{
			"html_part":    "This field is required",
		})*/

	// test filed to parse text part on put template
	auth.PUT("/api/templates/"+idStr).WithJSON(params.PutTemplate{
		Name:        "template 3",
		HTMLPart:    "<span>test 2 template<span>",
		TextPart:    "template {{{number}} 223",
		SubjectPart: "hello {{.name}}",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Unable to update template, failed to parse text_part")

	// test filed to parse subject part on put template
	auth.PUT("/api/templates/"+idStr).WithJSON(params.PutTemplate{
		Name:        "template 3",
		HTMLPart:    "<span>test 2 template<span>",
		TextPart:    "template {{number}} 223",
		SubjectPart: "hello {{{name}}",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Unable to update template, failed to parse subject_part")

	// test filed to parse html part on put template
	auth.PUT("/api/templates/"+idStr).WithJSON(params.PutTemplate{
		Name:        "template 3",
		HTMLPart:    "<span>test 2 template<span>{{{tesT}}",
		TextPart:    "template {{number}} 223",
		SubjectPart: "hello {{name}}",
	}).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Unable to update template, failed to parse html_part")

	// test get template with id not integer
	auth.GET("/api/templates/2.2").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Id must be an integer")

	// test template not found on get template
	auth.GET("/api/templates/94829342").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().
		ValueEqual("message", "Template not found.")

	// test html part not found on get template
	auth.GET("/api/templates/1").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().
		ValueEqual("message", "HTML part not found.")

	// test invalid state of html part on get template
	auth.GET("/api/templates/1").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().
		ValueNotEqual("message", "The state of the HTML part is invalid")

	// test invalid state of html part on get template
	auth.GET("/api/templates/1").
		Expect().
		Status(http.StatusUnprocessableEntity).
		JSON().Object().
		ValueEqual("message", "Unable to get template")

	// TODO add test for get template successfully
	/*auth.GET("/api/templates/" + idStr).
	Expect().
	Status(http.StatusOK).
	JSON().Object().
	ValueEqual("message", "Unable to get template")*/

	// test list templates
	collection := auth.GET("/api/templates").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("total", 2)

	collection.Value("links").Object().ContainsKey("previous").ContainsKey("next")
	collection.Value("collection").Array().NotEmpty().Length().Equal(2)

	// test list templates with scopes
	collection = auth.GET("/api/templates").
		WithQuery("scopes[name]", "template 1").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ValueEqual("total", 1)

	collection.Value("links").Object().ContainsKey("previous").ContainsKey("next")
	collection.Value("collection").Array().NotEmpty().Length().Equal(1)

	// test list templates with scopes
	collection = auth.GET("/api/templates").
		WithQuery("scopes[name]", "9999").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ValueEqual("total", 0)

	collection.Value("links").Object().ContainsKey("previous").ContainsKey("next")
	collection.Value("collection").Array().Empty().Length().Equal(0)

	// test get template with id not integer
	auth.DELETE("/api/templates/2.2").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Id must be an integer")

	// test delete template with non existing id
	auth.DELETE("/api/templates/94829342").
		Expect().
		Status(http.StatusNoContent)

	// test delete template with non existing id
	auth.DELETE("/api/templates/" + idStr).
		Expect().
		Status(http.StatusNoContent)
}
