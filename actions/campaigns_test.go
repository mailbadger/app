package actions_test

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
)

func TestGetCampaignAction(t *testing.T) {
	// Create new gin instance
	handler := gin.New()

	// Create httpexpect instance
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	testCampaigns(e)

}


// this func will run all actions from campaign
func testCampaigns(e *httpexpect.Expect) {
	type campaign struct {
		Name         string `json:"name"`
		TemplateName string `json:"template_name"`
	}


	// if we need to auth we can call auth api and use token from it to run
	// other campaigns actions like post campaign and etc.
	// example for this
	r := e.POST("/login").WithForm(Login{"username", "pw"}).
		Expect().
		Status(http.StatusOK).JSON().Object()

	// take token from auth and use it on others actions..
	r.Keys().ContainsOnly("token")

	// test post campaigns
	e.POST("/campaigns").WithForm(campaign{
		Name:         "test1",
		TemplateName: "test1",
	}).
		Expect().
		Status(http.StatusOK)

	// test get campaigns by id ( same campaign we create with post campaigns test.
	e.GET("/campaigns/1").
		Expect().
		Status(http.StatusOK)
}
