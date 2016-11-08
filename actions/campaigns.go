package actions

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/routes/middleware"
	"github.com/FilipNikolovski/news-maily/storage"
	"github.com/FilipNikolovski/news-maily/utils/pagination"
	"github.com/gin-gonic/gin"
)

func GetCampaigns(c *gin.Context) {
	val, ok := c.Get("pagination")
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("cannot create pagination object"))
		return
	}

	p, ok := val.(*pagination.Pagination)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("cannot cast pagination object"))
		return
	}

	storage.GetCampaigns(c, middleware.GetUser(c).Id, p)
	c.JSON(http.StatusOK, p)
}

func GetCampaign(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		if campaign, err := storage.GetCampaign(c, id, middleware.GetUser(c).Id); err == nil {
			c.JSON(http.StatusOK, campaign)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"reason": "Campaign not found",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Id must be an integer",
	})
	return
}

func PostCampaign(c *gin.Context) {
	name, subject, templateID := c.PostForm("name"), c.PostForm("subject"), c.PostForm("template_id")
	user := middleware.GetUser(c)

	if id, err := strconv.ParseInt(templateID, 10, 32); err == nil {
		t, err := storage.GetTemplate(c, id, user.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Template not found",
			})
			return
		}

		campaign := &entities.Campaign{
			Name:     name,
			Subject:  subject,
			UserId:   user.Id,
			Template: *t,
			Status:   entities.STATUS_DRAFT,
		}

		if !campaign.Validate() {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Invalid data",
				"errors": campaign.Errors,
			})
			return
		}

		err = storage.CreateCampaign(c, campaign)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, campaign)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Template id must be an integer",
	})
	return
}

func PutCampaign(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		user := middleware.GetUser(c)

		campaign, err := storage.GetCampaign(c, id, user.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Campaign not found",
			})
			return
		}

		name, subject, templateID := c.PostForm("name"), c.PostForm("subject"), c.PostForm("template_id")

		if tID, err := strconv.ParseInt(templateID, 10, 32); err == nil {
			t, err := storage.GetTemplate(c, tID, user.Id)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"reason": "Template not found",
				})
				return
			}

			campaign.Name = name
			campaign.Subject = subject
			campaign.Template = *t

			if !campaign.Validate() {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"reason": "Invalid data",
					"errors": campaign.Errors,
				})
				return
			}

			err = storage.UpdateCampaign(c, campaign)

			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"reason": err.Error(),
				})
				return
			}

			c.Status(http.StatusNoContent)
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "Template id must be an integer",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Id must be an integer",
	})
	return
}

func DeleteCampaign(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		user := middleware.GetUser(c)

		_, err := storage.GetCampaign(c, id, user.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Campaign not found",
			})
			return
		}

		err = storage.DeleteCampaign(c, id, user.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": err.Error(),
			})
			return
		}

		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Id must be an integer",
	})
	return
}
