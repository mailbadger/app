package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/queue"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/validator"
)

func StartCampaign(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
		return
	}

	body := &params.SendCampaign{}
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters, please try again",
		})
		return
	}

	// should bind supports only struct type so we need to take our map key value with PostFormMap before validating struct
	body.DefaultTemplateData = c.PostFormMap("default_template_data")

	if err := validator.Validate(body); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	u := middleware.GetUser(c)

	campaign, err := storage.GetCampaign(c, id, u.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign not found",
		})
		return
	}

	if campaign.Status != entities.StatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf(`Campaign has a status of '%s', cannot start the campaign.`, campaign.Status),
		})
		return
	}

	template, err := storage.GetTemplate(c, campaign.BaseTemplate.ID, u.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Template not found. Unable to send campaign.",
		})
		return
	}

	err = template.ValidateData(body.DefaultTemplateData)
	if err != nil {
		if errors.Is(err, entities.ErrMissingDefaultData) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Incomplete default template data. Unable to send campaign.",
			})
			return
		}
		logger.From(c).WithFields(logrus.Fields{
			"campaign_id": id,
			"user_id":     u.ID,
			"template_id": campaign.BaseTemplate.ID,
			"segment_ids": body.SegmentIDs,
		}).WithError(err).Warn("Unable to parse template")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Failed to parse template. Unable to send campaign.",
		})
		return
	}

	sesKeys, err := storage.GetSesKeys(c, u.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Amazon Ses keys are not set.",
		})
		return
	}

	lists, err := storage.GetSegmentsByIDs(c, u.ID, body.SegmentIDs)
	if err != nil || len(lists) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Subscriber lists are not found.",
		})
		return
	}

	sender, err := emails.NewSesSender(sesKeys.AccessKey, sesKeys.SecretKey, sesKeys.Region)
	if err != nil {
		logger.From(c).WithError(err).Warn("Unable to create SES sender.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	_, err = sender.DescribeConfigurationSet(&ses.DescribeConfigurationSetInput{
		ConfigurationSetName: aws.String(emails.ConfigurationSetName),
	})

	msg, err := json.Marshal(entities.CampaignerTopicParams{
		CampaignID:             id,
		SegmentIDs:             body.SegmentIDs,
		Source:                 fmt.Sprintf("%s <%s>", body.FromName, body.Source),
		TemplateData:           body.DefaultTemplateData,
		UserID:                 u.ID,
		UserUUID:               u.UUID,
		SesKeys:                *sesKeys,
		ConfigurationSetExists: err == nil,
	})
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"campaign_id": id,
			"user_id":     u.ID,
			"template_id": campaign.BaseTemplate.ID,
			"segment_ids": body.SegmentIDs,
		}).WithError(err).Error("Unable to marshal campaigner message body")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to publish campaign.",
		})
		return
	}

	err = queue.Publish(c, entities.CampaignerTopic, msg)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"campaign_id": id,
			"user_id":     u.ID,
			"template_id": campaign.BaseTemplate.ID,
			"segment_ids": body.SegmentIDs,
		}).WithError(err).Error("Unable to queue campaign for sending.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to publish campaign.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "The campaign has started. You can track the progress in the campaign details page.",
	})
}

func GetCampaigns(c *gin.Context) {
	val, ok := c.Get("cursor")
	if !ok {
		logger.From(c).Error("Unable to fetch pagination cursor from context.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch campaigns. Please try again.",
		})
		return
	}

	p, ok := val.(*storage.PaginationCursor)
	if !ok {
		logger.From(c).Error("Unable to cast pagination cursor from context value.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch campaigns. Please try again.",
		})
		return
	}

	scopeMap := c.QueryMap("scopes")
	err := storage.GetCampaigns(c, middleware.GetUser(c).ID, p, scopeMap)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"starting_after": p.StartingAfter,
			"ending_before":  p.EndingBefore,
		}).WithError(err).Warn("Unable to fetch campaigns collection.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch campaigns. Please try again.",
		})
		return
	}

	c.JSON(http.StatusOK, p)
}

func GetCampaign(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
		return
	}
	campaign, err := storage.GetCampaign(c, id, middleware.GetUser(c).ID)
	if err != nil {
		logrus.Info(err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign not found",
		})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

func PostCampaign(c *gin.Context) {
	body := &params.Campaign{}
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters, please try again",
		})
		return
	}

	if err := validator.Validate(body); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user := middleware.GetUser(c)

	_, err := storage.GetCampaignByName(c, body.Name, user.ID)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Campaign with that name already exists",
		})
		return
	}

	template, err := storage.GetTemplateByName(c, body.TemplateName, user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Template with that name does not exists",
		})
		return
	}

	campaign := &entities.Campaign{
		Name:         body.Name,
		UserID:       user.ID,
		BaseTemplate: template.GetBase(),
		Status:       entities.StatusDraft,
	}

	err = storage.CreateCampaign(c, campaign)
	if err != nil {
		logger.From(c).WithError(err).Warn("Unable to create campaign.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to create the campaign.",
		})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

func PutCampaign(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
		return
	}

	user := middleware.GetUser(c)

	campaign, err := storage.GetCampaign(c, id, user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Campaign not found",
		})
		return
	}

	body := &params.Campaign{}
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters, please try again",
		})
		return
	}

	if err := validator.Validate(body); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	campaign2, err := storage.GetCampaignByName(c, body.Name, middleware.GetUser(c).ID)
	if err == nil && campaign.ID != campaign2.ID {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Campaign with that name already exists",
		})
		return
	}

	template, err := storage.GetTemplateByName(c, body.TemplateName, user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Template with that name does not exists",
		})
		return
	}

	campaign.Name = body.Name
	campaign.BaseTemplate = template.GetBase()

	err = storage.UpdateCampaign(c, campaign)
	if err != nil {
		logger.From(c).
			WithError(err).
			WithField("campaign_id", id).
			Warn("Unable to update campaign.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to update campaign.",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteCampaign(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		user := middleware.GetUser(c)

		_, err := storage.GetCampaign(c, id, user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Campaign not found",
			})
			return
		}

		err = storage.DeleteCampaign(c, id, user.ID)
		if err != nil {
			logger.From(c).WithError(err).Warn("Unable to delete campaign.")
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Unable to delete the campaign.",
			})
			return
		}

		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Id must be an integer",
	})
}

func GetCampaignOpens(c *gin.Context) {
	val, ok := c.Get("cursor")
	if !ok {
		logger.From(c).Error("Unable to fetch pagination cursor from context.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch campaigns. Please try again.",
		})
		return
	}
	user := middleware.GetUser(c)

	p, ok := val.(*storage.PaginationCursor)
	if !ok {
		logger.From(c).Error("Unable to cast pagination cursor from context value.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch campaign opens. Please try again.",
		})
		return
	}

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		if err := storage.GetCampaignOpens(c, id, user.ID, p); err == nil {
			c.JSON(http.StatusOK, p)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign opens not found",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Id must be an integer",
	})
}

func GetCampaignStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
		return
	}
	user := middleware.GetUser(c)

	var campaignStats entities.CampaignStats
	campaignStats.TotalSent, err = storage.GetTotalSends(c, id, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign stats not found",
		})
		return
	}
	campaignStats.Delivered, err = storage.GetTotalDelivered(c, id, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign stats not found",
		})
		return
	}
	campaignStats.Opens, err = storage.GetOpensStats(c, id, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign stats not found",
		})
		return
	}
	campaignStats.Clicks, err = storage.GetClicksStats(c, id, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign stats not found",
		})
		return
	}
	campaignStats.Bounces, err = storage.GetTotalBounces(c, id, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign stats not found",
		})
		return
	}
	campaignStats.Complaints, err = storage.GetTotalComplaints(c, id, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign stats not found",
		})
		return
	}

	c.JSON(http.StatusOK, campaignStats)

}

func GetCampaignClicksStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
		return
	}

	stats, err := storage.GetCampaignClicksStats(c, id, middleware.GetUser(c).ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign clicks not found",
		})
		return
	}

	c.JSON(http.StatusOK, entities.CampaignClicksStats{
		Total:       int64(len(stats)),
		ClicksStats: stats,
	})
}

func GetCampaignComplaints(c *gin.Context) {
	val, ok := c.Get("cursor")
	if !ok {
		logger.From(c).Error("Unable to fetch pagination cursor from context.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch campaign complaints. Please try again.",
		})
		return
	}
	user := middleware.GetUser(c)

	p, ok := val.(*storage.PaginationCursor)
	if !ok {
		logger.From(c).Error("Unable to cast pagination cursor from context value.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch campaign complaints. Please try again.",
		})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
	}

	if err := storage.GetCampaignComplaints(c, id, user.ID, p); err == nil {
		c.JSON(http.StatusOK, p)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Campaign complaints not found",
	})
}

func GetCampaignBounces(c *gin.Context) {
	val, ok := c.Get("cursor")
	if !ok {
		logger.From(c).Error("Unable to fetch pagination cursor from context.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch campaigns bounces. Please try again.",
		})
		return
	}
	user := middleware.GetUser(c)

	p, ok := val.(*storage.PaginationCursor)
	if !ok {
		logger.From(c).Error("Unable to cast pagination cursor from context value.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch campaign bounces. Please try again.",
		})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
	}

	if err := storage.GetCampaignBounces(c, id, user.ID, p); err == nil {
		c.JSON(http.StatusOK, p)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Campaign bounces not found",
	})
}
