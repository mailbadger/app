package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/google/uuid"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/queue"
	"github.com/news-maily/api/routes/middleware"
	"github.com/news-maily/api/storage"
	"github.com/news-maily/api/utils/pagination"
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
			"reason": "Id must be an integer",
		})
		return
	}

	params := &sendCampaignParams{}
	c.Bind(params)

	v, err := valid.ValidateStruct(params)
	if err != nil || !v {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": err.Error(),
		})
		return
	}

	templateData := c.PostFormMap("default_template_data")

	u := middleware.GetUser(c)

	campaign, err := storage.GetCampaign(c, id, u.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": "Campaign not found",
		})
		return
	}

	if campaign.Status != entities.StatusDraft {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": fmt.Sprintf(`Campaign has a status of '%s', cannot start the campaign.`, campaign.Status),
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

	go func(
		userID int64,
		campaign *entities.Campaign,
		sesKeys *entities.SesKeys,
		params *sendCampaignParams,
		templateData map[string]string,
	) {
		// fetching subs that are active and that have not been blacklisted
		subs, err := storage.GetDistinctSubscribersByListIDs(c, params.Ids, userID, false, true)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id":  userID,
				"list_ids": params.Ids,
			}).Errorf("unable to fetch subscribers: %s", err.Error())
			return
		}

		// SES allows to send 50 emails in a bulk sending operation
		chunkSize := 50
		for i := 0; i < len(subs); i += chunkSize {
			end := i + chunkSize
			if end > len(subs) {
				end = len(subs)
			}

			// create
			var dest []*ses.BulkEmailDestination
			for _, s := range subs[i:end] {
				// marshal sub template data
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

			uuid, err := uuid.NewRandom()
			if err != nil {
				logrus.Errorf("unable to generate random uuid: %s", err.Error())
				continue
			}

			defaultData, err := json.Marshal(templateData)
			if err != nil {
				logrus.Errorln(err)
				continue
			}

			// prepare message for publishing to the queue
			msg, err := json.Marshal(entities.BulkSendMessage{
				UUID: uuid.String(),
				Input: &ses.SendBulkTemplatedEmailInput{
					Source:               aws.String(params.Source),
					Template:             aws.String(campaign.TemplateName),
					Destinations:         dest,
					ConfigurationSetName: aws.String("test"),
					DefaultTemplateData:  aws.String(string(defaultData)),
					DefaultTags: []*ses.MessageTag{
						&ses.MessageTag{
							Name:  aws.String("campaign_id"),
							Value: aws.String(strconv.Itoa(int(campaign.Id))),
						},
						&ses.MessageTag{
							Name:  aws.String("user_id"),
							Value: aws.String(strconv.Itoa(int(userID))),
						},
					},
				},
				CampaignID: campaign.Id,
				UserID:     u.Id,
				SesKeys:    sesKeys,
			})

			if err != nil {
				logrus.Errorln(err)
				continue
			}

			// publish the message to the queue
			err = queue.Publish(c, entities.CampaignsTopic, msg)
			if err != nil {
				logrus.Errorln(err)
			}
		}

		campaign.Status = entities.StatusSent
		err = storage.UpdateCampaign(c, campaign)
		if err != nil {
			logrus.WithField("campaign", campaign).Errorln(err)
			return
		}
	}(u.Id, campaign, sesKeys, params, templateData)

	c.JSON(http.StatusOK, gin.H{
		"reason": "The campaign has started. You can track the progress in the campaign details page.",
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
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
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
		Status:       entities.StatusDraft,
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
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
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
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
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
