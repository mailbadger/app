package actions_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/storage"
)

func TestCampaigns(t *testing.T) {
	s := storage.New("sqlite3", ":memory:")

	e := setup(t, s)
	auth, err := createAuthenticatedExpect(e, s)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// create template for fetching template in post/put camp
	template := &entities.Template{
		UserID:      1,
		Name:        "template_test_name",
		HTMLPart:    "",
		TextPart:    "",
		SubjectPart: "",
	}
	err = s.CreateTemplate(template)
	assert.Nil(t, err)

	e.POST("/api/campaigns").WithForm(params.Campaign{Name: "djale", TemplateName: "template_test_name"}).
		Expect().
		Status(http.StatusUnauthorized)

	auth.POST("/api/campaigns").WithForm(params.Campaign{Name: "", TemplateName: ""}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{"name": "This field is required", "template_name": "This field is required"})

	// test post campaign
	auth.POST("/api/campaigns").WithForm(params.Campaign{Name: "foo1", TemplateName: "template_test_name"}).
		Expect().
		Status(http.StatusCreated)

	auth.POST("/api/campaigns").WithForm(params.Campaign{Name: "foo2", TemplateName: "template_test_name"}).
		Expect().
		Status(http.StatusCreated)

	auth.POST("/api/campaigns").WithForm(params.Campaign{Name: "test-scopes", TemplateName: "template_test_name"}).
		Expect().
		Status(http.StatusCreated)

	// test scopes
	collection := auth.GET("/api/campaigns").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ValueEqual("total", 3)

	collection.Value("links").Object().ContainsKey("previous").ContainsKey("next")
	collection.Value("collection").Array().NotEmpty().Length().Equal(3)

	auth.GET("/api/campaigns").
		WithQuery("scopes[name]", "foo").WithQuery("scopes[template_name]", "template_test_name").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ValueEqual("total", 2)

	// test inserted campaign
	auth.GET("/api/campaigns/1").
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual("name", "foo1").
		ValueEqual("template_name", "template_test_name").
		ValueEqual("status", "draft")

	auth.PUT("/api/campaigns/1").WithForm(params.Campaign{Name: "djaleputtest", TemplateName: "template_test_name"}).
		Expect().
		Status(http.StatusNoContent)

	// test updated campaign
	auth.GET("/api/campaigns/1").
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual("name", "djaleputtest").
		ValueEqual("template_name", "template_test_name").
		ValueEqual("status", "draft")

	// delete campaign by id
	auth.DELETE("/api/campaigns/1").
		Expect().
		Status(http.StatusNoContent)
}
