package actions

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/templates"
	"github.com/mailbadger/app/validator"
)

func GetTemplate(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "AWS Ses keys not set.",
		})
		return
	}

	name := c.Param("name")

	store, err := templates.NewSesTemplateStore(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to create SES template store.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	res, err := store.GetTemplate(&ses.GetTemplateInput{
		TemplateName: aws.String(name),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Template not found.",
		})
		return
	}

	c.JSON(http.StatusOK, entities.Template{
		Name:        name,
		HTMLPart:    *res.Template.HtmlPart,
		TextPart:    *res.Template.TextPart,
		SubjectPart: *res.Template.SubjectPart,
	})
}

func GetTemplates(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "AWS Ses keys not set.",
		})
		return
	}

	nextToken := c.Query("next_token")

	store, err := templates.NewSesTemplateStore(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to create SES template store.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	res, err := store.ListTemplates(&ses.ListTemplatesInput{
		NextToken: aws.String(nextToken),
	})

	if err != nil {
		logger.From(c).WithField("token", nextToken).WithError(err).Error("Unable to list templates.")
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Templates not found, invalid page token.",
		})
		return
	}

	list := []entities.TemplateMeta{}

	for _, t := range res.TemplatesMetadata {
		list = append(list, entities.TemplateMeta{
			Name:      *t.Name,
			Timestamp: *t.CreatedTimestamp,
		})
	}

	var nt string
	if res.NextToken != nil {
		nt = *res.NextToken
	}

	c.JSON(http.StatusOK, entities.TemplateCollection{
		NextToken:  nt,
		Collection: list,
	})
}

func PostTemplate(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "AWS Ses keys not set.",
		})
		return
	}

	body := &params.PostTemplate{}
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters, please try again",
		})
		return
	}

	if err := validator.Validate(body); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	store, err := templates.NewSesTemplateStore(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to create SES template store.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	_, err = store.CreateTemplate(&ses.CreateTemplateInput{
		Template: &ses.Template{
			TemplateName: aws.String(body.Name),
			HtmlPart:     aws.String(body.Content),
			TextPart:     aws.String(body.Content),
			SubjectPart:  aws.String(body.Subject),
		},
	})
	if err != nil {
		reason := "Unable to create template."

		if awsErr, ok := err.(awserr.Error); ok {
			reason = fmt.Sprintf("%s: %s", awsErr.Code(), awsErr.Message())
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"message": reason,
		})
		return
	}

	c.JSON(http.StatusCreated, entities.Template{
		Name:        body.Name,
		HTMLPart:    body.Content,
		TextPart:    body.Content,
		SubjectPart: body.Subject,
	})
}

func PutTemplate(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "AWS Ses keys not set.",
		})
		return
	}

	body := &params.PutTemplate{}
	if err := c.ShouldBind(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid parameters, please try again",
		})
		return
	}

	if err := validator.Validate(body); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	name := c.Param("name")

	store, err := templates.NewSesTemplateStore(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to create SES template store.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}
	_, err = store.UpdateTemplate(&ses.UpdateTemplateInput{
		Template: &ses.Template{
			TemplateName: aws.String(name),
			HtmlPart:     aws.String(body.Content),
			TextPart:     aws.String(body.Content),
			SubjectPart:  aws.String(body.Subject),
		},
	})

	if err != nil {
		reason := "Unable to update template."

		if awsErr, ok := err.(awserr.Error); ok {
			reason = fmt.Sprintf("%s: %s", awsErr.Code(), awsErr.Message())
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"message": reason,
		})
		return
	}

	c.JSON(http.StatusOK, entities.Template{
		Name:        name,
		HTMLPart:    body.Content,
		TextPart:    body.Content,
		SubjectPart: body.Subject,
	})
}

// TODO change it after merging get template
func DeleteTemplate(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "AWS Ses keys not set.",
		})
		return
	}

	store, err := templates.NewSesTemplateStore(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logger.From(c).WithError(err).Error("Unable to create SES template store.")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	name := c.Param("name")

	_, err = store.DeleteTemplate(&ses.DeleteTemplateInput{
		TemplateName: aws.String(name),
	})

	if err != nil {
		reason := "Unable to delete template."

		if awsErr, ok := err.(awserr.Error); ok {
			reason = fmt.Sprintf("%s: %s", awsErr.Code(), awsErr.Message())
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"message": reason,
		})
		return
	}

	c.Status(http.StatusNoContent)
}
