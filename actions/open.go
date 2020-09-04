package actions

import (
	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/logger"
	"github.com/mailbadger/app/storage"
	"net/http"
	"strconv"
)

func GetCampaignOpens(c *gin.Context) {
	val, ok := c.Get("cursor")
	if !ok {
		logger.From(c).Error("Unable to fetch pagination cursor from context.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch campaigns. Please try again.",
		})
		return
	}

	p, ok := val.(*storage.PaginationCursor)
	if !ok {
		logger.From(c).Error("Unable to cast pagination cursor from context value.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to fetch campaign opens. Please try again.",
		})
		return
	}

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		if err := storage.GetCampaignOpens(c, id, p); err == nil {
			c.JSON(http.StatusOK, p)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Campaign opens not found",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Id must be an integer",
	})
}
