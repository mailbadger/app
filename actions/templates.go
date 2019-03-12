package actions

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/emails/sesclient"
	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/routes/middleware"
	"github.com/news-maily/api/storage"
)

func GetTemplate(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": "AWS Ses keys not set.",
		})
		return
	}

	name := c.Param("name")

	client, err := sesclient.NewSESClient(keys.AccessKey, keys.SecretKey, "eu-west-1")

	res, err := client.GetTemplate(&ses.GetTemplateInput{
		TemplateName: aws.String(name),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": "Template not found.",
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
			"reason": "AWS Ses keys not set.",
		})
		return
	}

	nextToken := c.Query("next_token")

	client, err := sesclient.NewSESClient(keys.AccessKey, keys.SecretKey, "eu-west-1")

	res, err := client.ListTemplates(&ses.ListTemplatesInput{
		NextToken: aws.String(nextToken),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": "Templates not found.",
		})
		return
	}

	var list []entities.TemplateMeta

	for _, t := range res.TemplatesMetadata {
		list = append(list, entities.TemplateMeta{
			Name:      *t.Name,
			Timestamp: *t.CreatedTimestamp,
		})
	}

	c.JSON(http.StatusOK, entities.TemplateCollection{
		NextToken: *res.NextToken,
		List:      list,
	})
}

func PostTemplate(c *gin.Context) {
	u := middleware.GetUser(c)

	keys, err := storage.GetSesKeys(c, u.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "AWS Ses keys not set.",
		})
		return
	}

	client, err := sesclient.NewSESClient(keys.AccessKey, keys.SecretKey, "eu-west-1")

	name := c.PostForm("name")
	html := c.PostForm("content")
	subject := c.PostForm("subject")

	_, err = client.CreateTemplate(&ses.CreateTemplateInput{
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
			"reason": reason,
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
			"reason": "AWS Ses keys not set.",
		})
		return
	}

	client, err := sesclient.NewSESClient(keys.AccessKey, keys.SecretKey, "eu-west-1")

	name := c.PostForm("name")
	html := c.PostForm("content")
	subject := c.PostForm("subject")

	_, err = client.UpdateTemplate(&ses.UpdateTemplateInput{
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
			"reason": reason,
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
			"reason": "AWS Ses keys not set.",
		})
		return
	}

	client, err := sesclient.NewSESClient(keys.AccessKey, keys.SecretKey, "eu-west-1")

	name := c.PostForm("name")

	_, err = client.DeleteTemplate(&ses.DeleteTemplateInput{
		TemplateName: aws.String(name),
	})

	if err != nil {
		reason := "Unable to delete template."

		if awsErr, ok := err.(awserr.Error); ok {
			reason = fmt.Sprintf("%s: %s", awsErr.Code(), awsErr.Message())
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"reason": reason,
		})
		return
	}

	c.Status(http.StatusNoContent)
}
