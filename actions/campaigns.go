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
	"github.com/news-maily/api/emails"
	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/routes/middleware"
	"github.com/news-maily/api/storage"
	"github.com/news-maily/api/utils/pagination"
	"github.com/sirupsen/logrus"
)

type listIds struct {
	Ids []int64 `form:"list_id[]"`
}

func StartCampaign(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Id must be an integer",
		})
		return
	}

	l := &listIds{}
	c.Bind(l)

	if len(l.Ids) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": "The list of ids is empty, cannot start the campaign.",
		})
		return
	}

	u := middleware.GetUser(c)

	campaign, err := storage.GetCampaign(c, id, u.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": "Campaign not found",
		})
		return
	}

	if campaign.Status != entities.STATUS_DRAFT {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf(`Campaign has a status of "%s", cannot start the campaign.`, campaign.Status),
		})
		return
	}

	sesKeys, err := storage.GetSesKeys(c, u.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": "Amazon Ses keys are not set.",
		})
		return
	}

	client, err := emails.NewSesSender(sesKeys.AccessKey, sesKeys.SecretKey, sesKeys.Region)
	if err != nil {
		logrus.Errorln(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "SES keys are incorrect.",
		})
		return
	}

	// fetching subs that are active and that have not been blacklisted
	subs, err := storage.GetDistinctSubscribersByListIDs(c, l.Ids, u.Id, false, true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": "Subscribers list is empty",
		})
		return
	}

	campaign.Status = entities.STATUS_SENDING
	err = storage.UpdateCampaign(c, campaign)
	if err != nil {
		logrus.Errorln(err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"reason": "Cannot update the campaign status, campaign sending is aborted.",
		})
		return
	}

	// SES allows to send 50 emails in a bulk sending operation
	chunkSize := 50
	for i := 0; i < len(subs); i += 50 {
		end := i + chunkSize

		if end > len(subs) {
			end = len(subs)
		}

		var dest []*ses.BulkEmailDestination
		for _, s := range subs[i:end] {
			s.Normalize()

			td, err := json.Marshal(s.TemplateData)
			if err != nil {
				logrus.Errorf("unable to marshal template data for subscriber %d - %s", s.Id, err.Error())
				continue
			}

			d := &ses.BulkEmailDestination{
				Destination: &ses.Destination{
					ToAddresses: []*string{aws.String(s.Email)},
				},
				ReplacementTemplateData: aws.String(string(td)),
			}

			dest = append(dest, d)
		}

		res, err := client.SendBulkTemplatedEmail(&ses.SendBulkTemplatedEmailInput{
			Source:       aws.String("me@filipnikolovski.com"),
			Template:     aws.String(campaign.TemplateName),
			Destinations: dest,
		})

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"campaign_id":   campaign.Id,
				"template_name": campaign.TemplateName,
			}).Errorln(err.Error())
			continue
		}

		for _, s := range res.Status {
			logrus.Info(s.GoString())
		}
	}

	campaign.Status = entities.STATUS_COMPLETED
	err = storage.UpdateCampaign(c, campaign)
	if err != nil {
		logrus.Errorln(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Cannot update the campaign status.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reason": "Campaign has started.",
	})
	return
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

	storage.GetCampaigns(c, middleware.GetUser(c).Id, p)
	c.JSON(http.StatusOK, p)
}

func GetCampaign(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		if campaign, err := storage.GetCampaign(c, id, middleware.GetUser(c).Id); err == nil {
			c.JSON(http.StatusOK, campaign)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"reason": "Campaign not found",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Id must be an integer",
	})
	return
}

func PostCampaign(c *gin.Context) {
	name, templateName := c.PostForm("name"), c.PostForm("template_name")
	user := middleware.GetUser(c)

	_, err := storage.GetCampaignByName(c, name, middleware.GetUser(c).Id)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": "Campaign with that name already exists",
		})
		return
	}

	campaign := &entities.Campaign{
		Name:         name,
		UserId:       user.Id,
		TemplateName: templateName,
		Status:       entities.STATUS_DRAFT,
	}

	if !campaign.Validate() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": "Invalid data",
			"errors": campaign.Errors,
		})
		return
	}

	err = storage.CreateCampaign(c, campaign)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, campaign)
	return
}

func PutCampaign(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		user := middleware.GetUser(c)

		campaign, err := storage.GetCampaign(c, id, user.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Campaign not found",
			})
			return
		}

		name, templateName := c.PostForm("name"), c.PostForm("template_name")

		campaign2, err := storage.GetCampaignByName(c, name, middleware.GetUser(c).Id)
		if err == nil && campaign.Id != campaign2.Id {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Campaign with that name already exists",
			})
			return
		}

		campaign.Name = name
		campaign.TemplateName = templateName

		if !campaign.Validate() {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Invalid data",
				"errors": campaign.Errors,
			})
			return
		}

		err = storage.UpdateCampaign(c, campaign)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": err.Error(),
			})
			return
		}

		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Id must be an integer",
	})
	return
}

func DeleteCampaign(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		user := middleware.GetUser(c)

		_, err := storage.GetCampaign(c, id, user.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Campaign not found",
			})
			return
		}

		err = storage.DeleteCampaign(c, id, user.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": err.Error(),
			})
			return
		}

		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Id must be an integer",
	})
	return
}
