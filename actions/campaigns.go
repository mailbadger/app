package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/queue"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/services/boundaries"
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

	body := &params.StartCampaign{}
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

	if campaign.Status != entities.StatusDraft && campaign.Status != entities.StatusScheduled {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Campaign can not be started its already processed",
		})
		return
	}

	campaign.Status = entities.StatusSending
	campaign.SetEventID()

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
		EventID:                *campaign.EventID, // this id is handled in campaigns SetEventID method
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
			"template_id": campaign.BaseTemplate.ID,
			"segment_ids": body.SegmentIDs,
		}).WithError(err).Error("Unable to marshal campaigner message body")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to start campaign.",
		})
		return
	}

	err = queue.Publish(c, entities.CampaignerTopic, msg)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"campaign_id": id,
			"template_id": campaign.BaseTemplate.ID,
			"segment_ids": body.SegmentIDs,
		}).WithError(err).Error("Unable to queue campaign for sending.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to start campaign.",
		})
		return
	}

	err = storage.UpdateCampaign(c, campaign)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"campaign_id": id,
			"template_id": campaign.BaseTemplate.ID,
			"segment_ids": body.SegmentIDs,
		}).WithError(err).Error("Unable to update campaign.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to start campaign.",
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
	body := &params.PostCampaign{}
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
	boundariesvc := boundaries.New(storage.GetFromContext(c))

	limitexceeded, err := boundariesvc.CampaignsLimitExceeded(user)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to check campaigns limit for user.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to check campaigns limit. Please try again.",
		})
		return
	}

	if limitexceeded {
		logger.From(c).Info("User has exceeded his campaigns limit.")
		c.JSON(http.StatusForbidden, gin.H{
			"message": "You have exceeded your campaigns limit, please upgrade to a bigger plan or contact support.",
		})
		return
	}

	_, err = storage.GetCampaignByName(c, body.Name, user.ID)
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
		logger.From(c).WithError(err).Error("Unable to create campaign.")
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

	body := &params.PutCampaign{}
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

	c.JSON(http.StatusOK, campaign)
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

func DeleteCampaignSchedule(c *gin.Context) {
	u := middleware.GetUser(c)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer.",
		})
		return
	}
	campaign, err := storage.GetCampaign(c, id, u.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Campaign not found, please try again.",
		})
		return
	}

	err = storage.DeleteCampaignSchedule(c, campaign.ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"campaign_id": campaign.ID,
			"user_id":     u.ID,
		}).WithError(err).Error("unable to delete campaign schedule")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to delete campaign, please try again.",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func PatchCampaignSchedule(c *gin.Context) {

	u := middleware.GetUser(c)

	if !u.Boundaries.ScheduleCampaignsEnabled {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "You do not have permission to schedule campaign, please upgrade to a bigger plan or contact support.",
		})
	}

	campaignID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
		return
	}

	campaign, err := storage.GetCampaign(c, campaignID, u.ID)
	if err != nil {
		logrus.Println(err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign not found",
		})
		return
	}

	body := &params.CampaignSchedule{}
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

	schAt, err := time.Parse("2006-01-02 15:04:05", body.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters, scheduled_at should be format: 2006-02-01 15:04:05",
		})
		return
	}

	defMetadata, err := json.Marshal(body.DefaultTemplateData)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to schedule campaign, invalid default metdata.",
		})
		return
	}

	segmentIDsJSON, err := json.Marshal(body.SegmentIDs)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to schedule campaign, please try again.",
		})
		return
	}

	// if schedule exist update.
	if campaign.Schedule != nil {
		campaign.Schedule.ScheduledAt = schAt
		campaign.Schedule.FromName = body.FromName
		campaign.Schedule.Source = body.Source
		campaign.Schedule.SegmentIDsJSON = segmentIDsJSON
		campaign.Schedule.DefaultTemplateDataJSON = defMetadata
	} else {
		// else create new campaign schedule
		campaign.Schedule = &entities.CampaignSchedule{
			ID:                      ksuid.New(),
			CampaignID:              campaign.ID,
			ScheduledAt:             schAt,
			UserID:                  u.ID,
			SegmentIDsJSON:          segmentIDsJSON,
			FromName:                body.FromName,
			Source:                  body.Source,
			DefaultTemplateDataJSON: defMetadata,
		}
	}

	err = storage.CreateCampaignSchedule(c, campaign.Schedule)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"schedule_id": campaign.Schedule.ID,
			"campaign_id": campaign.Schedule.CampaignID,
			"user_id":     u.ID,
		}).WithError(err).Error("unable to create campaign schedule")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to patch scheduled campaign, please try again.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Campaign %s successfully scheduled at %v", campaign.Name, body.ScheduledAt),
	})

}
