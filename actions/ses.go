package actions

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/events"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/validator"
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

	body := &params.PostSESKeys{}
	if err := c.ShouldBind(body); err != nil {
		logger.From(c).WithError(err).Error("Unable to bind ses keys params.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters, please try again",
		})
		return
	}

	if err := validator.Validate(body); err != nil {
		logger.From(c).WithError(err).Error("Invalid ses keys params.")
		c.JSON(http.StatusBadRequest, err)
		return
	}

	keys := &entities.SesKeys{
		AccessKey: body.AccessKey,
		SecretKey: body.SecretKey,
		Region:    body.Region,
		UserID:    u.ID,
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

	snsRes, err := snsClient.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(events.SNSTopicName),
	})
	if err != nil {
		return fmt.Errorf("ses keys: unable to create SNS topic: %w", err)
	}

	topicArn := *snsRes.TopicArn

	_, err = snsClient.Subscribe(&sns.SubscribeInput{
		Protocol: aws.String("https"),
		Endpoint: aws.String(hookURL),
		TopicArn: aws.String(topicArn),
	})
	if err != nil {
		return fmt.Errorf("ses keys: unable to subscribe to topic: %w", err)
	}

	// Check if the configuration set is already created
	cs, err := sender.DescribeConfigurationSet(&ses.DescribeConfigurationSetInput{
		ConfigurationSetName: aws.String(emails.ConfigurationSetName),
		ConfigurationSetAttributeNames: []*string{
			aws.String("eventDestinations"),
		},
	})

	if err != nil {
		_, err = sender.CreateConfigurationSet(&ses.CreateConfigurationSetInput{
			ConfigurationSet: &ses.ConfigurationSet{
				Name: aws.String(emails.ConfigurationSetName),
			},
		})
		if err != nil {
			return fmt.Errorf("ses keys: unable to create configuration set: %w", err)
		}
	}

	// Check if the event destination is set
	eventFound := false
	for _, e := range cs.EventDestinations {
		if e.Name != nil && *e.Name == events.SNSTopicName {
			eventFound = true
		}
	}

	if eventFound {
		return nil
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

func GetSESQuota(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "AWS Ses keys not set.",
		})
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

	res, err := sender.GetSendQuota(&ses.GetSendQuotaInput{})
	if err != nil {
		logger.From(c).WithError(err).Warn("Unable to fetch send quota.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to fetch send quota.",
		})
		return
	}

	c.JSON(http.StatusOK, entities.SendQuota{
		Max24HourSend:   *res.Max24HourSend,
		MaxSendRate:     *res.MaxSendRate,
		SentLast24Hours: *res.SentLast24Hours,
	})
}
