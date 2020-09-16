package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/queue"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/templates"
)

type sendCampaignParams struct {
	Ids      []int64 `form:"segment_id[]" valid:"required"`
	Source   string  `form:"source" valid:"email,required~Email is empty or in invalid format"`
	FromName string  `form:"from_name" valid:"required,stringlength(1|191)~From name is blank or exceeds maximum character limit."`
}

func StartCampaign(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
		return
	}

	params := &sendCampaignParams{}
	err = c.Bind(params)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to bind send campaign params.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid parameters, please try again.",
		})
		return
	}

	v, err := valid.ValidateStruct(params)
	if !v {
		msg := "Unable to start campaign, invalid request parameters."
		if err != nil {
			msg = err.Error()
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": msg,
		})
		return
	}

	templateData := c.PostFormMap("default_template_data")

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

	sesKeys, err := storage.GetSesKeys(c, u.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Amazon Ses keys are not set.",
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

	store, err := templates.NewSesTemplateStore(sesKeys.AccessKey, sesKeys.SecretKey, sesKeys.Region)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to create SES template store.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	lists, err := storage.GetSegmentsByIDs(c, u.ID, params.Ids)
	if err != nil || len(lists) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Subscriber lists are not found.",
		})
		return
	}

	_, err = store.GetTemplate(&ses.GetTemplateInput{
		TemplateName: aws.String(campaign.TemplateName),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Template not found. Unable to send campaign.",
		})
		return
	}

	_, err = sender.DescribeConfigurationSet(&ses.DescribeConfigurationSetInput{
		ConfigurationSetName: aws.String(emails.ConfigurationSetName),
	})

	csExists := err == nil

	msg, err := json.Marshal(entities.SendCampaignParams{
		SegmentIDs:             params.Ids,
		Source:                 fmt.Sprintf("%s <%s>", params.FromName, params.Source),
		TemplateData:           templateData,
		UserID:                 u.ID,
		UserUUID:               u.UUID,
		Campaign:               *campaign,
		SesKeys:                *sesKeys,
		ConfigurationSetExists: csExists,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to publish campaign.",
		})
		return
	}

	err = queue.Publish(c, entities.CampaignsTopic, msg)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"campaign_id": campaign.ID,
			"segment_ids": params.Ids,
		}).WithError(err).Error("Unable to queue campaign for sending.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to publish campaign.",
		})
		return
	}

	campaign.Status = entities.StatusSending
	err = storage.UpdateCampaign(c, campaign)
	if err != nil {
		logger.From(c).
			WithField("campaign_id", campaign.ID).
			WithError(err).Error("Unable to update campaign status.")
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

	err := storage.GetCampaigns(c, middleware.GetUser(c).ID, p)
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
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		if campaign, err := storage.GetCampaign(c, id, middleware.GetUser(c).ID); err == nil {
			c.JSON(http.StatusOK, campaign)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign not found",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Id must be an integer",
	})
}

func PostCampaign(c *gin.Context) {
	name, templateName := strings.TrimSpace(c.PostForm("name")), strings.TrimSpace(c.PostForm("template_name"))
	user := middleware.GetUser(c)

	_, err := storage.GetCampaignByName(c, name, middleware.GetUser(c).ID)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Campaign with that name already exists",
		})
		return
	}

	campaign := &entities.Campaign{
		Name:         name,
		UserID:       user.ID,
		TemplateName: templateName,
		Status:       entities.StatusDraft,
	}

	if !campaign.Validate() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid data",
			"errors":  campaign.Errors,
		})
		return
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
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		user := middleware.GetUser(c)

		campaign, err := storage.GetCampaign(c, id, user.ID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Campaign not found",
			})
			return
		}

		name, templateName := strings.TrimSpace(c.PostForm("name")), strings.TrimSpace(c.PostForm("template_name"))

		campaign2, err := storage.GetCampaignByName(c, name, middleware.GetUser(c).ID)
		if err == nil && campaign.ID != campaign2.ID {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Campaign with that name already exists",
			})
			return
		}

		campaign.Name = name
		campaign.TemplateName = templateName

		if !campaign.Validate() {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Invalid data",
				"errors":  campaign.Errors,
			})
			return
		}

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
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Id must be an integer",
	})
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

	p, ok := val.(*storage.PaginationCursor)
	if !ok {
		logger.From(c).Error("Unable to cast pagination cursor from context value.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch campaign opens. Please try again.",
		})
		return
	}

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		if err := storage.GetCampaignOpens(c, id, p); err == nil {
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
	return

}
