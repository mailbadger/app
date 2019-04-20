package actions

import (
	"fmt"
	"net/http"

	"github.com/news-maily/api/storage/templates"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/routes/middleware"
	"github.com/news-maily/api/storage"
)

func GetTemplate(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "AWS Ses keys not set.",
		})
		return
	}

	name := c.Param("name")

	store, err := templates.NewSesTemplateStore(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logrus.Errorln(err.Error())
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

	keys, err := storage.GetSesKeys(c, u.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "AWS Ses keys not set.",
		})
		return
	}

	nextToken := c.Query("next_token")

	store, err := templates.NewSesTemplateStore(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logrus.Errorln(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	res, err := store.ListTemplates(&ses.ListTemplatesInput{
		NextToken: aws.String(nextToken),
	})

	if err != nil {
		logrus.WithField("token", nextToken).Error(err)
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
		NextToken: nt,
		List:      list,
	})
}

func PostTemplate(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "AWS Ses keys not set.",
		})
		return
	}

	store, err := templates.NewSesTemplateStore(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logrus.Errorln(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	name := c.PostForm("name")
	html := c.PostForm("content")
	subject := c.PostForm("subject")

	_, err = store.CreateTemplate(&ses.CreateTemplateInput{
		Template: &ses.Template{
			TemplateName: aws.String(name),
			HtmlPart:     aws.String(html),
			TextPart:     aws.String(html),
			SubjectPart:  aws.String(subject),
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
		Name:        name,
		HTMLPart:    html,
		TextPart:    html,
		SubjectPart: subject,
	})
}

func PutTemplate(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "AWS Ses keys not set.",
		})
		return
	}

	store, err := templates.NewSesTemplateStore(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logrus.Errorln(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "SES keys are incorrect.",
		})
		return
	}

	name := c.Param("name")
	html := c.PostForm("content")
	subject := c.PostForm("subject")

	_, err = store.UpdateTemplate(&ses.UpdateTemplateInput{
		Template: &ses.Template{
			TemplateName: aws.String(name),
			HtmlPart:     aws.String(html),
			TextPart:     aws.String(html),
			SubjectPart:  aws.String(subject),
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
		HTMLPart:    html,
		TextPart:    html,
		SubjectPart: subject,
	})
}

func DeleteTemplate(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "AWS Ses keys not set.",
		})
		return
	}

	store, err := templates.NewSesTemplateStore(keys.AccessKey, keys.SecretKey, keys.Region)
	if err != nil {
		logrus.Errorln(err.Error())
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
