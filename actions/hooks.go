package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/entities"
	"github.com/sirupsen/logrus"
)

// SES Notification Types
const (
	SendType             = "Send"
	ClickType            = "Click"
	BounceType           = "Bounce"
	DeliveryType         = "Delivery"
	ComplaintType        = "Complaint"
	RenderingFailureType = "Rendering Failure"
	SubConfirmationType  = "SubscriptionConfirmation"
)

func HandleHook(c *gin.Context) {
	var sns entities.SNSMessage

	body, err := c.GetRawData()
	err = json.Unmarshal(body, &sns)
	if err != nil {
		logrus.Errorf("Cannot decode request: %s", err.Error())
		return
	}

	if sns.Type == SubConfirmationType {
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

	var notification entities.SesMessage

	s, _ := strconv.Unquote(string(sns.RawMessage))

	err = json.Unmarshal([]byte(s), &notification)
	if err != nil {
		logrus.Errorf("Cannot unmarshal SNS raw message: %s", err.Error())
		return
	}

	logrus.Info(s)

	switch notification.NotificationType {
	case BounceType:
		logrus.Infof("Received SES bounce: %+v %+v", notification.Mail, notification.Bounce)
	case ComplaintType:
		logrus.Infof("Received SES complaint: %+v", notification.Complaint)
	case DeliveryType:
		logrus.Infof("Received SES delivery: %+v", notification.Delivery)
	case SendType:
		logrus.Infof("Received SES send: %+v", notification.Mail)
	case RenderingFailureType:
		logrus.Infof("Received SES rendering failure: %+v", notification.RenderingFailure)
	case ClickType:
		logrus.Infof("Received SES click: %+v", notification.Click)
	default:
		logrus.Errorf("Received unknown AWS SES message: %s", notification.NotificationType)
	}
}
