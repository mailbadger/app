package routes

import (
	"net/http"
	"time"

	"github.com/FilipNikolovski/news-maily/actions"
	"github.com/FilipNikolovski/news-maily/routes/middleware"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
)

// New creates a new HTTP handler with the specified middleware.
func New() http.Handler {
	handler := gin.New()

	handler.Use(gin.Recovery())
	handler.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true))
	handler.Use(middleware.Storage())
	handler.Use(middleware.SetUser())

	// Guest routes
	handler.POST("/api/login", actions.PostLogin)

	// Authorized routes
	authorized := handler.Group("/api")
	authorized.Use(middleware.Authorized())
	{
		users := authorized.Group("/users")
		{
			users.GET("", actions.GetMe)
		}

		templates := authorized.Group("/templates")
		{
			templates.GET("", middleware.Paginate(), actions.GetTemplates)
			templates.GET("/:id", actions.GetTemplate)
			templates.POST("", actions.PostTemplate)
			templates.PUT("/:id", actions.PutTemplate)
			templates.DELETE("/:id", actions.DeleteTemplate)
		}

		campaigns := authorized.Group("/campaigns")
		{
			campaigns.GET("", middleware.Paginate(), actions.GetCampaigns)
			campaigns.GET("/:id", actions.GetCampaign)
			campaigns.POST("", actions.PostCampaign)
			campaigns.PUT("/:id", actions.PutCampaign)
			campaigns.DELETE("/:id", actions.DeleteCampaign)
		}

		lists := authorized.Group("/lists")
		{
			lists.GET("", middleware.Paginate(), actions.GetLists)
			lists.GET("/:id", actions.GetList)
			lists.POST("", actions.PostList)
			lists.PUT("/:id", actions.PutList)
			lists.DELETE("/:id", actions.DeleteList)
		}
	}

	return handler
}
