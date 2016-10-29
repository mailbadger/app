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
	handler.POST("/login", actions.PostLogin)

	// Authorized routes
	users := handler.Group("/api/users")
	{
		users.Use(middleware.Authorized())
		users.GET("", actions.GetMe)
	}

	templates := handler.Group("/api/templates")
	{
		templates.Use(middleware.Authorized())
		templates.GET("", middleware.Paginate(), actions.GetTemplates)
	}

	return handler
}
