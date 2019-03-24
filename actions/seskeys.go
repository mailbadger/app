package actions

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/emails"
	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/events"
	"github.com/news-maily/api/routes/middleware"
	"github.com/news-maily/api/storage"
)

func GetSESKeys(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": "AWS Ses keys not set.",
		})
		return
	}

	keys.SecretKey = "" //do not return the secret

	c.JSON(http.StatusOK, keys)
}

func PostSESKeys(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.Id)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "AWS Ses keys are already set.",
		})
		return
	}

	keys = &entities.SesKeys{
		AccessKey: c.PostForm("access_key"),
		SecretKey: c.PostForm("secret_key"),
		Region:    c.PostForm("region"),
		UserId:    u.Id,
	}

	if !keys.Validate() {
		c.JSON(http.StatusBadRequest, keys.Errors)
		return
	}

	sender, err := emails.NewSesSender(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logrus.Errorln(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "SES keys are incorrect.",
		})
		return
	}

	snsClient, err := events.NewEventsClient(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logrus.Errorln(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "SES keys are incorrect.",
		})
		return
	}

	// Check if the configuration set is already created
	topicArn := ""
	cs, err := sender.DescribeConfigurationSet(&ses.DescribeConfigurationSetInput{
		ConfigurationSetName: aws.String(emails.ConfigurationSetName),
		ConfigurationSetAttributeNames: []*string{
			aws.String("eventDestinations"),
		},
	})

	if err != nil {
		_, err := sender.CreateConfigurationSet(&ses.CreateConfigurationSetInput{
			ConfigurationSet: &ses.ConfigurationSet{
				Name: aws.String(emails.ConfigurationSetName),
			},
		})
		if err != nil {
			logrus.Errorln(err)
			return
		}

		snsRes, err := snsClient.CreateTopic(&sns.CreateTopicInput{
			Name: aws.String(events.SNSTopicName),
		})
		if err != nil {
			// rollback
			sender.DeleteConfigurationSet(&ses.DeleteConfigurationSetInput{
				ConfigurationSetName: aws.String(emails.ConfigurationSetName),
			})
			logrus.Errorln(err)
			return
		}

		topicArn = *snsRes.TopicArn

		_, err = snsClient.Subscribe(&sns.SubscribeInput{
			Protocol: aws.String("https"),
			Endpoint: aws.String(os.Getenv("DOMAIN_URL") + "/api/hooks"),
			TopicArn: aws.String(topicArn),
		})
		if err != nil {
			// rollback
			sender.DeleteConfigurationSet(&ses.DeleteConfigurationSetInput{
				ConfigurationSetName: aws.String(emails.ConfigurationSetName),
			})
			snsClient.DeleteTopic(&sns.DeleteTopicInput{TopicArn: snsRes.TopicArn})
			logrus.Errorln(err)
			return
		}
	}

	// Check if the event destination is set
	eventFound := false
	for _, e := range cs.EventDestinations {
		if e.Name != nil && *e.Name == events.SNSTopicName {
			eventFound = true
		}
	}

	if !eventFound {
		if topicArn == "" {
			snsRes, err := snsClient.CreateTopic(&sns.CreateTopicInput{
				Name: aws.String(events.SNSTopicName),
			})
			if err != nil {
				logrus.Errorln(err)
				return
			}

			topicArn = *snsRes.TopicArn

			_, err = snsClient.Subscribe(&sns.SubscribeInput{
				Protocol: aws.String("https"),
				Endpoint: aws.String("https://88950e6a.ngrok.io/api/hooks"),
				TopicArn: aws.String(topicArn),
			})
			if err != nil {
				logrus.Errorln(err)
				return
			}
		}

		_, err = sender.CreateConfigurationSetEventDestination(&ses.CreateConfigurationSetEventDestinationInput{
			ConfigurationSetName: aws.String(emails.ConfigurationSetName),
			EventDestination: &ses.EventDestination{
				Name:    aws.String(events.SNSTopicName),
				Enabled: aws.Bool(true),
				MatchingEventTypes: []*string{
					aws.String("send"),
					aws.String("open"),
					aws.String("click"),
					aws.String("bounce"),
					aws.String("reject"),
					aws.String("delivery"),
					aws.String("complaint"),
					aws.String("renderingFailure"),
				},
				SNSDestination: &ses.SNSDestination{
					TopicARN: aws.String(topicArn),
				},
			},
		})

		if err != nil {
			logrus.Errorln(err)
		}
	}

	err = storage.CreateSesKeys(c, keys)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"input": *keys,
		}).Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Unable to add ses keys.",
		})
		return
	}

	keys.SecretKey = ""

	c.JSON(http.StatusCreated, keys)
}

func DeleteSESKeys(c *gin.Context) {
	u := middleware.GetUser(c)

	err := storage.DeleteSesKeys(c, u.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "AWS Ses keys not set.",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
