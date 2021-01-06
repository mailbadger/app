package actions_test

import (
	"net/http"
	"testing"

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

	e.POST("/api/templates").WithForm(params.PostTemplate{Name: "", HTMLPart: "", TextPart: "", SubjectPart: ""}).Expect().Status(http.StatusUnauthorized)

	// test post template
	tempPostResp := auth.POST("/api/templates").WithForm(params.PostTemplate{Name: "test1", HTMLPart: "<html> bla </html>", TextPart: "txtpart", SubjectPart: "subpart"}).
		Expect().
		Status(http.StatusCreated)
	templateName := tempPostResp.JSON().Object().Value("name").String().Raw()

	// create template for test put case
	auth.POST("/api/templates").WithForm(params.PostTemplate{Name: "test2", HTMLPart: "<html> bla </html>", TextPart: "txtpart", SubjectPart: "subpart"}).
		Expect().
		Status(http.StatusCreated)

	// test put template
	auth.PUT("/api/templates/2").WithForm(params.PutTemplate{Name: "test3", HTMLPart: "<html> bla </html>", TextPart: "txtpart", SubjectPart: "subpart"}).
		Expect().
		Status(http.StatusOK)

	e.POST("/api/campaigns").WithForm(params.Campaign{Name: "djale", TemplateName: templateName}).
		Expect().
		Status(http.StatusUnauthorized)

	auth.POST("/api/campaigns").WithForm(params.Campaign{Name: "", TemplateName: ""}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{"name": "This field is required", "template_name": "This field is required"})

	// test post campaign
	auth.POST("/api/campaigns").WithForm(params.Campaign{Name: "foo1", TemplateName: templateName}).
		Expect().
		Status(http.StatusCreated)

	auth.POST("/api/campaigns").WithForm(params.Campaign{Name: "foo2", TemplateName: templateName}).
		Expect().
		Status(http.StatusCreated)

	auth.POST("/api/campaigns").WithForm(params.Campaign{Name: "test-scopes", TemplateName: templateName}).
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
		WithQuery("scopes[name]", "foo").
		Expect().
		Status(http.StatusOK).
		JSON().Object().ValueEqual("total", 2)

	// test inserted campaign
	auth.GET("/api/campaigns/1").
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual("name", "foo1").
		ValueEqual("status", "draft")

	auth.PUT("/api/campaigns/1").WithForm(params.Campaign{Name: "TESTputtest", TemplateName: templateName}).
		Expect().
		Status(http.StatusNoContent)

	// test updated campaign
	auth.GET("/api/campaigns/1").
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual("name", "TESTputtest").
		ValueEqual("status", "draft")

	// delete campaign by id
	auth.DELETE("/api/campaigns/1").
		Expect().
		Status(http.StatusNoContent)
}
