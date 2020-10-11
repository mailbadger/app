package actions_test

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func testCampaigns(e *httpexpect.Expect) {
	type campaign struct {
		Name         string `json:"name"`
		TemplateName string `json:"template_name"`
	}
	logrus.Info(e)

	// test post campaigns
	e.POST("/campaigns").WithForm(campaign{
		Name:         "test1",
		TemplateName: "test1",
	}).
		Expect().
		Status(http.StatusOK)

	// test get campaigns by id
	e.GET("/campaigns/1").
		Expect().
		Status(http.StatusOK)
}
