package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/emails"
	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/storage"
	sns "github.com/robbiet480/go.sns"
	"github.com/sirupsen/logrus"
)

func HandleHook(c *gin.Context) {
	var payload sns.Payload

	body, err := c.GetRawData()
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

		defer response.Body.Close()

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

	// fetch the user id from tags
	uidTag, ok := msg.Mail.Tags["user_id"]
	if !ok || len(uidTag) == 0 {
		logrus.WithFields(logrus.Fields{
			"message_id": msg.Mail.MessageID,
			"source":     msg.Mail.Source,
			"tags":       msg.Mail.Tags,
		}).Error("user id not found in mail tags")
		return
	}

	uid, err := strconv.ParseInt(uidTag[0], 10, 64)
	if err != nil {
		logrus.Errorf("unable to parse user id str to int: %s", err.Error())
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
				UserID:         uid,
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
				err = storage.BlacklistSubscriber(c, uid, recipient.EmailAddress)
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
				UserID:     uid,
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
				UserID:               uid,
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
	case emails.ClickType:
		if msg.Click == nil {
			logrus.WithField("notif", msg).Errorln("click is empty")
			return
		}

		err := storage.CreateClick(c, &entities.Click{
			UserID:     uid,
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
			UserID:     uid,
			CampaignID: cid,
			UserAgent:  msg.Open.UserAgent,
			IPAddress:  msg.Open.IPAddress,
			CreatedAt:  msg.Open.Timestamp,
		})
		if err != nil {
			logrus.WithField("notif", msg).Errorln(err.Error())
		}
	case emails.RenderingFailureType:
	default:
		logrus.WithField("sns", msg).Error("unknown AWS SES message")
	}
}
