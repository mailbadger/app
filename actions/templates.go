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

func GetTemplates(c *gin.Context) {
	all := c.Query("all")
	if all != "" {
		if t, err := storage.GetAllTemplates(c, middleware.GetUser(c).Id); err == nil {
			c.JSON(http.StatusOK, t)
			return
		}
	}

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

	storage.GetTemplates(c, middleware.GetUser(c).Id, p)
	c.JSON(http.StatusOK, p)
}

func GetTemplate(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		if t, err := storage.GetTemplate(c, id, middleware.GetUser(c).Id); err == nil {
			c.JSON(http.StatusOK, t)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"reason": "Template not found",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"reason": "Id must be an integer",
	})
	return
}

func PostTemplate(c *gin.Context) {
	name, content := c.PostForm("name"), c.PostForm("content")

	_, err := storage.GetTemplateByName(c, name, middleware.GetUser(c).Id)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": "Template with that name already exists",
		})
		return
	}

	t := &entities.Template{
		Name:    name,
		Content: content,
		UserId:  middleware.GetUser(c).Id,
	}

	if !t.Validate() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": "Invalid data",
			"errors": t.Errors,
		})
		return
	}

	err = storage.CreateTemplate(c, t)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"reason": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func PutTemplate(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		t, err := storage.GetTemplate(c, id, middleware.GetUser(c).Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Template not found",
			})
			return
		}

		name, content := c.PostForm("name"), c.PostForm("content")

		t.Name = name
		t.Content = content

		if !t.Validate() {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Invalid data",
				"errors": t.Errors,
			})
			return
		}

		t2, err := storage.GetTemplateByName(c, name, middleware.GetUser(c).Id)
		if err == nil && t2.Id != t.Id {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Template with that name already exists",
			})
			return
		}

		err = storage.UpdateTemplate(c, t)
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

func DeleteTemplate(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 32); err == nil {
		user := middleware.GetUser(c)
		campaigns, err := storage.GetCampaignsByTemplateId(c, id, user.Id)
		if err == nil && len(campaigns) > 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"reason": "Cannot delete template because it is used by campaigns.",
			})
			return
		}

		err = storage.DeleteTemplate(c, id, user.Id)
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
