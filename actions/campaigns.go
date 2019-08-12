package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/queue"
	"github.com/news-maily/app/routes/middleware"
	"github.com/news-maily/app/storage"
	"github.com/news-maily/app/utils/pagination"
	"github.com/sirupsen/logrus"
)

type sendCampaignParams struct {
	Ids    []int64 `form:"list_id[]" valid:"required"`
	Source string  `form:"source" valid:"email,required~Email is blank or in invalid format"`
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
	c.Bind(params)

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

	lists, err := storage.GetListsByIDs(c, u.ID, params.Ids)
	if err != nil || len(lists) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Subscriber lists are not found.",
		})
		return
	}

	msg, err := json.Marshal(entities.SendCampaignParams{
		ListIDs:      params.Ids,
		Source:       params.Source,
		TemplateData: templateData,
		UserID:       u.ID,
		Campaign:     *campaign,
		SesKeys:      *sesKeys,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to publish campaign.",
		})
		return
	}

	err = queue.Publish(c, entities.CampaignsTopic, msg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"campaign_id": campaign.ID,
			"user_id":     u.ID,
			"list_ids":    params.Ids,
		}).WithError(err).Error("unable to queue campaign for sending")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to publish campaign.",
		})
		return
	}

	campaign.Status = entities.StatusSending
	err = storage.UpdateCampaign(c, campaign)
	if err != nil {
		logrus.WithField("campaign", campaign).WithError(err).Error("unable to update campaign status")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "The campaign has started. You can track the progress in the campaign details page.",
	})
}

func GetCampaigns(c *gin.Context) {
	val, ok := c.Get("pagination")
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("cannot create pagination object"))
		return
	}

	p, ok := val.(*pagination.Pagination)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("cannot cast pagination object"))
		return
	}

	storage.GetCampaigns(c, middleware.GetUser(c).ID, p)
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
	name, templateName := c.PostForm("name"), c.PostForm("template_name")
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
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

		name, templateName := c.PostForm("name"), c.PostForm("template_name")

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
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": err.Error(),
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
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": err.Error(),
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
