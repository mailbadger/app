package actions

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/emails"
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/events"
	"github.com/news-maily/app/logger"
	"github.com/news-maily/app/routes/middleware"
	"github.com/news-maily/app/storage"
)

func GetSESKeys(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "AWS Ses keys not set.",
		})
		return
	}

	keys.SecretKey = "" //do not return the secret

	c.JSON(http.StatusOK, keys)
}

func PostSESKeys(c *gin.Context) {
	u := middleware.GetUser(c)

	_, err := storage.GetSesKeys(c, u.ID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "AWS Ses keys are already set.",
		})
		return
	}

	keys := &entities.SesKeys{
		AccessKey: strings.TrimSpace(c.PostForm("access_key")),
		SecretKey: strings.TrimSpace(c.PostForm("secret_key")),
		Region:    strings.TrimSpace(c.PostForm("region")),
		UserID:    u.ID,
	}

	if !keys.Validate() {
		c.JSON(http.StatusBadRequest, keys.Errors)
		return
	}

	sender, err := emails.NewSesSender(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logger.From(c).WithError(err).Warn("Unable to create SES sender.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	snsClient, err := events.NewEventsClient(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logger.From(c).WithError(err).Warn("Unable to create SNS client.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	//TODO: Move this to work queue to be more robust.
	//createAWSResources is a slow process and could fail periodically.
	go func(
		c *gin.Context,
		sender emails.Sender,
		snsClient events.EventsClient,
		keys *entities.SesKeys,
		uuid string,
	) {
		err := createAWSResources(sender, snsClient, uuid)
		if err != nil {
			logger.From(c).WithError(err).Warn("Unable to create AWS resources.")
			return
		}

		err = storage.CreateSesKeys(c, keys)
		if err != nil {
			logger.From(c).WithError(err).Warn("Unable to create SES keys.")
		}
	}(c.Copy(), sender, snsClient, keys, u.UUID)

	c.JSON(http.StatusOK, gin.H{
		"message": "We are currently processing the request.",
	})
}

func createAWSResources(
	sender emails.Sender,
	snsClient events.EventsClient,
	uuid string,
) error {
	hookURL := fmt.Sprintf("%s/api/hooks/%s", os.Getenv("APP_URL"), uuid)

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
			return fmt.Errorf("ses keys: unable to create configuration set: %w", err)
		}

		snsRes, err := snsClient.CreateTopic(&sns.CreateTopicInput{
			Name: aws.String(events.SNSTopicName),
		})
		if err != nil {
			return fmt.Errorf("ses keys: unable to create SNS topic: %w", err)
		}

		topicArn = *snsRes.TopicArn

		_, err = snsClient.Subscribe(&sns.SubscribeInput{
			Protocol: aws.String("https"),
			Endpoint: aws.String(hookURL),
			TopicArn: aws.String(topicArn),
		})
		if err != nil {
			return fmt.Errorf("ses keys: unable to subscribe to topic: %w", err)
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
				return fmt.Errorf("ses keys: unable to create SNS topic: %w", err)
			}

			topicArn = *snsRes.TopicArn

			_, err = snsClient.Subscribe(&sns.SubscribeInput{
				Protocol: aws.String("https"),
				Endpoint: aws.String(hookURL),
				TopicArn: aws.String(topicArn),
			})
			if err != nil {
				return fmt.Errorf("ses keys: unable to subscribe to topic: %w", err)
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
			return fmt.Errorf("ses keys: unable to set event destination: %w", err)
		}
	}

	return nil
}

func DeleteSESKeys(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.ID)
	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}

	sender, err := emails.NewSesSender(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logger.From(c).WithError(err).Warn("Unable to create SES sender.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	snsClient, err := events.NewEventsClient(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logger.From(c).WithError(err).Warn("Unable to create SNS client.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	cs, err := sender.DescribeConfigurationSet(&ses.DescribeConfigurationSetInput{
		ConfigurationSetName: aws.String(emails.ConfigurationSetName),
		ConfigurationSetAttributeNames: []*string{
			aws.String("eventDestinations"),
		},
	})

	if err == nil {
		// find and delete topic
		for _, e := range cs.EventDestinations {
			if e.Name != nil && *e.Name == events.SNSTopicName {
				_, err := snsClient.DeleteTopic(&sns.DeleteTopicInput{
					TopicArn: e.SNSDestination.TopicARN,
				})
				if err != nil {
					logger.From(c).WithError(err).Warn("Unable to delete topic.")
				}
				break
			}
		}

		_, err = sender.DeleteConfigurationSet(&ses.DeleteConfigurationSetInput{
			ConfigurationSetName: aws.String(emails.ConfigurationSetName),
		})
		if err != nil {
			logger.From(c).WithError(err).Warn("Unable to delete configuration set.")
		}
	}

	err = storage.DeleteSesKeys(c, u.ID)
	if err != nil {
		logger.From(c).WithError(err).Warn("Unable to delete SES keys.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to delete SES keys.",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
