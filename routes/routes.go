package routes

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/news-maily/api/actions"
	"github.com/news-maily/api/routes/middleware"
	"github.com/sirupsen/logrus"
)

// New creates a new HTTP handler with the specified middleware.
func New() http.Handler {
	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = logrus.InfoLevel
	}

	log := logrus.New()
	log.SetLevel(lvl)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	handler := gin.New()

	handler.Use(gin.Recovery())
	handler.Use(ginrus.Ginrus(log, time.RFC3339, true))
	handler.Use(middleware.Storage())
	handler.Use(middleware.Producer())
	handler.Use(middleware.SetUser())

	// Web app
	appDir := os.Getenv("APP_DIR")
	if appDir == "" {
		logrus.Panic("app directory not set")
	}

	handler.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Not found",
			})
			return
		}

		c.File(appDir + "/index.html")
	})

	// Assets
	handler.Static("/static", appDir+"/static")

	// Guest routes
	handler.POST("/api/authenticate", actions.PostLogin)
	handler.POST("/api/hooks", actions.HandleHook)

	// Authorized routes
	authorized := handler.Group("/api")
	authorized.Use(middleware.Authorized())
	{
		users := authorized.Group("/users")
		{
			users.GET("", actions.GetMe)
			users.POST("/password", actions.ChangePassword)
		}

		templates := authorized.Group("/templates")
		{
			templates.GET("", middleware.Paginate(), actions.GetTemplates)
			templates.GET("/:name", actions.GetTemplate)
			templates.POST("", actions.PostTemplate)
			templates.PUT("/:name", actions.PutTemplate)
			templates.DELETE("/:name", actions.DeleteTemplate)
		}

		campaigns := authorized.Group("/campaigns")
		{
			campaigns.GET("", middleware.Paginate(), actions.GetCampaigns)
			campaigns.GET("/:id", actions.GetCampaign)
			campaigns.POST("", actions.PostCampaign)
			campaigns.PUT("/:id", actions.PutCampaign)
			campaigns.DELETE("/:id", actions.DeleteCampaign)
			campaigns.POST("/:id/start", actions.StartCampaign)
		}

		lists := authorized.Group("/lists")
		{
			lists.GET("", middleware.Paginate(), actions.GetLists)
			lists.GET("/:id", actions.GetList)
			lists.POST("", actions.PostList)
			lists.PUT("/:id", actions.PutList)
			lists.DELETE("/:id", actions.DeleteList)
			lists.PUT("/:id/subscribers", actions.PutListSubscribers)
			lists.GET("/:id/subscribers", middleware.Paginate(), actions.GetListSubscribers)
			lists.DELETE("/:id/subscribers", actions.DetachListSubscribers)
		}

		subscribers := authorized.Group("/subscribers")
		{
			subscribers.GET("", middleware.Paginate(), actions.GetSubscribers)
			subscribers.GET("/:id", actions.GetSubscriber)
			subscribers.POST("", actions.PostSubscriber)
			subscribers.PUT("/:id", actions.PutSubscriber)
			subscribers.DELETE("/:id", actions.DeleteSubscriber)
		}

		sesKeys := authorized.Group(("/ses-keys"))
		{
			sesKeys.GET("", actions.GetSESKeys)
			sesKeys.POST("", actions.PostSESKeys)
			sesKeys.DELETE("", actions.DeleteSESKeys)
		}
	}

	return handler
}
