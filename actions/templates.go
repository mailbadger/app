package actions

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/entities/params"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/routes/middleware"
	templatesvc "github.com/mailbadger/app/services/templates"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/storage/s3"
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
	service := templatesvc.NewTemplateService(storage.GetFromContext(c), s3.GetFromContext(c))

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

	template := &entities.Template{
		UserID:      u.ID,
		Name:        body.Name,
		HTMLPart:    body.HTMLPart,
		TextPart:    body.TextPart,
		SubjectPart: body.Subject,
	}

	_, err := storage.GetTemplateByName(c, template.Name, u.ID)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Template with that name already exists",
		})
		return
	}

	err = service.AddTemplate(c, template)
	if err != nil {
		switch {
		case errors.Is(err, templatesvc.ErrParseHTMLPart):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to create template, failed to parse html_part",
			})
		case errors.Is(err, templatesvc.ErrParseTextPart):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to create template, failed to parse text_part",
			})
		case errors.Is(err, templatesvc.ErrParseSubjectPart):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to create template, failed to parse subject_part",
			})
		default:
			logger.From(c).WithFields(logrus.Fields{
				"template": template,
				"user_id":  u.ID,
			}).WithError(err).Error("Unable to create template")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to create template, please try again.",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, template)
}

func PutTemplate(c *gin.Context) {
	u := middleware.GetUser(c)
	service := templatesvc.NewTemplateService(storage.GetFromContext(c), s3.GetFromContext(c))

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
		return
	}

	template, err := storage.GetTemplate(c, id, u.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Template not found",
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

	template2, err := storage.GetTemplateByName(c, body.Name, u.ID)
	if err == nil && template.ID != template2.ID {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Template with that name already exists",
		})
		return
	}
	template.Name = body.Name
	template.HTMLPart = body.HTMLPart
	template.TextPart = body.TextPart
	template.SubjectPart = body.Subject

	err = service.UpdateTemplate(c, template)
	if err != nil {
		switch {
		case errors.Is(err, templatesvc.ErrParseHTMLPart):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to update template, failed to parse html_part",
			})
		case errors.Is(err, templatesvc.ErrParseTextPart):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to update template, failed to parse text_part",
			})
		case errors.Is(err, templatesvc.ErrParseSubjectPart):
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to update template, failed to parse subject_part",
			})
		default:
			logger.From(c).WithFields(logrus.Fields{
				"template": template,
				"user_id":  u.ID,
			}).WithError(err).Error("Unable to update template")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to update template, please try again.",
			})
		}
	}

	c.JSON(http.StatusOK, template)
}

func DeleteTemplate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Id must be an integer",
		})
		return
	}

	u := middleware.GetUser(c)
	service := templatesvc.NewTemplateService(storage.GetFromContext(c), s3.GetFromContext(c))

	err = service.DeleteTemplate(c, id, u.ID)
	if err != nil {
		logger.From(c).WithFields(logrus.Fields{
			"user_id":     u.ID,
			"template_id": id,
		}).WithError(err).Error("Unable to delete template.")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to delete template.",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
