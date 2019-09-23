package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/emails"
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/storage"
	sns "github.com/robbiet480/go.sns"
	"github.com/sirupsen/logrus"
)

func HandleHook(c *gin.Context) {
	var payload sns.Payload

	body, err := c.GetRawData()
	if err != nil {
		logrus.Errorf("Cannot fetch raw data: %s", err.Error())
		return
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		logrus.Errorf("Cannot decode request: %s", err.Error())
		return
	}

	err = payload.VerifyPayload()
	if err != nil {
		logrus.Error(err)
		return
	}

	if payload.Type == emails.SubConfirmationType {
		response, err := http.Get(payload.SubscribeURL)
		if err != nil {
			logrus.Errorf("AWS error while confirming the subscribe URL: %s", err.Error())
			return
		}

		defer func() {
			err := response.Body.Close()
			if err != nil {
				logrus.WithError(err).Error("Unable to close response body.")
			}
		}()

		if response.StatusCode >= http.StatusBadRequest {
			xml, _ := ioutil.ReadAll(response.Body)
			logrus.Errorf("AWS error while confirming the subscribe URL: %s", string(xml))
		} else {
			logrus.Infof("AWS SNS topic successfully subscribed: %s", payload.SubscribeURL)
		}

		return
	}

	var msg entities.SesMessage

	err = json.Unmarshal([]byte(payload.Message), &msg)
	if err != nil {
		logrus.Errorf("Cannot unmarshal SNS raw message: %s", err.Error())
		return
	}

	// fetch the campaign id from tags
	cidTag, ok := msg.Mail.Tags["campaign_id"]
	if !ok || len(cidTag) == 0 {
		logrus.WithFields(logrus.Fields{
			"message_id": msg.Mail.MessageID,
			"source":     msg.Mail.Source,
			"tags":       msg.Mail.Tags,
		}).Error("campaign id not found in mail tags")
		return
	}

	cid, err := strconv.ParseInt(cidTag[0], 10, 64)
	if err != nil {
		logrus.Errorf("unable to parse campaign id str to int: %s", err.Error())
		return
	}

	uuid := c.Param("uuid")
	u, err := storage.GetUserByUUID(c, uuid)
	if err != nil {
		logrus.WithField("uuid", uuid).WithError(err).Error("unable to fetch user")
		return
	}

	switch msg.NotificationType {
	case emails.BounceType:
		if msg.Bounce == nil {
			logrus.WithField("notif", msg).Errorln("bounce is empty")
			return
		}

		for _, recipient := range msg.Bounce.BouncedRecipients {
			err := storage.CreateBounce(c, &entities.Bounce{
				UserID:         u.ID,
				CampaignID:     cid,
				Recipient:      recipient.EmailAddress,
				Action:         recipient.Action,
				Status:         recipient.Status,
				DiagnosticCode: recipient.DiagnosticCode,
				Type:           msg.Bounce.BounceType,
				SubType:        msg.Bounce.BounceSubType,
				FeedbackID:     msg.Bounce.FeedbackID,
				CreatedAt:      msg.Bounce.Timestamp,
			})
			if err != nil {
				logrus.WithField("notif", msg).Errorln(err.Error())
			}

			if msg.Bounce.BounceType == "Permanent" {
				err = storage.BlacklistSubscriber(c, u.ID, recipient.EmailAddress)
				if err != nil {
					logrus.WithField("notif", msg).Errorln(err.Error())
				}
			}
		}
	case emails.ComplaintType:
		if msg.Complaint == nil {
			logrus.WithField("notif", msg).Errorln("complaint is empty")
			return
		}

		for _, recipient := range msg.Complaint.ComplainedRecipients {
			err := storage.CreateComplaint(c, &entities.Complaint{
				UserID:     u.ID,
				CampaignID: cid,
				Recipient:  recipient.EmailAddress,
				Type:       msg.Complaint.ComplaintFeedbackType,
				FeedbackID: msg.Complaint.FeedbackID,
				CreatedAt:  msg.Complaint.Timestamp,
			})
			if err != nil {
				logrus.WithField("notif", msg).Errorln(err.Error())
			}
		}
	case emails.DeliveryType:
		if msg.Delivery == nil {
			logrus.WithField("notif", msg).Errorln("delivery is empty")
			return
		}

		for _, r := range msg.Delivery.Recipients {
			err := storage.CreateDelivery(c, &entities.Delivery{
				UserID:               u.ID,
				CampaignID:           cid,
				Recipient:            r,
				ProcessingTimeMillis: msg.Delivery.ProcessingTimeMillis,
				ReportingMTA:         msg.Delivery.ReportingMTA,
				RemoteMtaIP:          msg.Delivery.RemoteMtaIP,
				SMTPResponse:         msg.Delivery.SMTPResponse,
				CreatedAt:            msg.Delivery.Timestamp,
			})
			if err != nil {
				logrus.WithField("notif", msg).Errorln(err.Error())
			}
		}
	case emails.SendType:
		for _, d := range msg.Mail.Destination {
			err := storage.CreateSend(c, &entities.Send{
				UserID:           u.ID,
				CampaignID:       cid,
				MessageID:        msg.Mail.MessageID,
				Source:           msg.Mail.Source,
				SendingAccountID: msg.Mail.SendingAccountID,
				Destination:      d,
				CreatedAt:        msg.Mail.Timestamp,
			})
			if err != nil {
				logrus.WithField("notif", msg).Errorln(err.Error())
			}
		}
	case emails.ClickType:
		if msg.Click == nil {
			logrus.WithField("notif", msg).Errorln("click is empty")
			return
		}

		err := storage.CreateClick(c, &entities.Click{
			UserID:     u.ID,
			CampaignID: cid,
			Link:       msg.Click.Link,
			UserAgent:  msg.Click.UserAgent,
			IPAddress:  msg.Click.IPAddress,
			CreatedAt:  msg.Click.Timestamp,
		})
		if err != nil {
			logrus.WithField("notif", msg).Errorln(err.Error())
		}
	case emails.OpenType:
		if msg.Open == nil {
			logrus.WithField("notif", msg).Errorln("open is empty")
			return
		}

		err := storage.CreateOpen(c, &entities.Open{
			UserID:     u.ID,
			CampaignID: cid,
			UserAgent:  msg.Open.UserAgent,
			IPAddress:  msg.Open.IPAddress,
			CreatedAt:  msg.Open.Timestamp,
		})
		if err != nil {
			logrus.WithField("notif", msg).Errorln(err.Error())
		}
	case emails.RenderingFailureType:
		logrus.WithFields(logrus.Fields{
			"campaign_id":   cid,
			"user_id":       u.ID,
			"error":         msg.RenderingFailure.ErrorMessage,
			"template_name": msg.RenderingFailure.TemplateName,
		}).Warn("rendering failure")
	default:
		logrus.WithField("sns", msg).Error("unknown AWS SES message")
	}
}
