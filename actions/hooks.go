package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/emails"
	"github.com/news-maily/api/entities"
	"github.com/sirupsen/logrus"
)

func HandleHook(c *gin.Context) {
	var sns entities.SNSMessage

	body, err := c.GetRawData()
	err = json.Unmarshal(body, &sns)
	if err != nil {
		logrus.Errorf("Cannot decode request: %s", err.Error())
		return
	}

	if sns.Type == emails.SubConfirmationType {
		response, err := http.Get(sns.SubscribeURL)
		if err != nil {
			logrus.Errorf("AWS error while confirming the subscribe URL: %s", err.Error())
			return
		}

		defer response.Body.Close()

		if response.StatusCode >= http.StatusBadRequest {
			xml, _ := ioutil.ReadAll(response.Body)
			logrus.Errorf("AWS error while confirming the subscribe URL: %s", string(xml))
		} else {
			logrus.Infof("AWS SNS topic successfully subscribed: %s", sns.SubscribeURL)
		}

		return
	}

	var msg entities.SesMessage

	s, _ := strconv.Unquote(string(sns.RawMessage))

	logrus.WithFields(logrus.Fields{
		"type": sns.Type,
	}).Infof("SNS Raw Message")

	err = json.Unmarshal([]byte(s), &msg)
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
	}

	cid, err := strconv.ParseInt(cidTag[0], 10, 64)
	if err != nil {
		logrus.Errorf("unable to parse campaign id str to int: %s", err.Error())
	}

	// fetch the user id from tags
	uidTag, ok := msg.Mail.Tags["user_id"]
	if !ok || len(uidTag) == 0 {
		logrus.WithFields(logrus.Fields{
			"message_id": msg.Mail.MessageID,
			"source":     msg.Mail.Source,
			"tags":       msg.Mail.Tags,
		}).Error("user id not found in mail tags")
	}

	uid, err := strconv.ParseInt(uidTag[0], 10, 64)
	if err != nil {
		logrus.Errorf("unable to parse user id str to int: %s", err.Error())
	}

	logrus.WithFields(logrus.Fields{
		"type":        msg.NotificationType,
		"mail":        msg.Mail,
		"campaign_id": cid,
		"user_id":     uid,
	}).Infof("sns message")

	// todo: insert data into proper tables
	switch msg.NotificationType {
	case emails.BounceType:
	case emails.ComplaintType:
	case emails.DeliveryType:
	case emails.SendType:
	case emails.RenderingFailureType:
	case emails.ClickType:
	default:
		logrus.Errorf("Received unknown AWS SES message: %s", msg.NotificationType)
	}
}
