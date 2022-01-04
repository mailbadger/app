package actions_test

import (
	"net/http"
	"strconv"
	"testing"

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

func TestCampaigns(t *testing.T) {
	db := storage.New(config.Config{
		Database: config.Database{
			Driver:        "sqlite3",
			Sqlite3Source: ":memory:",
		},
	})
	s := storage.From(db)
	sess := session.New(s, "jXn2r5u8x/A?D(G+KbPeSgVkYp3s6v9y", "jXn2r5u8x/A?D(G+KbPeSgVkYp3s6v9y", true)

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

	// test post template
	templateName := auth.POST("/api/templates").WithJSON(params.PostTemplate{Name: "test1", HTMLPart: "<html> bla </html>", TextPart: "txtpart", SubjectPart: "subpart"}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().Value("name").String().Raw()

	e.POST("/api/campaigns").WithJSON(params.PostCampaign{Name: "djale", TemplateName: templateName}).
		Expect().
		Status(http.StatusUnauthorized)

	auth.POST("/api/campaigns").WithJSON(params.PostCampaign{Name: "", TemplateName: ""}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{"name": "This field is required", "template_name": "This field is required"})

	// test post campaign
	id := auth.POST("/api/campaigns").WithJSON(params.PostCampaign{Name: "foo1", TemplateName: templateName}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().Value("id")

	idStr := strconv.FormatFloat(id.Raw().(float64), 'f', 0, 64)

	auth.POST("/api/campaigns").WithJSON(params.PostCampaign{Name: "foo2", TemplateName: templateName}).
		Expect().
		Status(http.StatusCreated)

	auth.POST("/api/campaigns").WithJSON(params.PostCampaign{Name: "test-scopes", TemplateName: templateName}).
		Expect().
		Status(http.StatusForbidden).JSON().Object().
		ValueEqual("message", "You have exceeded your campaigns limit, please upgrade to a bigger plan or contact support.")

	// test scopes
	collection := auth.GET("/api/campaigns").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("total", 2)

	collection.Value("links").Object().ContainsKey("previous").ContainsKey("next")
	collection.Value("collection").Array().NotEmpty().Length().Equal(2)

	auth.GET("/api/campaigns").
		WithQuery("scopes[name]", "foo").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("total", 2)

	// test get campaign
	auth.GET("/api/campaigns/"+idStr).
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual("name", "foo1").
		ValueEqual("status", "draft")

	auth.PUT("/api/campaigns/" + idStr).WithJSON(params.PutCampaign{Name: "TESTputtest", TemplateName: templateName}).
		Expect().
		Status(http.StatusOK)

	// test updated campaign
	auth.GET("/api/campaigns/"+idStr).
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual("name", "TESTputtest").
		ValueEqual("status", "draft")

	// start campaign
	auth.POST("/api/campaigns/"+idStr+"/start").WithJSON(params.StartCampaign{}).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{
			"from_name":   "This field is required",
			"segment_ids": "This field is required",
			"source":      "This field is required",
		})

	// test campaign not found
	auth.POST("/api/campaigns/2223/start").WithJSON(params.StartCampaign{FromName: "gl", Source: "gudgl@example.com", SegmentIDs: []int64{1}}).
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().
		ValueEqual("message", "Campaign not found")

	// todo add ses keys (need mocking library methods)
	/*auth.POST("/api/ses/keys").WithJSON(params.PostSESKeys{
		AccessKey: "testAccessKey",
		SecretKey: "test secret key",
		Region:    "test region",
	}).Expect().
		Status(http.StatusOK)*/

	// test without ses keys
	auth.POST("/api/campaigns/"+idStr+"/start").WithJSON(params.StartCampaign{FromName: "gl", Source: "gudgl@example.com", SegmentIDs: []int64{1}}).
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().
		ValueEqual("message", "Amazon Ses keys are not set.")

	// successful patch campaign schedule.
	auth.PATCH("/api/campaigns/1/schedule").WithJSON(params.CampaignSchedule{
		FromName:    "gl",
		Source:      "gudgl@example.com",
		SegmentIDs:  []int64{1},
		ScheduledAt: "2020-04-04 15:04:03",
	}).
		Expect().
		Status(http.StatusOK).JSON().Object().
		ValueEqual("message", "Campaign TESTputtest successfully scheduled at 2020-04-04 15:04:03")

	// wrong time format patch campaign schedule.
	auth.PATCH("/api/campaigns/1/schedule").WithJSON(params.CampaignSchedule{}).
		Expect().
		Status(http.StatusBadRequest).JSON().Object().
		ValueEqual("message", "Invalid parameters, please try again").
		ValueEqual("errors", map[string]string{
			"from_name":    "This field is required",
			"scheduled_at": "This field is required",
			"segment_ids":  "This field is required",
			"source":       "This field is required",
		})

	// delete campaign by id
	auth.DELETE("/api/campaigns/" + idStr).
		Expect().
		Status(http.StatusNoContent)

}
